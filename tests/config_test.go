package tests

import (
	"testing"

	"github.com/RamanSharma100/go-autoloader/core"
)

func TestDefaultConfig(t *testing.T) {
	config := core.DefaultConfig()

	if !config.AutoParse {
		t.Fatal("autoparse should be enabled")
	}
}
