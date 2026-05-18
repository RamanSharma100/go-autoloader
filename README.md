# go-autoloader

> A hybrid AST + Plugin-based runtime module system for Go — built because nothing else did what was actually needed.

**Status:** WIP — Core working on Linux/macOS. Windows runtime calling in progress (see limitations).

---

## Why this exists

While building **Workflow Studio** — an API orchestration platform in Go — I needed a way to load Go files dynamically at runtime and call their functions without recompiling the whole binary or manually wiring every import.

Go is statically compiled. There is no `require("./routes/auth")` like in Node. There is no `importlib` like in Python. You write the import, you recompile, you wire it up manually. Every single time.

I looked for a library that could:

- Scan a directory of Go files
- Understand what functions, structs, and variables exist in them
- Let me call those functions at runtime by name — like `engine.Modules.Call("auth.Login")`
- Without me changing a single line in those Go files

Nothing did all of this together.

---

## What already exists — and why it wasn't enough

| Library                            | What it does                                | Why it wasn't enough                                                                                 |
| ---------------------------------- | ------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `plugin` (stdlib)                  | Loads `.so` files, looks up symbols         | Requires you to manually compile each file first, no scanning, no registry, no AST                   |
| `pkujhd/goloader`                  | Loads compiled `.o` object files at runtime | Works at object level, not source level — requires `go tool compile` manually, no directory scanning |
| `rainycape/dl`                     | dlopen/dlsym wrapper for C shared libs      | C libraries only, not Go source files                                                                |
| `golang.org/x/tools/go/loader`     | Loads and type-checks Go packages           | Analysis only — no execution, no plugin building, no function calling                                |
| `vladimirvivien/go-plugin-example` | Demo of Go's plugin package                 | Example code, not a library — no registry, no AST scanning                                           |

None of them give you a single `engine.Modules.Call("auth.Login")` from a raw `.go` file.

`go-autoloader` combines all three layers — scanning, parsing, and execution — into one cohesive runtime system.

---

## Architecture

```
Directory of .go files
        ↓
   [ Scanner ]          → walks files, skips .plugins/, builds module list
        ↓
   [ AST Parser ]       → extracts functions, structs, variables, constants, interfaces
        ↓
   [ Plugin Builder ]   → compiles .go → .so via go build -buildmode=plugin
        ↓
   [ Module Registry ]  → indexes everything by name with dot-notation access
        ↓
engine.Modules.Get("auth.Login")     → *Symbol
engine.Modules.Call("auth.Login")    → executes
engine.Modules.HasFunction("auth.Login") → bool
```

Three layers:

**1. Scanner** — recursively walks a directory, skips `.plugins/` and binary artifacts, builds a module list with `RuntimeName` and `RuntimePath` for every file.

**2. Parser (AST)** — uses `go/ast` to extract all exported and unexported symbols from each `.go` file. Builds a symbol table in memory. This is metadata only — nothing executes here.

**3. Runtime** — two execution modes:

- **Plugin mode (Linux/macOS):** compiles each `.go` to a `.so` using `go build -buildmode=plugin`, loads it via `plugin.Open`, and executes functions via `plugin.Lookup` + reflection.
- **Windows fallback (WIP):** rewrites the source file with a `main` package wrapper and streams it through `go run -` via stdin. Runtime function calling on Windows is currently unreliable — tracked below.

---

## Installation

```bash
go get github.com/RamanSharma100/go-autoloader
```

Requires Go 1.18+ and CGO enabled on Linux/macOS for plugin mode.

---

## Quick Start

```go
package main

import (
    autoloader "github.com/RamanSharma100/go-autoloader"
    "github.com/RamanSharma100/go-autoloader/core"
)

func main() {
    engine := autoloader.Load("./examples/routes")

    // get a module
    result := engine.Modules.Get("auth")
    handle := result.(*core.ModuleHandle)

    // check and call
    if handle.HasFunction("Login") {
        handle.Call("Login", "user@example.com")
    }

    // dot-notation symbol lookup
    sym := engine.Modules.Get("auth.Login")
    symbol := sym.(*core.Symbol)
    println(symbol.Type) // "function"

    // call directly from registry
    engine.Modules.Call("auth.Login", "user@example.com")

    // chain multiple functions on a module
    handle.Chain("Login", "Logout").
        With("Login", "user@example.com").
        Exec()

    // chain across modules from registry
    engine.Modules.Chain("auth.Login", "users.GetProfile").
        With("auth.Login", "user@example.com").
        Exec()
}
```

---

## Example Go File

No changes needed to your Go files. Just write normal Go:

```go
package routes

func Login(name string) {
    println("Login called with:", name)
}

func Logout() {
    println("Logout called")
}
```

---

## API Reference

### Load engine

```go
engine := autoloader.Load("./routes")

// or with config
engine := autoloader.LoadWithConfig("./routes", core.Config{
    AutoParse: true,
})
```

### Module access

```go
// get module handle
result := engine.Modules.Get("auth")
handle := result.(*core.ModuleHandle)

// get by filename
result := engine.Modules.Get("auth.go")

// get symbol directly
sym := engine.Modules.Get("auth.Login")
symbol := sym.(*core.Symbol)
println(symbol.Name, symbol.Type)
```

### Function checks

```go
engine.Modules.HasFunction("auth.Login")  // true/false from registry
handle.HasFunction("Login")               // true/false from module handle
```

### Calling functions

```go
// from registry
engine.Modules.Call("auth.Login", "arg1")

// from module handle
handle.Call("Login", "arg1")
```

### Chaining

```go
// single module chain
handle.Chain("Login", "Validate", "Logout").
    With("Login", "user@example.com").
    With("Validate", "token_xyz").
    Exec()

// cross-module chain from registry
engine.Modules.Chain("auth.Login", "users.GetProfile", "audit.Log").
    With("auth.Login", "user@example.com").
    Exec()

// check chain error
chain := handle.Chain("Login", "Logout")
if err := chain.Exec(); err != nil {
    println("chain failed:", err.Error())
}
```

---

## Windows — Current Status

**Plugin mode does not work on Windows.** Go's `plugin` package requires CGO and ELF shared object support — neither of which Windows provides natively.

The current Windows fallback rewrites each `.go` file with a `main` package wrapper and pipes it through `go run -` via stdin. This works for simple standalone function calls but has known issues:

- Functions with external package imports fail because the rewritten wrapper doesn't carry dependencies
- No argument type safety — everything is passed as `fmt.Sprintf("%#v", arg)` expressions
- `go run` subprocess overhead per call makes it impractical for chains or loops
- No symbol registry on Windows since plugins never build — `HasFunction` works (AST still parses), but `Call` may fail silently

**What is being worked on for Windows:**

- Pre-generating a proper `main.go` wrapper per module with correct imports extracted from AST
- Using `go run` with a temp directory instead of stdin streaming to support multi-file dependencies
- Alternatively: exploring `pkujhd/goloader` as a Windows-compatible execution backend that works at the object file level without CGO

If you are on Windows and only need metadata (scanning, AST, `HasFunction`, `Get`), everything works fine. Runtime execution via `Call` and `Chain` is unreliable until the above is resolved.

---

## Running Tests

```bash
# Linux/macOS
go test ./... -v

# Windows (metadata tests only — execution tests will fail)
go test ./... -v -run "TestGetModule|TestGetModuleByFilename|TestGetFunction|TestHasFunction"
```

---

## Running the Example

```bash
go run ./examples
```

## Output of Example

```bash
/app # cd examples/
/app/examples # go run main.go
[DETECTED JSON FILE]

[LOADED /app/examples/config/app.json][DETECTED CODE FILE]

[LOADED /app/examples/main.go]
[DETECTED GOLANG]
[PARSING GOLANG]
2026/05/18 08:34:18 [PLUGIN BUILT] /app/examples/.plugins/main.so
pluginPath /app/examples/.plugins/main.so
[MAKING PLUGIN]
[DETECTED TEXT FILE]

[LOADED /app/examples/notes/readme.txt][DETECTED CODE FILE]

[LOADED /app/examples/routes/admin.go]
[DETECTED GOLANG]
[PARSING GOLANG]
2026/05/18 08:34:18 [PLUGIN BUILT] /app/examples/.plugins/admin.so
pluginPath /app/examples/.plugins/admin.so
[MAKING PLUGIN]
[DETECTED CODE FILE]

[LOADED /app/examples/routes/auth.go]
[DETECTED GOLANG]
[PARSING GOLANG]
2026/05/18 08:34:18 [PLUGIN BUILT] /app/examples/.plugins/auth.so
pluginPath /app/examples/.plugins/auth.so
[MAKING PLUGIN]
[DETECTED CODE FILE]

[LOADED /app/examples/routes/users.go]
[DETECTED GOLANG]
[PARSING GOLANG]
2026/05/18 08:34:18 [PLUGIN BUILT] /app/examples/.plugins/users.so
pluginPath /app/examples/.plugins/users.so
[MAKING PLUGIN]

=== MODULES ===

Name: app
Path: /app/examples/config/app.json
Kind: json
Extension: .json

Name: main
Path: /app/examples/main.go
Kind: code
Extension: .go
Functions:
 - main

Name: readme
Path: /app/examples/notes/readme.txt
Kind: text
Extension: .txt

Name: admin
Path: /app/examples/routes/admin.go
Kind: code
Extension: .go
Functions:
 - CreateAdmin

Name: auth
Path: /app/examples/routes/auth.go
Kind: code
Extension: .go
Found Login
Calling Login
[RUNTIME] calling: Login
Login Worked!
Functions:
 - Login
 - Logout

Name: users
Path: /app/examples/routes/users.go
Kind: code
Extension: .go
Functions:
 - CreateUser

/app/examples
```

---

## Roadmap

**v0.1 — current**

- File scanning ✅
- AST parsing ✅
- Plugin build system ✅
- Module registry ✅
- Dot-notation symbol lookup ✅
- Function chaining ✅
- Runtime execution — Linux/macOS ✅
- Runtime execution — Windows 🔧 in progress

**v0.2**

- File watcher / hot reload
- Plugin caching (skip rebuild if source unchanged)
- Dependency graph between modules
- Windows execution via proper temp-dir runner

**v0.3**

- Middleware system (before/after call hooks)
- Module lifecycle events
- Exported symbol type introspection with full signatures

**v1.0**

- Full runtime orchestration
- Multi-module execution graph
- Production plugin ecosystem

---

## Design Principles

- AST is metadata only — it never executes anything
- Plugin is the execution engine on supported platforms
- Module is the runtime bridge between discovery and execution
- The registry is the single source of truth — one `Get()` call, consistent behavior
- No magic — everything is explicit, traceable, and debuggable

---

## Limitations

- Plugin mode requires Linux or macOS (Go stdlib constraint — `plugin` package is not supported on Windows)
- Functions must be exported (capital letter) for plugin execution via `plugin.Lookup`
- All files in the scanned directory get compiled — there is no selective include yet
- No hot reload yet (planned v0.2)
- Windows `Call` is unreliable for functions with external imports

---

## Philosophy

> Go should not feel static when you are building systems that need to grow at runtime. The compiler is a tool — not a ceiling.
