package alertadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
)

// MockAlertEvaluator 是 alert.AlertEvaluator 的模擬實作。
// zh: 用於測試的假告警判斷器，可自訂回傳結果與錯誤。
type MockAlertEvaluator struct {
	MockResult alert.AlertResult // zh: 預設回傳的結果
	MockError  error             // zh: 預設回傳的錯誤
}

// NewMockAlertEvaluator 建立一個 MockAlertEvaluator 實例。
// zh: 可傳入欲回傳的結果與錯誤，模擬告警行為。
func NewMockAlertEvaluator(result alert.AlertResult, err error) *MockAlertEvaluator {
	return &MockAlertEvaluator{
		MockResult: result,
		MockError:  err,
	}
}

// Evaluate 回傳預設的結果與錯誤。
// zh: 不執行任何邏輯，僅回傳初始化時指定的內容。
func (m *MockAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	return m.MockResult, m.MockError
}
