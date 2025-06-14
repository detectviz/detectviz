package contracts

import (
	"context"
	"time"
)

// Importer defines the interface for data importers (Telegraf Input pattern).
// zh: Importer 定義資料匯入器介面（遵循 Telegraf Input 模式）。
type Importer interface {
	Plugin
	Import(ctx context.Context) error
}

// StreamingImporter defines the interface for continuous data importers.
// zh: StreamingImporter 定義持續資料匯入器介面。
type StreamingImporter interface {
	Importer
	StartStreaming(ctx context.Context) (<-chan ImportData, error)
	StopStreaming() error
}

// BatchImporter defines the interface for batch data importers.
// zh: BatchImporter 定義批次資料匯入器介面。
type BatchImporter interface {
	Importer
	ImportBatch(ctx context.Context, batchSize int) ([]ImportData, error)
}

// ImportData represents the standard data structure for imported data.
// zh: ImportData 代表匯入資料的標準資料結構。
type ImportData struct {
	Type      string                 `json:"type"`   // metrics, logs, traces, alerts
	Source    string                 `json:"source"` // 來源 plugin
	Timestamp time.Time              `json:"timestamp"`
	Labels    map[string]string      `json:"labels"`
	Fields    map[string]interface{} `json:"fields"`
	Raw       []byte                 `json:"raw,omitempty"`
}

// ImporterConfig defines the configuration structure for importers.
// zh: ImporterConfig 定義匯入器的配置結構。
type ImporterConfig struct {
	Name        string            `yaml:"name" json:"name"`
	Type        string            `yaml:"type" json:"type"`
	Endpoint    string            `yaml:"endpoint" json:"endpoint"`
	Interval    time.Duration     `yaml:"interval" json:"interval"`
	Timeout     time.Duration     `yaml:"timeout" json:"timeout"`
	BatchSize   int               `yaml:"batch_size" json:"batch_size"`
	Credentials map[string]string `yaml:"credentials" json:"credentials"`
	ExtraConfig map[string]any    `yaml:"extra_config" json:"extra_config"`
}

// ImporterRegistry defines the interface for importer registration.
// zh: ImporterRegistry 定義匯入器註冊介面。
type ImporterRegistry interface {
	RegisterImporter(name string, factory func(config ImporterConfig) (Importer, error)) error
	GetImporter(name string) (Importer, error)
	ListImporters() []string
}
