package main

import (
	"fmt"

	qmq "github.com/rqure/qmq/src"
)

type AudioRequestProcessor struct {
	app         *qmq.QMQApplication
	audioPlayer IAudioPlayer
}

func NewAudioRequestProcessor(app *qmq.QMQApplication, audioPlayer IAudioPlayer) IQueueProcessor {
	return &AudioRequestProcessor{
		app:         app,
		audioPlayer: audioPlayer,
	}
}

func (a *AudioRequestProcessor) Tick() {
	request := &qmq.QMQAudioRequest{}
	popped := a.app.Consumer("audio-player:file:queue").Pop(request)

	if popped == nil {
		return
	}

	a.app.Logger().Advise(fmt.Sprintf("Playing audio file: %s", request.Filename))
	popped.Ack()

	err := a.audioPlayer.Play(request.Filename)
	if err != nil {
		a.app.Logger().Panic(fmt.Sprintf("Failed to play audio: %v", err))
		return
	}

	a.app.Logger().Advise(fmt.Sprintf("Finished playing audio file %s", request.Filename))
}
