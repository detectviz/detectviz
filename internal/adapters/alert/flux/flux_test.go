package fluxadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/stretchr/testify/assert"
)

func TestFluxAlertEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		logger logger.Logger
	}
	type args struct {
		ctx  context.Context
		cond alert.AlertCondition
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFiring bool
		wantMsg    string
		wantErr    bool
	}{
		{
			name: "not implemented",
			fields: fields{
				logger: testutil.NewTestLogger(),
			},
			args: args{
				ctx: context.Background(),
				cond: alert.AlertCondition{
					Expr:      `from(bucket:"test") |> range(start:-1h)`,
					Threshold: 80,
					Labels:    map[string]string{"host": "dev"},
				},
			},
			wantFiring: false,
			wantMsg:    "flux evaluation not implemented",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvaluator(tt.fields.logger)
			result, err := e.Evaluate(tt.args.ctx, tt.args.cond)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantFiring, result.Firing)
			assert.Equal(t, tt.wantMsg, result.Message)
		})
	}
}
