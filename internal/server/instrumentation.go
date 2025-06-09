package server

import (
	"net/http"
	"net/http/pprof"
)

// RegisterInstrumentation mounts instrumentation routes into the given ServeMux.
// zh: 將監控與除錯端點註冊至 HTTP multiplexer。
func RegisterInstrumentation(mux *http.ServeMux) {
	// Prometheus metrics endpoint (佔位用，可替換為 promhttp.Handler())
	mux.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# HELP detectviz_metrics_placeholder Placeholder for Prometheus metrics\n"))
		w.Write([]byte("# TYPE detectviz_metrics_placeholder counter\n"))
		w.Write([]byte("detectviz_metrics_placeholder 1\n"))
	}))

	// Basic health check endpoint
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// Register pprof endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
