package main

import (
	"context"

	"github.com/rqure/qlib/pkg/app"
	"github.com/rqure/qlib/pkg/data"
	"github.com/rqure/qlib/pkg/data/notification"
	"github.com/rqure/qlib/pkg/log"
	"github.com/rqure/qlib/pkg/signalslots"
)

type AudioFileRequestHandler struct {
	store              data.Store
	isLeader           bool
	notificationTokens []data.NotificationToken

	NewRequest signalslots.Signal
}

func NewAudioFileRequestHandler(store data.Store) *AudioFileRequestHandler {
	return &AudioFileRequestHandler{
		store: store,
	}
}

func (w *AudioFileRequestHandler) OnBecameLeader(ctx context.Context) {
	w.isLeader = true

	w.notificationTokens = append(w.notificationTokens, w.store.Notify(
		ctx,
		notification.NewConfig().
			SetEntityType("AudioController").
			SetFieldName("AudioFile").
			SetContextFields([]string{
				"AudioFile->Description",
				"AudioFile->Content",
			},
			),
		notification.NewCallback(w.ProcessNotification),
	))
}

func (w *AudioFileRequestHandler) OnLostLeadership(ctx context.Context) {
	w.isLeader = false

	for _, token := range w.notificationTokens {
		token.Unbind(ctx)
	}

	w.notificationTokens = []data.NotificationToken{}
}

func (w *AudioFileRequestHandler) Init(context.Context, app.Handle) {

}

func (w *AudioFileRequestHandler) Deinit(context.Context) {

}

func (w *AudioFileRequestHandler) DoWork(context.Context) {

}

func (w *AudioFileRequestHandler) ProcessNotification(ctx context.Context, n data.Notification) {
	if !w.isLeader {
		return
	}

	log.Info("Received audio file request: %v", n)

	if n.GetContextCount() < 2 {
		w.NewRequest.Emit(ctx, "")
		return
	}

	description := n.GetContext(0).GetValue().GetString()
	content := n.GetContext(1).GetValue().GetBinaryFile()

	log.Info("Adding request play audio file: %s", description)

	w.NewRequest.Emit(ctx, content)
}
