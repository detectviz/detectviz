package server

import (
	"context"
	"fmt"
	"net/http"
)

// Run starts all modules and the HTTP server.
// zh: 啟動所有模組並啟動 HTTP 伺服器。
func (s *Server) Run(ctx context.Context) error {
	// 啟動模組
	if err := s.ModuleEngine.RunAll(ctx); err != nil {
		return fmt.Errorf("failed to run modules: %w", err)
	}

	// 建立 HTTP Server
	s.HTTPServer = &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux, // TODO: 可替換為自定義 multiplexer
	}

	// 啟動 HTTP Server
	go func() {
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("HTTP server error", "error", err)
		}
	}()

	s.Logger.Info("Server started on :8080")
	<-ctx.Done()
	return nil
}

// Shutdown gracefully shuts down the HTTP server and all modules.
// zh: 優雅地關閉 HTTP 伺服器與所有模組。
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server...")

	if s.HTTPServer != nil {
		if err := s.HTTPServer.Shutdown(ctx); err != nil {
			s.Logger.Warn("HTTP server shutdown error", "error", err)
		}
	}

	if err := s.ModuleEngine.ShutdownAll(ctx); err != nil {
		return fmt.Errorf("failed to shutdown modules: %w", err)
	}

	s.Logger.Info("Server shutdown complete")
	return nil
}
