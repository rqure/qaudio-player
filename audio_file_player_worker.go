package main

import (
	"encoding/base64"
	"os"

	qdb "github.com/rqure/qdb/src"
)

type AudioFilePlayerWorker struct {
	db              qdb.IDatabase
	isLeader        bool
	subscriptionIds []string

	audioPlayer IAudioPlayer
}

func NewAudioFilePlayerWorker(db qdb.IDatabase) *AudioFilePlayerWorker {
	return &AudioFilePlayerWorker{
		db:          db,
		audioPlayer: NewAudioPlayer(),
	}
}

func (w *AudioFilePlayerWorker) OnBecameLeader() {
	w.subscriptionIds = []string{}

	w.isLeader = true

	w.db.Notify(&qdb.DatabaseNotificationConfig{
		Type:  "AudioController",
		Field: "AudioFile",
		ContextFields: []string{
			"AudioFile->Description",
			"AudioFile->Content",
		},
	}, w.ProcessNotification)
}

func (w *AudioFilePlayerWorker) OnLostLeadership() {
	w.isLeader = false
}

func (w *AudioFilePlayerWorker) Init() {

}

func (w *AudioFilePlayerWorker) Deinit() {

}

func (w *AudioFilePlayerWorker) DoWork() {

}

func (w *AudioFilePlayerWorker) ProcessNotification(notification *qdb.DatabaseNotification) {
	if !w.isLeader {
		return
	}

	if len(notification.Context) == 0 {
		w.audioPlayer.Cancel()
		return
	}

	for _, c := range notification.Context {
		switch c.Name {
		case "AudioFile->Description":
			descriptionField, err := c.Value.UnmarshalNew()
			if err != nil {
				qdb.Error("[AudioFilePlayerWorker::ProcessNotification] Failed to unmarshal audio file description: %s", err)
				return
			}

			switch description := descriptionField.(type) {
			case *qdb.String:
				qdb.Info("[AudioFilePlayerWorker::ProcessNotification] Playing audio file: %s", description.Raw)
			default:
				qdb.Error("[AudioFilePlayerWorker::ProcessNotification] Unknown audio file description type: %T", descriptionField)
				return
			}
		case "AudioFile->Content":
			fileContent, err := c.Value.UnmarshalNew()
			if err != nil {
				qdb.Error("[AudioFilePlayerWorker::ProcessNotification] Failed to unmarshal audio file content: %s", err)
				return
			}

			switch content := fileContent.(type) {
			case *qdb.String:
				decoded, err := base64.StdEncoding.DecodeString(content.Raw)
				if err != nil {
					qdb.Error("[AudioFilePlayerWorker::ProcessNotification] Failed to decode audio file content: %s", err)
					return
				}

				os.WriteFile("temp.mp3", decoded, 0644)
				w.audioPlayer.Cancel()
				w.audioPlayer.Play("temp.mp3")
			default:
				qdb.Error("[AudioFilePlayerWorker::ProcessNotification] Unknown audio file content type: %T", fileContent)
				return
			}
		default:
			qdb.Error("[AudioFilePlayerWorker::ProcessNotification] Unknown context field: %s", c.Name)
		}
	}
}
