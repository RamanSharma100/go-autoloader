package tests

import (
	"testing"

	autoloader "github.com/RamanSharma100/go-autoloader"
)

func TestCallFunction(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	err := engine.Modules.Call("auth.Login")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCallMissingFunction(t *testing.T) {
	engine := autoloader.Load("../examples/routes")

	err := engine.Modules.Call("auth.Unknown")

	if err == nil {
		t.Fatal("expected error")
	}
}
