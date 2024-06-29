package main

import (
	qdb "github.com/rqure/qdb/src"
)

type AudioFileRequestHandlerSignals struct {
	NewRequest qdb.Signal
}

type AudioFileRequestHandler struct {
	db              qdb.IDatabase
	isLeader        bool
	subscriptionIds []string

	Signals AudioFileRequestHandlerSignals
}

func NewAudioFileRequestHandler(db qdb.IDatabase) *AudioFileRequestHandler {
	return &AudioFileRequestHandler{
		db: db,
	}
}

func (w *AudioFileRequestHandler) OnBecameLeader() {
	w.isLeader = true

	w.subscriptionIds = append(w.subscriptionIds, w.db.Notify(&qdb.DatabaseNotificationConfig{
		Type:  "AudioController",
		Field: "AudioFile",
		ContextFields: []string{
			"AudioFile->Description",
			"AudioFile->Content",
		},
	}, w.ProcessNotification))
}

func (w *AudioFileRequestHandler) OnLostLeadership() {
	w.isLeader = false

	for _, id := range w.subscriptionIds {
		w.db.Unnotify(id)
	}

	w.subscriptionIds = []string{}
}

func (w *AudioFileRequestHandler) Init() {

}

func (w *AudioFileRequestHandler) Deinit() {

}

func (w *AudioFileRequestHandler) DoWork() {

}

func (w *AudioFileRequestHandler) ProcessNotification(notification *qdb.DatabaseNotification) {
	if !w.isLeader {
		return
	}

	qdb.Info("[AudioFileRequestHandler::ProcessNotification] Received audio file request: %v", notification)

	if len(notification.Context) == 0 {
		w.Signals.NewRequest.Emit("")
		return
	}

	for _, c := range notification.Context {
		switch c.Name {
		case "AudioFile->Description":
			descriptionField, err := c.Value.UnmarshalNew()
			if err != nil {
				qdb.Error("[AudioFileRequestHandler::ProcessNotification] Failed to unmarshal audio file description: %s", err)
				return
			}

			switch description := descriptionField.(type) {
			case *qdb.String:
				qdb.Info("[AudioFileRequestHandler::ProcessNotification] Playing audio file: %s", description.Raw)
			default:
				qdb.Error("[AudioFileRequestHandler::ProcessNotification] Unknown audio file description type: %T", descriptionField)
				return
			}
		case "AudioFile->Content":
			fileContent, err := c.Value.UnmarshalNew()
			if err != nil {
				qdb.Error("[AudioFileRequestHandler::ProcessNotification] Failed to unmarshal audio file content: %s", err)
				return
			}

			switch content := fileContent.(type) {
			case *qdb.BinaryFile:
				w.Signals.NewRequest.Emit(content.Raw)
			default:
				qdb.Error("[AudioFileRequestHandler::ProcessNotification] Unknown audio file content type: %T", fileContent)
				return
			}
		default:
			qdb.Error("[AudioFileRequestHandler::ProcessNotification] Unknown context field: %s", c.Name)
		}
	}
}
