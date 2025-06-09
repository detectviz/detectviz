package metricsadapter

import (
	"context"
	"errors"
)

// InfluxMetricWriter implements MetricWriter for sending metrics to InfluxDB.
// zh: InfluxMetricWriter 是將指標寫入 InfluxDB 的 MetricWriter 實作。
type InfluxMetricWriter struct {
	// TODO: 注入 InfluxDB 寫入 client
	// Client influxdb2.WriteAPI
}

func (w *InfluxMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	// TODO: 將 measurement、value、labels 組裝成 Point 資料
	// TODO: 呼叫 InfluxDB client 寫入方法
	return errors.New("influx write not implemented")
}

// PushgatewayMetricWriter implements MetricWriter for pushing metrics to Prometheus Pushgateway.
// zh: PushgatewayMetricWriter 是將指標推送至 Prometheus Pushgateway 的 MetricWriter 實作。
type PushgatewayMetricWriter struct {
	// TODO: 注入 Pushgateway HTTP client 或設定參數
	// Endpoint string
}

func (w *PushgatewayMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	// TODO: 組裝成 Pushgateway 支援的格式並發送 HTTP 請求
	return errors.New("pushgateway write not implemented")
}

// MockMetricWriter is a mock implementation of MetricWriter for testing.
// zh: MockMetricWriter 是測試用的 MetricWriter 實作，可記錄寫入行為或模擬錯誤。
type MockMetricWriter struct {
	LastMeasurement string
	LastValue       float64
	LastLabels      map[string]string
	Err             error
}

func (m *MockMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	if m.Err != nil {
		return m.Err
	}
	m.LastMeasurement = measurement
	m.LastValue = value
	m.LastLabels = labels
	return nil
}
