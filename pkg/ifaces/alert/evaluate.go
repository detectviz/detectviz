package alert

import (
	"fmt"
)

// Evaluate 判斷查詢結果是否觸發告警條件。
// zh: 根據 AlertCondition 中指定的運算子（Operator）與閾值（Threshold），比對查詢結果 value 是否觸發告警。
func Evaluate(value float64, cond AlertCondition) (AlertResult, error) {
	operator := cond.Operator
	if operator == "" {
		operator = "ge" // 預設為 >=
	}

	var firing bool
	var message string

	switch operator {
	case "ge":
		firing = value >= cond.Threshold
		message = "threshold exceeded (>=)"
	case "gt":
		firing = value > cond.Threshold
		message = "threshold exceeded (>)"
	case "le":
		firing = value <= cond.Threshold
		message = "threshold under (<=)"
	case "lt":
		firing = value < cond.Threshold
		message = "threshold under (<)"
	case "eq":
		firing = value == cond.Threshold
		message = "threshold matched (==)"
	case "ne":
		firing = value != cond.Threshold
		message = "threshold not matched (!=)"
	default:
		return AlertResult{
			Firing:  false,
			Message: fmt.Sprintf("unsupported operator: %s", cond.Operator),
			Value:   value,
		}, fmt.Errorf("unsupported operator: %s", cond.Operator)
	}

	return AlertResult{
		Firing:  firing,
		Message: message,
		Value:   value,
	}, nil
}
