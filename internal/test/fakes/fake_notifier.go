package fakes

import (
	"context"

	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// FakeNotifier is a simple notifier implementation used for tests.
// zh: 測試用 Notifier 假實作，用於驗證通知流程。
type FakeNotifier struct {
	NameVal     string
	NotifyCalls []NotifyCall
	SendCalls   []SendCall
	NotifyErr   error
	SendErr     error
}

type NotifyCall struct {
	Title   string
	Message string
}

type SendCall struct {
	Ctx context.Context
	Msg notifieriface.Message
}

func (f *FakeNotifier) Name() string { return f.NameVal }

func (f *FakeNotifier) Notify(title, message string) error {
	f.NotifyCalls = append(f.NotifyCalls, NotifyCall{Title: title, Message: message})
	return f.NotifyErr
}

func (f *FakeNotifier) Send(ctx context.Context, msg notifieriface.Message) error {
	f.SendCalls = append(f.SendCalls, SendCall{Ctx: ctx, Msg: msg})
	return f.SendErr
}
