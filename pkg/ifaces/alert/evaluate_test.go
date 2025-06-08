package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		cond     AlertCondition
		wantFire bool
		wantMsg  string
		wantErr  bool
	}{
		{"default operator (ge)", 80, AlertCondition{Threshold: 70}, true, "threshold exceeded (>=)", false},
		{"gt operator", 75, AlertCondition{Threshold: 70, Operator: "gt"}, true, "threshold exceeded (>)", false},
		{"lt operator", 60, AlertCondition{Threshold: 70, Operator: "lt"}, true, "threshold under (<)", false},
		{"eq operator", 100, AlertCondition{Threshold: 100, Operator: "eq"}, true, "threshold matched (==)", false},
		{"ne operator", 95, AlertCondition{Threshold: 100, Operator: "ne"}, true, "threshold not matched (!=)", false},
		{"unsupported operator", 42, AlertCondition{Threshold: 50, Operator: "bad"}, false, "unsupported operator: bad", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.value, tt.cond)
			assert.Equal(t, tt.wantFire, result.Firing)
			assert.Contains(t, result.Message, tt.wantMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
