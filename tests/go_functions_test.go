package tests

import (
	"testing"

	"github.com/RamanSharma100/go-autoloader/core"
	"github.com/RamanSharma100/go-autoloader/parsers"
)

func TestGoFunctions(t *testing.T) {
	module := &core.Module{
		Path: "../examples/routes/auth.go",
	}

	err := parsers.ParseGoFile(module)
	if err != nil {
		t.Fatal(err)
	}

	if len(module.Functions) == 0 {
		t.Fatal("expected functions")
	}
}
