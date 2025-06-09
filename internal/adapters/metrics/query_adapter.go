package metricsadapter

import (
	"context"
	"errors"
)

// FluxQueryAdapter implements metric.MetricQueryAdapter for querying InfluxDB via Flux.
// zh: FluxQueryAdapter 是透過 Flux 查詢 InfluxDB 的 MetricQueryAdapter 實作。
type FluxQueryAdapter struct {
	// TODO: Inject InfluxDB client or query executor.
	// zh: 注入 InfluxDB 查詢用的 client。
}

func (f *FluxQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	// TODO: Assemble Flux query from expr and labels.
	// zh: 根據 expr 與 labels 組裝 Flux 語法。

	// TODO: Execute InfluxDB query and parse response.

	// TODO: Return the first value or an appropriate error.
	return 0, errors.New("flux query not implemented")
}

// PromQueryAdapter implements metric.MetricQueryAdapter for querying Prometheus using PromQL.
// zh: PromQueryAdapter 是使用 PromQL 查詢 Prometheus 的 MetricQueryAdapter 實作。
type PromQueryAdapter struct {
	// TODO: Inject Prometheus query client, e.g., prometheus.Client.
	// zh: 注入 Prometheus 查詢用的 client。
}

func (p *PromQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	// TODO: Assemble PromQL query from expr and labels.
	// zh: 根據 expr 與 labels 組裝 PromQL 語法。

	// TODO: Execute Prometheus API query and parse response.

	// TODO: Return the first value or an appropriate error.
	return 0, errors.New("prometheus query not implemented")
}

// MockQueryAdapter is a stub implementation of MetricQueryAdapter for testing.
// zh: MockQueryAdapter 是 MetricQueryAdapter 的測試用假實作，回傳固定值。
type MockQueryAdapter struct {
	FixedValue float64

	Err error
}

func (m *MockQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	if m.Err != nil {
		return 0, m.Err
	}

	return m.FixedValue, nil
}
