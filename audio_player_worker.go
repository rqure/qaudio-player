package main

import (
	"context"
	"os"

	"github.com/gammazero/deque"
	"github.com/rqure/qlib/pkg/app"
	"github.com/rqure/qlib/pkg/data"
	"github.com/rqure/qlib/pkg/log"
	"github.com/rqure/qtts"
	"github.com/rqure/qtts/voices"
)

type ttsRequest struct {
	text     string
	language string
}

type AudioPlayerWorker struct {
	audioFileQueue  *deque.Deque[string]
	ttsQueue        *deque.Deque[ttsRequest]
	audioPlayer     IAudioPlayer
	requestToCancel bool
}

func NewAudioPlayerWorker() *AudioPlayerWorker {
	audioPlayer := NewAudioPlayer()
	return &AudioPlayerWorker{
		audioFileQueue: new(deque.Deque[string]),
		ttsQueue:       new(deque.Deque[ttsRequest]),
		audioPlayer:    audioPlayer,
	}
}

func (w *AudioPlayerWorker) Init(context.Context, app.Handle) {
}

func (w *AudioPlayerWorker) Deinit(context.Context) {
}

func (w *AudioPlayerWorker) DoWork(context.Context) {
	if w.requestToCancel {
		w.requestToCancel = false
		w.audioPlayer.Cancel()
		return
	}

	if w.audioPlayer.IsPlaying() {
		return
	}

	if w.ttsQueue.Len() > 0 {
		req := w.ttsQueue.PopFront()
		w.audioPlayer.Cancel()

		language := voices.English
		if req.language != "" {
			language = req.language
		}

		tts := &qtts.Speech{
			Folder:   "/",
			Language: language,
			Handler:  w.audioPlayer}

		err := tts.Speak(req.text)
		if err != nil {
			log.Error("Error while playing TTS %v", err)
		}

		return
	}

	if w.audioFileQueue.Len() > 0 {
		content := w.audioFileQueue.PopFront()
		decoded := data.FileDecode(content)

		if len(decoded) == 0 {
			return
		}

		os.WriteFile("temp.mp3", decoded, 0644)
		w.audioPlayer.Cancel()

		err := w.audioPlayer.Play("temp.mp3")
		if err != nil {
			log.Error("Error while playing AudioFile %v", err)
		}

		return
	}
}

func (w *AudioPlayerWorker) OnAddAudioFileToQueue(ctx context.Context, args ...interface{}) {
	content := args[0].(string)
	if content == "" {
		w.requestToCancel = true
		return
	}

	w.audioFileQueue.PushBack(content)
}

func (w *AudioPlayerWorker) OnAddTtsToQueue(ctx context.Context, args ...interface{}) {
	text := args[0].(string)
	language := args[1].(string)
	w.ttsQueue.PushBack(ttsRequest{
		text:     text,
		language: language,
	})
}
