package file

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RamanSharma100/go-autoloader/core"
	"github.com/RamanSharma100/go-autoloader/parsers"
	"github.com/RamanSharma100/go-autoloader/runtime"
)

func Run(m *core.Module, autoParse bool, loader runtime.Loader) {
	if m.Extension == "" {
		m.Extension = filepath.Ext(m.Path)
	}

	if m.Kind == "" {
		m.Kind = Detect(m.Path)
	}

	if m.Content == nil {
		data, err := os.ReadFile(m.Path)
		if err == nil {
			m.Content = data
		}
	}

	switch m.Kind {
	case core.CodeKind:
		fmt.Println("[DETECTED CODE FILE]")
		fmt.Printf("\n[LOADED %s] \n", m.Path)

		if m.Extension == ".go" {
			fmt.Println("[DETECTED GOLANG]")
			fmt.Println("[PARSING GOLANG]")
			m.Loader = &loader
			pluginPath, err := loader.BuildAndLoad(m.Path)
			println("pluginPath", pluginPath)
			if err == nil {
				p, err := loader.LoadPlugin(pluginPath)
				if err == nil {
					m.Plugin = p
					m.PluginPath = pluginPath
				}
			}

			fmt.Println("[MAKING PLUGIN]")

			if autoParse {
				err := parsers.ParseGoFile(m)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	case core.JSONKind:
		fmt.Println("[DETECTED JSON FILE]")
		fmt.Printf("\n[LOADED %s]", m.Path)
	case core.TextKind:
		fmt.Println("[DETECTED TEXT FILE]")
		fmt.Printf("\n[LOADED %s]", m.Path)
	case core.BinaryKind:
		fmt.Println("[DETECTED BINARY FILE]")
		fmt.Printf("\n[LOADED %s]", m.Path)
	default:
		fmt.Println("[DETECTED UNKOWN FILE]")
		fmt.Printf("\n[LOADED Unkown %s]", m.Path)
	}
}
