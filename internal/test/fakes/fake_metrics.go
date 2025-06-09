package fakes

import (
	"context"

	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// FakeMetricWriter is a fake implementation of metric.MetricWriter for testing.
// zh: FakeMetricWriter 是測試用的 MetricWriter 假實作。
type FakeMetricWriter struct {
	WritePointCalls []WritePointCall
	WritePointError error
}

type WritePointCall struct {
	Ctx         context.Context
	Measurement string
	Value       float64
	Labels      map[string]string
}

func (f *FakeMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	f.WritePointCalls = append(f.WritePointCalls, WritePointCall{
		Ctx:         ctx,
		Measurement: measurement,
		Value:       value,
		Labels:      labels,
	})
	return f.WritePointError
}

// FakeMetricTransformer is a fake implementation of metric.MetricTransformer for testing.
// zh: FakeMetricTransformer 是測試用的 MetricTransformer 假實作。
type FakeMetricTransformer struct {
	TransformCalls []TransformCall
	TransformError error
}

type TransformCall struct {
	Measurement string
	Value       float64
	Labels      map[string]string
}

func (f *FakeMetricTransformer) Transform(measurement *string, value *float64, labels map[string]string) error {
	f.TransformCalls = append(f.TransformCalls, TransformCall{
		Measurement: *measurement,
		Value:       *value,
		Labels:      labels,
	})
	return f.TransformError
}

// FakeMetricSeriesReader is a fake implementation of metric.MetricSeriesReader for testing.
// zh: FakeMetricSeriesReader 是測試用的 MetricSeriesReader 假實作。
type FakeMetricSeriesReader struct {
	ReadSeriesCalls  []ReadSeriesCall
	ReadSeriesResult []metric.TimePoint
	ReadSeriesError  error
}

type ReadSeriesCall struct {
	Ctx    context.Context
	Expr   string
	Labels map[string]string
	Start  int64
	End    int64
}

func (f *FakeMetricSeriesReader) ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]metric.TimePoint, error) {
	f.ReadSeriesCalls = append(f.ReadSeriesCalls, ReadSeriesCall{
		Ctx:    ctx,
		Expr:   expr,
		Labels: labels,
		Start:  start,
		End:    end,
	})
	return f.ReadSeriesResult, f.ReadSeriesError
}
