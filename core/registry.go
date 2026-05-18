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

func isFileExtension(s string) bool {
	knownExts := map[string]bool{
		"go": true, "so": true, "json": true,
		"txt": true, "md": true, "yaml": true,
		"toml": true, "yml": true, "env": true,
		"js": true, "ts": true, "py": true,
	}
	return knownExts[strings.ToLower(s)]
}

func (r *ModuleRegistry) matchesModule(m *Module, query string) bool {
	q := strings.TrimSuffix(query, ".go")
	return query == m.RuntimeName ||
		query == m.RuntimePath ||
		query == m.Name ||
		q == m.RuntimeName ||
		q == m.RuntimePath ||
		query == m.RuntimePath+".go" ||
		query == m.RuntimeName+".go"
}

func (r *ModuleRegistry) bestModule(query string) *Module {
	var best *Module
	for _, m := range r.Items {
		if m == nil {
			continue
		}
		if r.matchesModule(m, query) {
			if best == nil {
				best = m
			} else if len(m.Functions)+len(m.Structs)+len(m.Variables)+len(m.Constants) >
				len(best.Functions)+len(best.Structs)+len(best.Variables)+len(best.Constants) {
				best = m
			}
		}
	}
	return best
}

func (r *ModuleRegistry) Get(query string) any {
	if query == "" {
		return nil
	}

	if !strings.Contains(query, ".") {
		m := r.bestModule(query)
		if m != nil {
			return &ModuleHandle{Module: m}
		}
		return nil
	}

	parts := strings.SplitN(query, ".", 2)
	if len(parts) != 2 {
		return nil
	}

	modulePart := parts[0]
	symbolPart := parts[1]

	if isFileExtension(symbolPart) {
		m := r.bestModule(query)
		if m != nil {
			return &ModuleHandle{Module: m}
		}
		return nil
	}

	var bestModule *Module
	for _, m := range r.Items {
		if m == nil {
			continue
		}
		if m.RuntimeName == modulePart || m.Name == modulePart {
			if bestModule == nil {
				bestModule = m
			} else if len(m.Functions)+len(m.Structs)+len(m.Variables)+len(m.Constants) >
				len(bestModule.Functions)+len(bestModule.Structs)+len(bestModule.Variables)+len(bestModule.Constants) {
				bestModule = m
			}
		}
	}

	if bestModule != nil {
		sym := bestModule.Get(symbolPart)
		if sym != nil {
			return sym
		}
	}

	return nil
}

func (r *ModuleRegistry) getSymbol(query string) *Symbol {
	if query == "" {
		return nil
	}

	result := r.Get(query)
	if result == nil {
		return nil
	}
	if sym, ok := result.(*Symbol); ok {
		return sym
	}
	return nil
}

func (r *ModuleRegistry) GetModule(query string) *ModuleHandle {
	result := r.Get(query)
	if result == nil {
		return nil
	}
	if h, ok := result.(*ModuleHandle); ok {
		return h
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
		return fmt.Errorf("not callable: %s is a %s", query, sym.Type)
	}

	fmt.Printf("[RUNTIME] calling function: %s\n", sym.Name)

	m := sym.Module
	if m == nil {
		return fmt.Errorf("symbol has no associated module: %s", sym.Name)
	}

	if m.Plugin != nil && m.Loader != nil {
		return m.Loader.CallPlugin(m.Plugin, m.Path, sym.Name, args...)
	}
	if m.Path != "" && m.Loader != nil {
		return m.Loader.CallPlugin(nil, m.Path, sym.Name, args...)
	}
	if sym.Value != nil {
		if fn, ok := sym.Value.(func(...any)); ok {
			fn(args...)
			return nil
		}
	}
	return fmt.Errorf("no executor available for: %s", sym.Name)
}

func (r *ModuleRegistry) Chain(queries ...string) *RegistryChainHandle {
	return &RegistryChainHandle{
		registry: r,
		calls:    queries,
		args:     make(map[string][]any),
	}
}

type RegistryChainHandle struct {
	registry *ModuleRegistry
	calls    []string
	args     map[string][]any
	err      error
}

func (c *RegistryChainHandle) With(query string, args ...any) *RegistryChainHandle {
	if c.args == nil {
		c.args = make(map[string][]any)
	}
	c.args[query] = args
	return c
}

func (c *RegistryChainHandle) Exec() error {
	for _, query := range c.calls {
		args := c.args[query]
		if err := c.registry.Call(query, args...); err != nil {
			c.err = fmt.Errorf("chain stopped at %s: %w", query, err)
			return c.err
		}
	}
	return nil
}

func (c *RegistryChainHandle) Error() error {
	return c.err
}
