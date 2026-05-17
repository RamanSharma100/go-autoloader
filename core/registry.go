package core

import (
	"fmt"
	"strings"
)

type ModuleRegistry struct {
	Items []*Module
}

func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{Items: []*Module{}}
}

func (r *ModuleRegistry) Add(module *Module) {
	r.Items = append(r.Items, module)
}

func (r *ModuleRegistry) getModule(query string) *Module {
	for _, module := range r.Items {
		if module == nil {
			continue
		}

		q := strings.TrimSuffix(query, ".go")

		if query == module.RuntimePath ||
			q == module.RuntimePath ||
			query == module.RuntimeName ||
			query == module.Name {
			return module
		}
	}
	return nil
}

func (r *ModuleRegistry) getSymbol(query string) *Symbol {
	for _, module := range r.Items {

		prefixes := []string{
			module.RuntimePath,
			module.RuntimePath + ".go",
			module.RuntimeName,
		}

		for _, prefix := range prefixes {
			if strings.HasPrefix(query, prefix+".") {
				name := strings.TrimPrefix(query, prefix+".")
				return module.Get(name)
			}
		}
	}
	return nil
}

func (r *ModuleRegistry) Get(query string) any {
	for _, m := range r.Items {
		if m == nil {
			continue
		}

		if query == m.RuntimeName ||
			query == m.RuntimePath ||
			query == m.Name ||
			query == m.RuntimePath+".go" {
			return &ModuleHandle{Module: m}
		}

		if strings.Contains(query, ".") {
			parts := strings.Split(query, ".")
			if len(parts) == 2 {
				if m.RuntimeName == parts[0] {
					return m.Get(parts[1])
				}
			}
		}
	}

	return nil
}

func (r *ModuleRegistry) HasFunction(query string) bool {
	sym := r.getSymbol(query)
	return sym != nil && sym.Type == "function"
}

func (r *ModuleRegistry) Call(query string, args ...any) error {
	sym := r.getSymbol(query)

	if sym == nil {
		return fmt.Errorf("function not found: %s", query)
	}

	if sym.Type != "function" {
		return fmt.Errorf("not callable: %s", query)
	}

	fmt.Printf("[RUNTIME] calling function: %s\n", sym.Name)

	m := sym.Module
	if m.Plugin != nil {
		err := m.Loader.CallPlugin(m.Plugin, m.Path, sym.Name, args...)
		if err != nil {
			return fmt.Errorf("plugin execution failed: %w", err)
		}
		return nil
	}

	if m.Plugin == nil && m.Path != "" && m.Loader != nil {
		err := m.Loader.CallPlugin(nil, m.Path, sym.Name, args...)
		if err != nil {
			return fmt.Errorf("windows fallback execution failed: %w", err)
		}
		return nil
	}

	s := m.Get(sym.Name)
	if s == nil {
		return fmt.Errorf("function not found: %s", sym.Name)
	}

	if s.Type != "function" {
		return fmt.Errorf("not callable: %s", sym.Name)
	}

	if s.Value != nil {
		if fn, ok := s.Value.(func(...any)); ok {
			fn(args...)
			return nil
		}
	}
	return nil
}
