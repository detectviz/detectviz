// Package alert 提供 AlertEvaluator 的預設實作，用於根據規則比對指標狀態。
package alertadapter

import (
	"context"
	"fmt"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// DefaultAlertEvaluator 是 alert.AlertEvaluator 的預設實作。
// zh: 預設告警判斷器，用於依據查詢結果與規則閾值進行比對。
type DefaultAlertEvaluator struct {
	Log   logger.Logger             // zh: 用於記錄評估過程與錯誤資訊
	Query metric.MetricQueryAdapter // zh: 指標查詢介面（需由外部注入）
}

// NewDefaultAlertEvaluator 建立一個 DefaultAlertEvaluator 實例。
// zh: 可傳入 logger 與 metric.QueryAdapter 元件以追蹤告警觸發過程並查詢指標。
func NewDefaultAlertEvaluator(log logger.Logger, query metric.MetricQueryAdapter) *DefaultAlertEvaluator {
	return &DefaultAlertEvaluator{
		Log:   log.Named("alertevaluator"),
		Query: query,
	}
}

// Evaluate 根據告警條件進行比對。
// zh: 根據 AlertCondition 查詢指標結果，並依據閾值與運算子進行告警判斷。
func (e *DefaultAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	e.Log.Debug("evaluating alert condition", "expr", cond.Expr, "threshold", cond.Threshold)

	// 查詢指標資料
	value, err := e.Query.Query(ctx, cond.Expr, cond.Labels)
	if err != nil {
		e.Log.Error("query failed", "error", err)
		return alert.AlertResult{
			Firing:  false,
			Message: fmt.Sprintf("failed to evaluate '%s': %v", cond.Expr, err),
			Value:   0,
		}, err
	}

	e.Log.Debug("evaluation result", "expr", cond.Expr, "value", value, "threshold", cond.Threshold)

	// 呼叫通用比對邏輯
	result, evalErr := alert.Evaluate(value, cond)
	return result, evalErr
}
