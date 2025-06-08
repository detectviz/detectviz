package alert

import (
	"context"
	"errors"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/stretchr/testify/assert"
)

// TestMockAlertEvaluator_ReturnsExpectedResult 驗證模擬評估器能正確回傳指定結果。
// zh: 模擬告警條件成功觸發時的行為。
func TestMockAlertEvaluator_ReturnsExpectedResult(t *testing.T) {
	want := alert.AlertResult{
		Firing:  true,
		Message: "mock triggered",
		Value:   42.0,
	}

	mock := NewMockAlertEvaluator(want, nil)

	got, err := mock.Evaluate(context.Background(), alert.AlertCondition{
		Expr:      "mock_expr",
		Threshold: 40,
	})

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

// TestMockAlertEvaluator_ReturnsError 驗證模擬評估器在預期錯誤時能正確回傳 error。
// zh: 模擬告警評估過程中出現錯誤時的行為。
func TestMockAlertEvaluator_ReturnsError(t *testing.T) {
	mockErr := errors.New("mock error")
	mock := NewMockAlertEvaluator(alert.AlertResult{}, mockErr)

	_, err := mock.Evaluate(context.Background(), alert.AlertCondition{})

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
}
