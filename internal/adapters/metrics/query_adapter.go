package metricsadapter

import (
	"context"
	"errors"
)

// FluxQueryAdapter implements MetricQueryAdapter for querying InfluxDB via Flux.
// zh: FluxQueryAdapter 是透過 Flux 查詢 InfluxDB 的 MetricQueryAdapter 實作。
type FluxQueryAdapter struct {
	// TODO: 注入 InfluxDB client 或查詢執行器
	// Client influxdb2.Client
}

func (f *FluxQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	// TODO: 根據 expr 與 labels 組裝 Flux 語法
	// TODO: 呼叫 InfluxDB 查詢 API 並解析回傳資料
	// TODO: 回傳查詢結果中的第一筆數值或適當的錯誤
	return 0, errors.New("flux query not implemented")
}

// PromQueryAdapter implements metric.MetricQueryAdapter for querying Prometheus using PromQL.
// PromQueryAdapter 是使用 PromQL 查詢 Prometheus 的 metric.MetricQueryAdapter 實作。
type PromQueryAdapter struct {
	// TODO: 注入 Prometheus 查詢 client，例如 prometheus.Client
	// Client *promapi.Client
}

func (p *PromQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	// TODO: 根據 expr 與 labels 組裝 PromQL 語法
	// TODO: 呼叫 Prometheus API 並解析回傳資料
	// TODO: 回傳查詢結果中的第一筆數值或適當的錯誤
	return 0, errors.New("prometheus query not implemented")
}

// MockQueryAdapter is a stub implementation of MetricQueryAdapter for testing.
// zh: MockQueryAdapter 是 MetricQueryAdapter 的測試用假實作，回傳固定值。
type MockQueryAdapter struct {
	FixedValue float64
	Err        error
}

func (m *MockQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return m.FixedValue, nil
}
