package tests

import (
	"testing"

	"github.com/RamanSharma100/go-autoloader/file"
)

func TestScan(t *testing.T) {
	modules, err := file.Scan("../examples")

	if err != nil {
		t.Fatal(err)
	}

	if len(modules) == 0 {
		t.Fatal("expected modules")
	}
}
