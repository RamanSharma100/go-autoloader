package tests

import (
	"testing"

	autoloader "github.com/RamanSharma100/go-autoloader"
)

func TestLoad(t *testing.T) {
	engine := autoloader.Load("../examples")

	if engine == nil {
		t.Fatal("expected engine")
	}

	if len(engine.Modules.Items) == 0 {
		t.Fatal("expected modules")
	}
}
