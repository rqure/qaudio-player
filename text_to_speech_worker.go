package main

import (
	qdb "github.com/rqure/qdb/src"
	"github.com/rqure/qtts"
	"github.com/rqure/qtts/voices"
)

type TextToSpeechWorker struct {
	db              qdb.IDatabase
	isLeader        bool
	subscriptionIds []string
	tts             *qtts.Speech
}

func NewTextToSpeechWorker(db qdb.IDatabase) *TextToSpeechWorker {
	return &TextToSpeechWorker{
		db: db,
		tts: &qtts.Speech{
			Folder:   "/",
			Language: voices.English,
			Handler:  audioPlayer},
	}
}

func (w *TextToSpeechWorker) OnBecameLeader() {
	w.isLeader = true

	w.db.Notify(&qdb.DatabaseNotificationConfig{
		Type:  "AudioController",
		Field: "TextToSpeech",
	}, w.ProcessNotification)
}

func (w *TextToSpeechWorker) OnLostLeadership() {
	w.isLeader = false
	w.subscriptionIds = []string{}
}

func (w *TextToSpeechWorker) Init() {

}

func (w *TextToSpeechWorker) Deinit() {

}

func (w *TextToSpeechWorker) DoWork() {

}

func (w *TextToSpeechWorker) ProcessNotification(notification *qdb.DatabaseNotification) {
	if !w.isLeader {
		return
	}

	textToSpeechField, err := notification.Current.Value.UnmarshalNew()
	if err != nil {
		qdb.Error("[TextToSpeechWorker::ProcessNotification] Failed to unmarshal text to speech field: %v", err)
		return
	}

	switch textToSpeech := textToSpeechField.(type) {
	case *qdb.String:
		err := e.Tts.Speak(textToSpeech.Raw)
		if err != nil {
			qdb.Error("[TextToSpeechWorker::ProcessNotification] Failed to play text to speech (%s): %v", textToSpeech.Raw, err)
		} else {
			qdb.Error("[TextToSpeechWorker::ProcessNotification] Played text to speech: %s", textToSpeech.Raw)
		}
	}
}
