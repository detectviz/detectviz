package alertlog

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

func TestAlertPluginHandler_HandleAlertTriggered(t *testing.T) {
	handler := &AlertPluginHandler{}

	t.Run("handle normal alert", func(t *testing.T) {
		event := event.AlertTriggeredEvent{
			AlertID:    "test-alert",
			RuleName:   "test-rule",
			Level:      "critical",
			Instance:   "test-instance",
			Metric:     "test-metric",
			Comparison: ">",
			Value:      100,
			Threshold:  100,
			Message:    "alert test",
		}

		err := handler.HandleAlertTriggered(context.Background(), event)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("handle empty fields", func(t *testing.T) {
		event := event.AlertTriggeredEvent{
			AlertID:    "",
			RuleName:   "",
			Level:      "",
			Instance:   "",
			Metric:     "",
			Comparison: "",
			Value:      0,
			Threshold:  0,
			Message:    "",
		}

		err := handler.HandleAlertTriggered(context.Background(), event)
		if err != nil {
			t.Errorf("expected no error for empty fields, got %v", err)
		}
	})
}
