package main

import (
	"fmt"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/voices"
	qmq "github.com/rqure/qmq/src"
)

type TextToSpeechProcessor struct {
	app         *qmq.QMQApplication
	audioPlayer IAudioPlayer
}

func NewTextToSpeechProcessor(app *qmq.QMQApplication, audioPlayer IAudioPlayer) IQueueProcessor {
	return &TextToSpeechProcessor{
		app:         app,
		audioPlayer: audioPlayer,
	}
}

func (t *TextToSpeechProcessor) Tick() {
	request := &qmq.QMQTextToSpeechRequest{}
	popped := t.app.Consumer("audio-player:tts:queue").Pop(request)

	if popped == nil {
		return
	}

	t.app.Logger().Advise(fmt.Sprintf("Playing text-to-speech: %s", request.Text))
	popped.Ack()

	speech := htgotts.Speech{
		Folder:   "audio",
		Language: voices.English,
		Handler:  t.audioPlayer}
	err := speech.Speak(request.Text)
	if err != nil {
		t.app.Logger().Panic(fmt.Sprintf("Failed to play text-to-speech: %v", err))
		return
	}

	t.app.Logger().Advise(fmt.Sprintf("Finished playing text-to-speech: %s", request.Text))
}