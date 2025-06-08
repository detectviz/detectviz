package metric

// TimePoint represents a single point in a time series.
// zh: TimePoint 表示時間序列中的一個時間點。
type TimePoint struct {
	Timestamp int64
	Value     float64
}

// TimeRange represents a start and end timestamp for querying time series data.
// zh: TimeRange 表示查詢時間序列資料時使用的起始與結束時間範圍。
type TimeRange struct {
	Start int64
	End   int64
}
