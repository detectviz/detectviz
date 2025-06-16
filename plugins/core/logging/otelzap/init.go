package otelzap

import (
	"sync"

	"detectviz/internal/platform/registry"
	"detectviz/pkg/platform/contracts"
)

var (
	// globalRegistry holds the singleton registry instance
	globalRegistry contracts.Registry
	registryOnce   sync.Once
)

// getOrCreateRegistry returns the global registry instance
func getOrCreateRegistry() contracts.Registry {
	registryOnce.Do(func() {
		globalRegistry = registry.NewManager()
	})
	return globalRegistry
}

// GetGlobalRegistry returns the global registry for external access
func GetGlobalRegistry() contracts.Registry {
	return getOrCreateRegistry()
}

// init automatically registers the otelzap plugin
func init() {
	registry := getOrCreateRegistry()

	// Register plugin with the name expected by lifecycle manager
	if err := registry.RegisterPlugin("otelzap", NewOtelZapPlugin); err != nil {
		// Silently continue if registration fails to avoid startup issues
		return
	}

	// Register metadata
	metadata := &contracts.PluginMetadata{
		Name:        "otelzap",
		Version:     "1.0.0",
		Type:        "logger",
		Category:    "core",
		Description: "OpenTelemetry integrated Zap logger plugin for DetectViz",
		Author:      "DetectViz Team",
		License:     "MIT",
		Config: map[string]any{
			"enabled":         true,
			"level":           "info",
			"format":          "json",
			"output_type":     "console",
			"service_name":    "detectviz",
			"service_version": "1.0.0",
			"environment":     "development",
		},
		Enabled: true,
	}

	// Register metadata if the registry supports it
	if registryManager, ok := registry.(interface {
		RegisterMetadata(string, *contracts.PluginMetadata) error
	}); ok {
		registryManager.RegisterMetadata("otelzap", metadata)
	}
}
