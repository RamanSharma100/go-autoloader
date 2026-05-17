package tests

import (
	"testing"

	autoloader "github.com/RamanSharma100/go-autoloader"
	"github.com/RamanSharma100/go-autoloader/core"
)

func TestModule(t *testing.T) {
	module := &core.Module{
		Name: "auth.go",
	}

	if module.Name != "auth.go" {
		t.Fatal("invalid module")
	}
}

func TestModuleChaining(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	handle := engine.Modules.Get("auth")

	mod, ok := handle.(*core.ModuleHandle)
	if !ok {
		t.Fatal("expected ModuleHandle")
	}

	if !mod.HasFunction("Login") {
		t.Fatal("expected Login function")
	}
}
