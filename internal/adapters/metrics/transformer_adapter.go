package metricsadapter

// NoopTransformer implements MetricTransformer without any modification.
// zh: NoopTransformer 是不進行任何轉換的預設實作。
type NoopTransformer struct{}

// Transform fulfills the MetricTransformer interface without changing input.
// zh: 不進行任何修改，原樣傳回
func (t *NoopTransformer) Transform(measurement *string, value *float64, labels map[string]string) error {
	return nil
}
