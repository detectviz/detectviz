package fakes

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// FakeLifecycleModule is a fake implementation of modules.LifecycleModule for testing.
// zh: FakeLifecycleModule 是測試用的 LifecycleModule 假實作。
type FakeLifecycleModule struct {
	RunCalls      []context.Context
	ShutdownCalls []context.Context
	RunError      error
	ShutdownError error
	IsHealthy     bool
}

func (f *FakeLifecycleModule) Run(ctx context.Context) error {
	f.RunCalls = append(f.RunCalls, ctx)
	return f.RunError
}

func (f *FakeLifecycleModule) Shutdown(ctx context.Context) error {
	f.ShutdownCalls = append(f.ShutdownCalls, ctx)
	return f.ShutdownError
}

func (f *FakeLifecycleModule) Healthy() bool {
	return f.IsHealthy
}

// FakeModuleEngine is a fake implementation of modules.ModuleEngine for testing.
// zh: FakeModuleEngine 是測試用的 ModuleEngine 假實作。
type FakeModuleEngine struct {
	RegisterCalls    []modules.LifecycleModule
	RunAllCalls      []context.Context
	ShutdownAllCalls []context.Context
	RunAllError      error
	ShutdownAllError error
}

func (f *FakeModuleEngine) Register(m modules.LifecycleModule) {
	f.RegisterCalls = append(f.RegisterCalls, m)
}

func (f *FakeModuleEngine) RunAll(ctx context.Context) error {
	f.RunAllCalls = append(f.RunAllCalls, ctx)
	return f.RunAllError
}

func (f *FakeModuleEngine) ShutdownAll(ctx context.Context) error {
	f.ShutdownAllCalls = append(f.ShutdownAllCalls, ctx)
	return f.ShutdownAllError
}

// FakeModuleRegistry is a fake implementation of modules.ModuleRegistry for testing.
// zh: FakeModuleRegistry 是測試用的 ModuleRegistry 假實作。
type FakeModuleRegistry struct {
	RegisterCalls []RegisterCall
	GetCalls      []string
	ListCalls     int
	RegisterError error
	GetResult     modules.LifecycleModule
	GetFound      bool
	ListResult    []string
}

type RegisterCall struct {
	Name   string
	Module modules.LifecycleModule
}

func (f *FakeModuleRegistry) Register(name string, m modules.LifecycleModule) error {
	f.RegisterCalls = append(f.RegisterCalls, RegisterCall{
		Name:   name,
		Module: m,
	})
	return f.RegisterError
}

func (f *FakeModuleRegistry) Get(name string) (modules.LifecycleModule, bool) {
	f.GetCalls = append(f.GetCalls, name)
	return f.GetResult, f.GetFound
}

func (f *FakeModuleRegistry) List() []string {
	f.ListCalls++
	return f.ListResult
}
