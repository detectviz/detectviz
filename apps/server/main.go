package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"detectviz/internal/platform/composition"
	"detectviz/internal/platform/registry"
	"detectviz/pkg/config/loader"
	"detectviz/pkg/platform/contracts"

	// Import plugins to register them
	prometheusPlugin "detectviz/plugins/community/importers/prometheus"
	jwtPlugin "detectviz/plugins/core/auth/jwt"
)

const (
	defaultCompositionPath = "../../compositions/minimal-platform/composition.yaml"
	defaultTimeout         = 30 * time.Second
)

// Server represents the main server application
// zh: Server 代表主要的伺服器應用程式
type Server struct {
	registry         *registry.Manager
	lifecycleManager *composition.LifecycleManager
	configLoader     *loader.ConfigLoader
	resolver         contracts.DependencyResolver
	ctx              context.Context
	cancel           context.CancelFunc
}

// NewServer creates a new server instance
// zh: NewServer 建立新的伺服器實例
func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		registry:     registry.NewManager(),
		resolver:     composition.NewDependencyResolver(),
		configLoader: loader.NewConfigLoader(""), // Use empty base path to avoid path joining issues
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Initialize initializes the server components
// zh: Initialize 初始化伺服器組件
func (s *Server) Initialize() error {
	log.Println("Initializing DetectViz minimal server...")

	// Register plugins first
	log.Println("Registering plugins...")
	if err := s.registerPlugins(); err != nil {
		return fmt.Errorf("failed to register plugins: %w", err)
	}

	// Load composition configuration
	compositionPath := os.Getenv("DETECTVIZ_COMPOSITION")
	if compositionPath == "" {
		// Try different possible paths
		possiblePaths := []string{
			"compositions/minimal-platform/composition.yaml",
			"../../compositions/minimal-platform/composition.yaml",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				compositionPath = path
				break
			}
		}

		if compositionPath == "" {
			return fmt.Errorf("cannot find composition file in any of the expected locations")
		}
	}

	log.Printf("Loading composition from: %s", compositionPath)
	config, err := s.configLoader.LoadComposition(compositionPath)
	if err != nil {
		return fmt.Errorf("failed to load composition: %w", err)
	}

	// Validate configuration
	if err := s.configLoader.ValidateConfiguration(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Loaded composition: %s v%s", config.Metadata.Name, config.Metadata.Version)

	// Get enabled plugins
	pluginConfigs, err := s.configLoader.GetPluginConfigs()
	if err != nil {
		return fmt.Errorf("failed to get plugin configs: %w", err)
	}

	log.Printf("Found %d enabled plugins", len(pluginConfigs))

	// Register plugin metadata
	for _, pluginConfig := range pluginConfigs {
		if err := s.registry.RegisterMetadata(pluginConfig.Name, &pluginConfig); err != nil {
			return fmt.Errorf("failed to register metadata for plugin %s: %w", pluginConfig.Name, err)
		}
		log.Printf("Registered plugin metadata: %s (type: %s, category: %s)",
			pluginConfig.Name, pluginConfig.Type, pluginConfig.Category)
	}

	// Initialize lifecycle manager
	s.lifecycleManager = composition.NewLifecycleManager(s.resolver)

	// Initialize lifecycle manager
	if err := s.lifecycleManager.Initialize(s.ctx); err != nil {
		return fmt.Errorf("failed to initialize lifecycle manager: %w", err)
	}

	log.Println("Server initialization completed successfully")
	return nil
}

// registerPlugins registers all available plugins
// zh: registerPlugins 註冊所有可用的插件
func (s *Server) registerPlugins() error {
	// Register JWT plugin
	if err := jwtPlugin.Register(s.registry); err != nil {
		return fmt.Errorf("failed to register JWT plugin: %w", err)
	}
	log.Println("Registered JWT authenticator plugin")

	// Register Prometheus plugin
	if err := prometheusPlugin.Register(s.registry); err != nil {
		return fmt.Errorf("failed to register Prometheus plugin: %w", err)
	}
	log.Println("Registered Prometheus importer plugin")

	return nil
}

// Start starts the server
// zh: Start 啟動伺服器
func (s *Server) Start() error {
	log.Println("Starting DetectViz server...")

	// Start lifecycle manager with all plugins
	if err := s.lifecycleManager.StartAll(s.ctx, s.registry); err != nil {
		return fmt.Errorf("failed to start plugins: %w", err)
	}

	log.Println("All plugins started successfully")

	// Perform initial health check
	healthStatus := s.lifecycleManager.HealthCheck(s.ctx, s.registry)
	log.Printf("Initial health check completed:")
	for pluginName, status := range healthStatus {
		log.Printf("  - %s: %s (%s)", pluginName, status.Status, status.Message)
	}

	log.Println("DetectViz server started successfully")
	return nil
}

// Stop stops the server
// zh: Stop 停止伺服器
func (s *Server) Stop() error {
	log.Println("Stopping DetectViz server...")

	// Cancel context to signal shutdown
	s.cancel()

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Shutdown all plugins
	if err := s.lifecycleManager.ShutdownAll(shutdownCtx, s.registry); err != nil {
		log.Printf("Error during shutdown: %v", err)
		return err
	}

	log.Println("DetectViz server stopped successfully")
	return nil
}

// Run runs the server until termination signal
// zh: Run 運行伺服器直到終止信號
func (s *Server) Run() error {
	// Initialize server
	if err := s.Initialize(); err != nil {
		return fmt.Errorf("server initialization failed: %w", err)
	}

	// Start server
	if err := s.Start(); err != nil {
		return fmt.Errorf("server start failed: %w", err)
	}

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run health check periodically
	healthTicker := time.NewTicker(30 * time.Second)
	defer healthTicker.Stop()

	for {
		select {
		case <-sigChan:
			log.Println("Received termination signal")
			return s.Stop()

		case <-healthTicker.C:
			// Perform periodic health check
			healthStatus := s.lifecycleManager.HealthCheck(s.ctx, s.registry)
			healthyCount := 0
			for _, status := range healthStatus {
				if status.Status == "healthy" {
					healthyCount++
				}
			}
			log.Printf("Health check: %d/%d plugins healthy", healthyCount, len(healthStatus))

		case <-s.ctx.Done():
			log.Println("Server context cancelled")
			return nil
		}
	}
}

// main is the entry point for the server
// zh: main 是伺服器的入口點
func main() {
	log.Println("Starting DetectViz Minimal Server...")

	// Create server instance
	server := NewServer()

	// Run server
	if err := server.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server exited cleanly")
}
