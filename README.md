# go-autoloader [WIP => Fixing Runtime Function Calling for windows]

A hybrid **AST + Plugin-based runtime module system for Go**

`go-autoloader` scans Go projects, parses source files using AST, builds runtime metadata, and optionally compiles and executes modules using Go plugins.

It bridges the gap between:

- Static Go compilation
- Dynamic module discovery
- Runtime plugin execution

---

# 🚀 Why this exists

Go normally requires:

- manual imports
- manual wiring of modules
- static execution paths
- no runtime module discovery

This project enables:

> “Load Go code like a runtime system while still staying within Go’s constraints.”

---

# 🧠 Architecture

The system has 3 layers:

## 1. Scanner (Discovery Layer)

- Recursively scans files
- Detects Go source files
- Builds module registry

## 2. Parser (AST Layer)

- Parses Go using `go/ast`
- Extracts:
  - Functions
  - Structs
  - Interfaces
  - Variables
  - Constants

- Builds symbol table (metadata only)

## 3. Runtime Layer (Execution)

Two execution modes:

### A. Plugin Mode (Production)

- Compiles `.go` → `.so`
- Loads using `plugin.Open`
- Executes real runtime functions

### B. Symbol Fallback (Dev Mode)

- Uses in-memory symbol registry
- Executes bound function if available

---

# 📦 Installation

```bash
go get github.com/RamanSharma100/go-autoloader
```

---

# ⚡ Quick Start

```go
package main

import autoloader "github.com/RamanSharma100/go-autoloader"

func main() {

	engine := autoloader.Load("./examples")

	mod := engine.Modules.Get("auth").(*core.ModuleHandle)

	mod.Call("Login", "hello", 123)

	if mod.HasFunction("Login") {
		println("Login exists")
	}
}
```

---

# 📂 Example Go File

```go
package routes

func Login(name string) {
	println("Login called")
}

func Logout() {
	println("Logout called")
}
```

---

# ⚙️ Auto Behavior

When a Go file is loaded:

### AST parsing extracts:

- function names
- struct names
- variables
- interfaces

### Plugin system:

- compiles file using:

```bash
go build -buildmode=plugin
```

- loads `.so` dynamically
- enables runtime execution

---

# 🔧 Core API

## Load Engine

```go
engine := autoloader.Load("./routes")
```

---

## Get Module

```go
mod := engine.Modules.Get("auth").(*core.ModuleHandle)
```

---

## Call Function

```go
mod.Call("Login", "arg1", 123)
```

---

## Check Function

```go
mod.HasFunction("Login")
```

---

# 🔌 Plugin System

Internally uses:

```bash
go build -buildmode=plugin
```

Then loads:

```go
plugin.Open(path)
```

Execution flow:

```
Call()
  ↓
Check Plugin
  ↓
plugin.Lookup("Login")
  ↓
Execute function
```

---

# 🧪 Running Tests

```bash
go test ./...
```

---

# 🏗 Running Example

```bash
go run ./examples
```

---

# 📌 Design Principles

- AST = metadata only
- Plugin = execution engine
- Module = runtime bridge
- No duplication of execution logic
- Minimal abstraction layers
- Go-native plugin runtime

---

# 🧭 Limitations (Important)

- Go plugins only work on Linux/macOS/compatible builds
- Windows support is limited
- Functions must be exported for plugin execution
- No hot reload yet (planned)

---

# 🛣 Roadmap

## v0.1

- File scanning [Done]
- AST parsing [Done]
- Plugin build system [Done]
- Module registry [Done]
- Runtime execution [Under Progress for Windows]

## v0.2 [Upcoming]

- File watcher (hot reload)
- Plugin caching
- Dependency graph

## v0.3

- Middleware system
- Hooks (before/after call)
- Module lifecycle events

## v1.0

- Full runtime orchestration
- Multi-module execution graph
- Production plugin ecosystem

---

# 💡 Philosophy

> “Go should not feel static when building systems — it should behave like a runtime engine when needed.”
