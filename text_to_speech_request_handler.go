package main

import (
	qdb "github.com/rqure/qdb/src"
)

type TextToSpeechRequestHandlerSignals struct {
	NewRequest qdb.Signal
}

type TextToSpeechRequestHandler struct {
	db              qdb.IDatabase
	isLeader        bool
	subscriptionIds []string

	Signals TextToSpeechRequestHandlerSignals
}

func NewTextToSpeechRequestHandler(db qdb.IDatabase) *TextToSpeechRequestHandler {
	return &TextToSpeechRequestHandler{
		db: db,
	}
}

func (w *TextToSpeechRequestHandler) OnBecameLeader() {
	w.isLeader = true

	w.db.Notify(&qdb.DatabaseNotificationConfig{
		Type:  "AudioController",
		Field: "TextToSpeech",
	}, w.ProcessNotification)
}

func (w *TextToSpeechRequestHandler) OnLostLeadership() {
	w.isLeader = false
	w.subscriptionIds = []string{}
}

func (w *TextToSpeechRequestHandler) Init() {

}

func (w *TextToSpeechRequestHandler) Deinit() {

}

func (w *TextToSpeechRequestHandler) DoWork() {

}

func (w *TextToSpeechRequestHandler) ProcessNotification(notification *qdb.DatabaseNotification) {
	if !w.isLeader {
		return
	}

	textToSpeechField, err := notification.Current.Value.UnmarshalNew()
	if err != nil {
		qdb.Error("[TextToSpeechRequestHandler::ProcessNotification] Failed to unmarshal text to speech field: %v", err)
		return
	}

	switch textToSpeech := textToSpeechField.(type) {
	case *qdb.String:
		qdb.Info("[TextToSpeechRequestHandler::ProcessNotification] Adding request to play text to speech: %s", textToSpeech.Raw)
		w.Signals.NewRequest.Emit(textToSpeech.Raw)
	}
}
