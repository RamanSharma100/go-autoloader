package core

import "fmt"

type Symbol struct {
	Name   string
	Type   string
	Module *Module

	Value any
}

func (m *Module) Get(name string) *Symbol {
	for _, fn := range m.Functions {
		if fn.Name == name {
			return &Symbol{
				Name:   fn.Name,
				Type:   "function",
				Module: m,
			}
		}
	}

	for _, s := range m.Structs {
		if s.Name == name {
			return &Symbol{
				Name:   s.Name,
				Type:   "struct",
				Module: m,
			}
		}
	}

	for _, v := range m.Variables {
		if v.Name == name {
			return &Symbol{
				Name:   v.Name,
				Type:   "variable",
				Module: m,
			}
		}
	}

	for _, c := range m.Constants {
		if c.Name == name {
			return &Symbol{
				Name:   c.Name,
				Type:   "constant",
				Module: m,
			}
		}
	}

	for _, i := range m.Interfaces {
		if i.Name == name {
			return &Symbol{
				Name:   i.Name,
				Type:   "interface",
				Module: m,
			}
		}
	}

	return nil
}

func (m *Module) HasFunction(name string) bool {
	if m == nil {
		return false
	}

	s := m.Get(name)
	return s != nil && s.Type == "function"
}

func (m *Module) Call(name string, args ...any) error {
	s := m.Get(name)

	if s == nil {
		return fmt.Errorf("function not found: %s", name)
	}

	if s.Type != "function" {
		return fmt.Errorf("not callable: %s", name)
	}

	fmt.Printf("[RUNTIME] calling: %s\n", name)

	if m.Plugin != nil {
		err := m.Loader.CallPlugin(m.Plugin, m.Path, name, args...)
		if err != nil {
			return fmt.Errorf("plugin execution failed: %w", err)
		}
		return nil
	}

	if m.Plugin == nil && m.Path != "" && m.Loader != nil {
		err := m.Loader.CallPlugin(nil, m.Path, name, args...)
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
