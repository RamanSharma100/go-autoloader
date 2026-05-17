package core

import (
	"fmt"
	"plugin"

	"github.com/RamanSharma100/go-autoloader/runtime"
)

type ModuleHandle struct {
	Module *Module
}

type Function struct {
	Name       string
	Exported   bool
	Parameters []string
	Returns    []string
}

type Struct struct {
	Name     string
	Exported bool
}

type Variable struct {
	Name     string
	Exported bool
}

type Constant struct {
	Name     string
	Exported bool
}

type Interface struct {
	Name     string
	Exported bool
}

type Module struct {
	Name        string
	Path        string
	RelativeDir string
	RuntimePath string
	RuntimeName string
	PackageName string

	Extension string
	Kind      FileKind
	Content   []byte

	Functions  []Function
	Structs    []Struct
	Variables  []Variable
	Constants  []Constant
	Interfaces []Interface

	Symbols    []*Symbol
	PluginPath string
	Plugin     *plugin.Plugin
	Loader     *runtime.Loader
}

func (h *ModuleHandle) Get(name string) *Symbol {
	return h.Module.Get(name)
}

func (h *ModuleHandle) HasFunction(name string) bool {
	s := h.Module.Get(name)
	return s != nil && s.Type == "function"
}

func (m *ModuleHandle) Call(name string, args ...any) error {
	s := m.Get(name)

	if s == nil {
		return fmt.Errorf("function not found: %s", name)
	}

	if s.Type != "function" {
		return fmt.Errorf("not callable: %s", name)
	}

	fmt.Printf("[RUNTIME] calling: %s\n", name)

	if m.Module.Plugin != nil {
		err := m.Module.Loader.CallPlugin(m.Module.Plugin, m.Module.Path, name, args...)
		if err != nil {
			return fmt.Errorf("plugin execution failed: %w", err)
		}
		return nil
	}

	if m.Module.Plugin == nil && m.Module.Path != "" && m.Module.Loader != nil {
		err := m.Module.Loader.CallPlugin(nil, m.Module.Path, name, args...)
		if err != nil {
			return fmt.Errorf("windows fallback execution failed: %w", err)
		}
		return nil
	}

	s = m.Get(name)
	if s == nil {
		return fmt.Errorf("function not found: %s", name)
	}

	if s.Type != "function" {
		return fmt.Errorf("not callable: %s", name)
	}

	if s.Value != nil {
		if fn, ok := s.Value.(func(...any)); ok {
			fn(args...)
			return nil
		}
	}

	return nil
}
