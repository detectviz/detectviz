package promadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/stretchr/testify/assert"
)

func TestPromAlertEvaluator_Evaluate(t *testing.T) {
	e := NewEvaluator(testutil.NewTestLogger())

	tests := []struct {
		name        string
		condition   alert.AlertCondition
		wantFiring  bool
		wantMessage string
		wantErr     bool
	}{
		{
			name: "basic evaluation",
			condition: alert.AlertCondition{
				Expr:      "up{job=\"node\"} == 0",
				Threshold: 1,
				Labels:    map[string]string{"job": "node"},
			},
			wantFiring:  false,
			wantMessage: "prom evaluation not implemented",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Evaluate(context.Background(), tt.condition)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantFiring, result.Firing)
			assert.Equal(t, tt.wantMessage, result.Message)
		})
	}
}
