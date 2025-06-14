package contracts

import (
	"context"
	"time"
)

// Exporter defines the interface for data exporters (Telegraf Output pattern).
// zh: Exporter 定義資料匯出器介面（遵循 Telegraf Output 模式）。
type Exporter interface {
	Plugin
	Export(ctx context.Context, data any) error
}

// StreamingExporter defines the interface for continuous data exporters.
// zh: StreamingExporter 定義持續資料匯出器介面。
type StreamingExporter interface {
	Exporter
	StartStreaming(ctx context.Context) (chan<- ExportData, error)
	StopStreaming() error
}

// BatchExporter defines the interface for batch data exporters.
// zh: BatchExporter 定義批次資料匯出器介面。
type BatchExporter interface {
	Exporter
	ExportBatch(ctx context.Context, data []ExportData) error
	Flush(ctx context.Context) error
}

// ExportData represents the standard data structure for exported data.
// zh: ExportData 代表匯出資料的標準資料結構。
type ExportData struct {
	Type      string                 `json:"type"`   // metrics, logs, traces, alerts
	Target    string                 `json:"target"` // 目標系統
	Timestamp time.Time              `json:"timestamp"`
	Labels    map[string]string      `json:"labels"`
	Fields    map[string]interface{} `json:"fields"`
	Raw       []byte                 `json:"raw,omitempty"`
}

// ExporterConfig defines the configuration structure for exporters.
// zh: ExporterConfig 定義匯出器的配置結構。
type ExporterConfig struct {
	Name         string            `yaml:"name" json:"name"`
	Type         string            `yaml:"type" json:"type"`
	Endpoint     string            `yaml:"endpoint" json:"endpoint"`
	BatchSize    int               `yaml:"batch_size" json:"batch_size"`
	FlushTimeout time.Duration     `yaml:"flush_timeout" json:"flush_timeout"`
	RetryCount   int               `yaml:"retry_count" json:"retry_count"`
	Credentials  map[string]string `yaml:"credentials" json:"credentials"`
	ExtraConfig  map[string]any    `yaml:"extra_config" json:"extra_config"`
}

// ExporterRegistry defines the interface for exporter registration.
// zh: ExporterRegistry 定義匯出器註冊介面。
type ExporterRegistry interface {
	RegisterExporter(name string, factory func(config ExporterConfig) (Exporter, error)) error
	GetExporter(name string) (Exporter, error)
	ListExporters() []string
}

// DataPipeline defines the interface for data processing pipeline.
// zh: DataPipeline 定義資料處理管道介面。
type DataPipeline interface {
	Send(data ExportData) error
	Receive() <-chan ExportData
	Transform(transformer DataTransformer) DataPipeline
}

// DataTransformer defines the interface for data transformation.
// zh: DataTransformer 定義資料轉換介面。
type DataTransformer interface {
	Transform(data ExportData) (ExportData, error)
}
