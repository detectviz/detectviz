package metric

import "context"

// MetricQueryAdapter defines the interface for querying metric values from various data sources.
// zh: MetricQueryAdapter 定義從各種資料來源查詢監控指標的介面。
type MetricQueryAdapter interface {
	// Query executes a query expression with optional label filters and returns a numeric result.
	// zh: 執行查詢語句並搭配可選的標籤篩選，回傳單一數值結果。
	QueryValue(ctx context.Context, expr string, labels map[string]string) (float64, error)
}
