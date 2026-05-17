package tests

import (
	"testing"

	"github.com/RamanSharma100/go-autoloader/core"
	"github.com/RamanSharma100/go-autoloader/parsers"
)

func TestGoSymbolParsing(t *testing.T) {
	module := &core.Module{
		Path: "../examples/routes/users.go",
	}

	err := parsers.ParseGoFile(module)
	if err != nil {
		t.Fatal(err)
	}

	if len(module.Functions) == 0 {
		t.Fatal("expected functions")
	}

	if len(module.Structs) == 0 {
		t.Fatal("expected structs")
	}
}
