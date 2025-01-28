package module

type Registry interface {
	Register(name string, factory ModuleFactory)
	CreateModules(cfg *ModuleConfig) []Module
}

type ModuleRegistry struct {
	modules map[string]ModuleFactory
}

type ModuleFactory func(*ModuleConfig) Module

var NewRegistry = func() Registry {
	return &ModuleRegistry{
		modules: make(map[string]ModuleFactory),
	}
}

func (r *ModuleRegistry) Register(name string, factory ModuleFactory) {
	r.modules[name] = factory
}

func (r *ModuleRegistry) CreateModules(cfg *ModuleConfig) []Module {
	var modules []Module
	for _, factory := range r.modules {
		modules = append(modules, factory(cfg))
	}
	return modules
}
