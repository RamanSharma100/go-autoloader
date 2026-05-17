package tests

import (
	"testing"

	autoloader "github.com/RamanSharma100/go-autoloader"
	"github.com/RamanSharma100/go-autoloader/core"
)

func TestGetModule(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	module := engine.Modules.Get("auth")

	if module == nil {
		t.Fatal("expected module")
	}
}

func TestGetModuleByFilename(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	module := engine.Modules.Get("auth.go")

	if module == nil {
		t.Fatal("expected module")
	}
}

func TestGetFunction(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	sym := engine.Modules.Get("auth.Login")

	s, ok := sym.(*core.Symbol)
	if !ok {
		t.Fatalf("expected Symbol got %T", sym)
	}

	if s.Type != "function" {
		t.Fatal("expected function")
	}
}
func TestHasFunction(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	if !engine.Modules.HasFunction("auth.Login") {
		t.Fatal("expected function")
	}

	mod := engine.Modules.Get("auth")
	handle, ok := mod.(*core.ModuleHandle)
	if !ok {
		t.Fatal("expected ModuleHandle")
	}

	if !handle.HasFunction("Login") {
		t.Fatal("expected function")
	}
}
