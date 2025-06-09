package metricsadapter

import (
	"errors"

	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// AggregationType defines the available aggregation methods.
// zh: AggregationType 定義可用的統計聚合類型。
type AggregationType string

const (
	SumAggregation     AggregationType = "sum"
	AverageAggregation AggregationType = "avg"
	MaxAggregation     AggregationType = "max"
	MinAggregation     AggregationType = "min"
)

// SimpleAggregator provides basic aggregation logic over time series.
// zh: SimpleAggregator 提供對時間序列資料的基本統計運算邏輯。
type SimpleAggregator struct{}

// Aggregate performs the specified aggregation on a list of TimePoints.
// zh: 根據指定類型對一組時間點進行聚合計算。
func (a *SimpleAggregator) Aggregate(points []metric.TimePoint, aggType AggregationType) (float64, error) {
	if len(points) == 0 {
		return 0, errors.New("no points to aggregate")
	}

	switch aggType {
	case SumAggregation:
		var sum float64
		for _, p := range points {
			sum += p.Value
		}
		return sum, nil

	case AverageAggregation:
		var sum float64
		for _, p := range points {
			sum += p.Value
		}
		return sum / float64(len(points)), nil

	case MaxAggregation:
		max := points[0].Value
		for _, p := range points {
			if p.Value > max {
				max = p.Value
			}
		}
		return max, nil

	case MinAggregation:
		min := points[0].Value
		for _, p := range points {
			if p.Value < min {
				min = p.Value
			}
		}
		return min, nil

	default:
		return 0, errors.New("unsupported aggregation type")
	}
}
