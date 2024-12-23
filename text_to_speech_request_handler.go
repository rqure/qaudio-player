package main

import (
	"context"

	"github.com/rqure/qlib/pkg/app"
	"github.com/rqure/qlib/pkg/data"
	"github.com/rqure/qlib/pkg/data/notification"
	"github.com/rqure/qlib/pkg/signalslots"
	"github.com/rqure/qlib/pkg/signalslots/signal"
)

type TextToSpeechRequestHandler struct {
	store              data.Store
	isLeader           bool
	notificationTokens []data.NotificationToken

	NewRequest signalslots.Signal
}

func NewTextToSpeechRequestHandler(store data.Store) *TextToSpeechRequestHandler {
	return &TextToSpeechRequestHandler{
		store:      store,
		NewRequest: signal.New(),
	}
}

func (w *TextToSpeechRequestHandler) OnBecameLeader(ctx context.Context) {
	w.isLeader = true

	w.notificationTokens = append(w.notificationTokens, w.store.Notify(
		ctx,
		notification.NewConfig().
			SetEntityType("AudioController").
			SetFieldName("TextToSpeech").
			SetContextFields("TTSLanguage", "TTSGender"),
		notification.NewCallback(w.ProcessNotification)))
}

func (w *TextToSpeechRequestHandler) OnLostLeadership(ctx context.Context) {
	w.isLeader = false

	for _, token := range w.notificationTokens {
		token.Unbind(ctx)
	}

	w.notificationTokens = []data.NotificationToken{}
}

func (w *TextToSpeechRequestHandler) Init(context.Context, app.Handle) {

}

func (w *TextToSpeechRequestHandler) Deinit(context.Context) {

}

func (w *TextToSpeechRequestHandler) DoWork(context.Context) {

}

func (w *TextToSpeechRequestHandler) ProcessNotification(ctx context.Context, n data.Notification) {
	if !w.isLeader {
		return
	}

	text := n.GetCurrent().GetValue().GetString()
	lang := n.GetContext(0).GetValue().GetString()
	gender := n.GetContext(1).GetValue().GetString()
	w.NewRequest.Emit(ctx, text, lang, gender)
}
