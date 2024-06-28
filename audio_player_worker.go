package main

import (
	"encoding/base64"
	"os"

	"github.com/gammazero/deque"
	qdb "github.com/rqure/qdb/src"
	"github.com/rqure/qtts"
	"github.com/rqure/qtts/voices"
)

type AudioPlayerWorker struct {
	audioFileQueue  *deque.Deque[string]
	ttsQueue        *deque.Deque[string]
	audioPlayer     IAudioPlayer
	tts             *qtts.Speech
	requestToCancel bool
}

func NewAudioPlayerWorker() *AudioPlayerWorker {
	audioPlayer := NewAudioPlayer()
	return &AudioPlayerWorker{
		audioFileQueue: deque.New[string](),
		ttsQueue:       deque.New[string](),
		audioPlayer:    audioPlayer,
		tts: &qtts.Speech{
			Folder:   "/",
			Language: voices.English,
			Handler:  audioPlayer},
	}
}

func (w *AudioPlayerWorker) Init() {
}

func (w *AudioPlayerWorker) Deinit() {
}

func (w *AudioPlayerWorker) DoWork() {
	if w.requestToCancel {
		w.requestToCancel = false
		w.audioPlayer.Cancel()
		return
	}

	if w.audioPlayer.IsPlaying() {
		return
	}

	if w.ttsQueue.Len() > 0 {
		content := w.ttsQueue.PopFront()
		w.audioPlayer.Cancel()
		w.tts.Speak(content)
		return
	}

	if w.audioFileQueue.Len() > 0 {
		content := w.audioFileQueue.PopFront()
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			qdb.Error("[AudioFileRequestHandler::ProcessNotification] Failed to decode audio file content: %s", err)
			return
		}

		os.WriteFile("temp.mp3", decoded, 0644)
		w.audioPlayer.Cancel()
		w.audioPlayer.Play("temp.mp3")
		return
	}
}

func (w *AudioPlayerWorker) OnAddAudioFileToQueue(args ...interface{}) {
	content := args[0].(string)
	if content == "" {
		w.requestToCancel = true
		return
	}

	w.audioFileQueue.PushBack(content)
}

func (w *AudioPlayerWorker) OnAddTtsToQueue(args ...interface{}) {
	content := args[0].(string)
	w.ttsQueue.PushBack(content)
}
