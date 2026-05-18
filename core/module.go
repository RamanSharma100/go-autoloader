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
	if h == nil || h.Module == nil || name == "" {
		return nil
	}
	return h.Module.Get(name)
}

func (h *ModuleHandle) HasFunction(name string) bool {
	if h == nil || h.Module == nil || name == "" {
		return false
	}
	s := h.Module.Get(name)
	return s != nil && s.Type == "function"
}

func (h *ModuleHandle) Call(name string, args ...any) error {
	if h == nil || h.Module == nil {
		return fmt.Errorf("module handle is nil")
	}

	s := h.Get(name)
	if s == nil {
		return fmt.Errorf("function not found: %s", name)
	}
	if s.Type != "function" {
		return fmt.Errorf("not callable: %s is a %s", name, s.Type)
	}

	fmt.Printf("[RUNTIME] calling: %s\n", name)

	m := h.Module
	if m.Plugin != nil && m.Loader != nil {
		return m.Loader.CallPlugin(m.Plugin, m.Path, name, args...)
	}
	if m.Path != "" && m.Loader != nil {
		return m.Loader.CallPlugin(nil, m.Path, name, args...)
	}
	if s.Value != nil {
		if fn, ok := s.Value.(func(...any)); ok {
			fn(args...)
			return nil
		}
	}

	return fmt.Errorf("no executor available for: %s", name)
}

func (h *ModuleHandle) Chain(names ...string) *ChainHandle {
	return &ChainHandle{
		handle: h,
		calls:  names,
		args:   make(map[string][]any),
	}
}

type ChainHandle struct {
	handle *ModuleHandle
	calls  []string
	args   map[string][]any
	err    error
}

func (c *ChainHandle) With(name string, args ...any) *ChainHandle {
	if c.args == nil {
		c.args = make(map[string][]any)
	}
	c.args[name] = args
	return c
}

func (c *ChainHandle) Exec() error {
	if c.handle == nil || c.handle.Module == nil {
		return fmt.Errorf("chain: module handle is nil")
	}
	for _, name := range c.calls {
		args := c.args[name]
		if err := c.handle.Call(name, args...); err != nil {
			c.err = fmt.Errorf("chain stopped at %s: %w", name, err)
			return c.err
		}
	}
	return nil
}

func (c *ChainHandle) Error() error {
	return c.err
}
