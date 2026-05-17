package tests

import (
	"testing"

	"github.com/RamanSharma100/go-autoloader/file"
)

func TestDetectGoFile(t *testing.T) {
	kind := file.Detect("../examples/routes/auth.go")

	if kind != "code" {
		t.Fatal("expected code file")
	}
}
