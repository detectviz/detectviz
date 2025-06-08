// Package metric provides shared metric-related interfaces for the detectviz project.
// zh: 提供 Detectviz 專案中與指標資料相關的共用介面。

package metric

import (
	"context"
)

// MetricWriter defines the interface for sending metric data to external systems.
// zh: MetricWriter 定義寫入指標資料至外部系統的介面（例如 InfluxDB、Pushgateway）。
type MetricWriter interface {
	// WritePoint writes a single metric point with measurement name, value, and labels.
	// zh: 寫入單筆指標資料，包含量測名稱、數值與標籤。
	WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error
}

// MetricTransformer defines the interface for preprocessing metric data before evaluation or storage.
// zh: MetricTransformer 定義指標資料在評估或儲存前的預處理邏輯（例如單位轉換、標籤增補）。
type MetricTransformer interface {
	// Transform modifies the measurement name, value, and labels in-place before processing.
	// zh: 對指標名稱、數值與標籤進行轉換處理，會就地修改輸入值。
	Transform(measurement *string, value *float64, labels map[string]string) error
}

// MetricSeriesReader defines the interface for reading a time series of metric data.
// zh: MetricSeriesReader 定義讀取時間序列資料的介面，常用於報表、圖表或趨勢分析。
type MetricSeriesReader interface {
	// ReadSeries returns a list of timestamped values for a given expression and labels within a time range.
	// zh: 讀取指定表達式與標籤條件的時間序列資料，回傳時間戳與對應數值清單。
	ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]TimePoint, error)
}
