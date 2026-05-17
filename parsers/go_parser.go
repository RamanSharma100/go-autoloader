package parsers

import (
	"go/ast"
	goparser "go/parser"
	"go/token"

	"github.com/RamanSharma100/go-autoloader/core"
)

func ParseGoFile(module *core.Module) error {
	fset := token.NewFileSet()

	node, err := goparser.ParseFile(fset, module.Path, nil, 0)
	if err != nil {
		return err
	}

	if module.Symbols == nil {
		module.Symbols = []*core.Symbol{}
	}

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			module.Functions = append(module.Functions, core.Function{
				Name: d.Name.Name,
			})
			symbol := &core.Symbol{
				Name:   d.Name.Name,
				Type:   "function",
				Module: module,
			}

			module.Symbols = append(module.Symbols, symbol)
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					switch s.Type.(type) {
					case *ast.StructType:
						module.Structs = append(
							module.Structs,
							core.Struct{
								Name:     s.Name.Name,
								Exported: s.Name.IsExported(),
							},
						)
					case *ast.InterfaceType:
						module.Interfaces = append(
							module.Interfaces,
							core.Interface{
								Name:     s.Name.Name,
								Exported: s.Name.IsExported(),
							},
						)
					}
				case *ast.ValueSpec:
					for _, name := range s.Names {
						switch d.Tok {
						case token.VAR:
							module.Variables = append(
								module.Variables,
								core.Variable{
									Name:     name.Name,
									Exported: name.IsExported(),
								},
							)
						case token.CONST:
							module.Constants = append(
								module.Constants,
								core.Constant{
									Name:     name.Name,
									Exported: name.IsExported(),
								},
							)
						}
					}
				}
			}
		}
	}

	return nil
}
