package core

type Engine struct {
	Modules *ModuleRegistry
	Config  Config
}

func New(config Config) *Engine {
	return &Engine{
		Modules: NewModuleRegistry(),
		Config:  config,
	}
}
