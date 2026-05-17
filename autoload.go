package goautoloader

import (
	"log"
	goPath "path"
	"path/filepath"

	"github.com/RamanSharma100/go-autoloader/core"
	"github.com/RamanSharma100/go-autoloader/file"
	"github.com/RamanSharma100/go-autoloader/runtime"
)

func Load(path string) *core.Engine {
	return LoadWithConfig(
		path,
		core.DefaultConfig(),
	)
}

func LoadWithConfig(
	path string,
	config core.Config,
) *core.Engine {
	engine := core.New(config)

	rootAbs, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	loader := runtime.NewLoader(path, goPath.Join(rootAbs, "./.plugins"))

	modules, err := file.Scan(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, module := range modules {
		engine.Modules.Add(module)
	}

	for _, module := range modules {
		file.Run(module, config.AutoParse, *loader)
	}

	return engine
}
