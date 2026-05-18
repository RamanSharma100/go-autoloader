package tests

import (
	"testing"

	autoloader "github.com/RamanSharma100/go-autoloader"
	"github.com/RamanSharma100/go-autoloader/core"
)

func TestGetModule(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	result := engine.Modules.Get("auth")
	if result == nil {
		t.Fatal("expected module, got nil")
	}

	module, ok := result.(*core.ModuleHandle)
	if !ok {
		t.Fatalf("expected *core.ModuleHandle, got %T", result)
	}

	if module.Module == nil {
		t.Fatal("expected module.Module to be non-nil")
	}
}

func TestGetModuleByFilename(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	result := engine.Modules.Get("auth.go")
	if result == nil {
		t.Fatal("expected module, got nil")
	}

	_, ok := result.(*core.ModuleHandle)
	if !ok {
		t.Fatalf("expected *core.ModuleHandle, got %T", result)
	}
}

func TestGetFunction(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	result := engine.Modules.Get("auth.Login")
	if result == nil {
		t.Fatal("expected symbol for auth.Login, got nil — ensure auth.go exports a Login function and AutoParse is true")
	}

	s, ok := result.(*core.Symbol)
	if !ok {
		t.Fatalf("expected *core.Symbol, got %T", result)
	}

	if s.Type != "function" {
		t.Fatalf("expected type 'function', got '%s'", s.Type)
	}
}

func TestHasFunction(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	if !engine.Modules.HasFunction("auth.Login") {
		t.Fatal("expected HasFunction to return true for auth.Login")
	}

	result := engine.Modules.Get("auth")
	if result == nil {
		t.Fatal("expected module for 'auth', got nil")
	}

	handle, ok := result.(*core.ModuleHandle)
	if !ok {
		t.Fatalf("expected *core.ModuleHandle, got %T", result)
	}

	if !handle.HasFunction("Login") {
		t.Fatal("expected handle.HasFunction to return true for Login")
	}
}

func TestChainModuleFunctions(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	result := engine.Modules.Get("auth")
	if result == nil {
		t.Fatal("expected module, got nil")
	}

	handle, ok := result.(*core.ModuleHandle)
	if !ok {
		t.Fatalf("expected *core.ModuleHandle, got %T", result)
	}

	err := handle.Chain("Login", "Logout").Exec()
	if err != nil {
		t.Fatalf("chain execution failed: %v", err)
	}
}

func TestRegistryChain(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	err := engine.Modules.
		Chain("auth.Login", "auth.Logout").
		Exec()

	if err != nil {
		t.Fatalf("registry chain failed: %v", err)
	}
}
