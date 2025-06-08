package alert

import "context"

// AlertCondition represents the input to the evaluator, typically a rule or threshold.
// zh: AlertCondition 表示要評估的告警條件，例如規則或閾值。
type AlertCondition struct {
	RuleID    string
	Expr      string
	Threshold float64
	Labels    map[string]string
}

// AlertResult represents the outcome of the evaluation.
// zh: AlertResult 表示評估結果，包含是否觸發及原因。
type AlertResult struct {
	Firing  bool
	Message string
	Value   float64
}

// AlertEvaluator defines the interface for evaluating alert conditions.
// zh: AlertEvaluator 定義告警條件評估器的介面，用於根據輸入條件判斷是否觸發告警。
type AlertEvaluator interface {
	// Evaluate analyzes a condition and returns the result.
	// zh: 根據輸入條件進行評估，回傳是否觸發與原因。
	Evaluate(ctx context.Context, cond AlertCondition) (AlertResult, error)
}
