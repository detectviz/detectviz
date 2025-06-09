package metricsadapter

import (
	"context"
	"errors"
	"time"

	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// MockSeriesReader implements MetricSeriesReader with dummy time series data.
// zh: MockSeriesReader 是 MetricSeriesReader 的模擬實作，用於回傳虛擬的時間序列資料。
type MockSeriesReader struct {
	Points []metric.TimePoint // zh: 預設回傳的資料點清單
	Err    error              // zh: 若設定錯誤，則每次查詢都會回傳此錯誤
}

// ReadSeries returns a sequence of metric points based on the configured mock data.
// zh: 根據預設資料回傳一段時間序列資料，若未設定則產生 5 筆模擬資料。
func (m *MockSeriesReader) ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]metric.TimePoint, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if len(m.Points) > 0 {
		return m.Points, nil
	}
	// Default mock: generate 5 points, spaced 1 minute apart
	var points []metric.TimePoint
	t := time.Now().Unix()
	for i := 0; i < 5; i++ {
		points = append(points, metric.TimePoint{
			Timestamp: t - int64(60*i),
			Value:     float64(i),
		})
	}
	return points, nil
}

// InfluxSeriesReader implements MetricSeriesReader for reading time series from InfluxDB.
// zh: InfluxSeriesReader 是用於從 InfluxDB 讀取時間序列資料的實作。
type InfluxSeriesReader struct {
	// TODO: 注入 InfluxDB 的查詢 client，例如 influxdb2.QueryAPI
	// Client influxdb2.QueryAPI
}

// ReadSeries reads a sequence of metric points from InfluxDB using Flux.
// zh: 從 InfluxDB 讀取指定條件的時間序列資料，並轉換為 TimePoint 格式回傳。
func (r *InfluxSeriesReader) ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]metric.TimePoint, error) {
	// TODO: 組裝 flux 查詢語法
	// TODO: 執行查詢，解析查詢結果為 TimePoint 陣列
	return nil, errors.New("influx series reader not implemented")
}
