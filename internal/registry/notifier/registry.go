package notifier

import (
	"errors"
	"fmt"
	"net/http"

	notifieradapter "github.com/detectviz/detectviz/internal/adapters/notifier"
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// NewNotifierRegistry 建立 notifier 實例清單，根據 config 載入。
// zh: 根據設定檔產生可用的 notifier 實作清單。
func NewNotifierRegistry(cfgs []configtypes.NotifierConfig, log logger.Logger) []notifieriface.Notifier {
	var list []notifieriface.Notifier

	for _, cfg := range cfgs {
		if !cfg.Enable {
			continue
		}
		n, err := buildNotifier(cfg, log)
		if err != nil {
			log.Warn("Failed to build notifier", "type", cfg.Type, "name", cfg.Name, "err", err)
			continue
		}
		list = append(list, n)
	}

	return list
}

// buildNotifier 根據設定建立對應的 Notifier 實例。
// zh: 根據 notifier 類型回傳實作，支援 email, slack, webhook。
func buildNotifier(cfg configtypes.NotifierConfig, log logger.Logger) (notifieriface.Notifier, error) {
	switch cfg.Type {
	case "email":
		return notifieradapter.NewEmailNotifier(cfg.Name, cfg.Target, log), nil
	case "slack":
		return notifieradapter.NewSlackNotifier(cfg.Name, cfg.Target, log), nil
	case "webhook":
		return notifieradapter.NewWebhookNotifier(cfg.Name, log, http.DefaultClient), nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown notifier type: %s", cfg.Type))
	}
}

// ProvideNotifier 根據設定與 logger 提供整合後的 Notifier 實例。
// zh: 整合所有 notifier 為一個 Notifier，用於統一發送通知。
func ProvideNotifier(cfgs []configtypes.NotifierConfig, log logger.Logger) notifieriface.Notifier {
	list := NewNotifierRegistry(cfgs, log)
	if len(list) == 0 {
		return notifieradapter.NewNopNotifier()
	}
	return notifieradapter.NewMultiNotifier(list...)
}
