package tests

import (
	"testing"

	"github.com/RamanSharma100/go-autoloader/core"
	"github.com/RamanSharma100/go-autoloader/file"
	"github.com/RamanSharma100/go-autoloader/runtime"
)

func TestRun(t *testing.T) {
	module := &core.Module{
		Path: "../examples/routes/auth.go",
	}

	file.Run(module, true, *runtime.NewLoader("../examples/routes/", ".plugins"))

	if len(module.Functions) == 0 {
		t.Fatal("expected parsed functions")
	}
}
