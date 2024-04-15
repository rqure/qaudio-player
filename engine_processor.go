package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	qmq "github.com/rqure/qmq/src"
	"github.com/rqure/qtts"
	"github.com/rqure/qtts/voices"
)

type EngineProcessor struct {
	AudioPlayer AudioPlayer
	Tts         qtts.Speech
}

func NewEngineProcessor(audioPlayer AudioPlayer) qmq.EngineProcessor {
	return &EngineProcessor{
		AudioPlayer: audioPlayer,
		Tts: qtts.Speech{
			Folder:   "audio",
			Language: voices.English,
			Handler:  audioPlayer},
	}
}

func (e *EngineProcessor) Process(cp qmq.EngineComponentProvider) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(10 * time.Minute)

	for {
		select {
		case <-quit:
			return
		case c := <-cp.WithConsumer("audio-player:file:queue").Pop():
			c.Ack()
			r := c.Data().(*qmq.AudioRequest)

			cp.WithLogger().Advise(fmt.Sprintf("Playing audio file: %s", r.Filename))

			err := e.AudioPlayer.Play(r.Filename)
			if err != nil {
				cp.WithLogger().Error(fmt.Sprintf("Failed to play audio: %v", err))
			} else {
				cp.WithLogger().Advise(fmt.Sprintf("Finished playing audio file: %s", r.Filename))
			}
		case c := <-cp.WithConsumer("audio-player:tts:queue").Pop():
			c.Ack()
			r := c.Data().(*qmq.TextToSpeechRequest)

			cp.WithLogger().Advise(fmt.Sprintf("Playing text-to-speech: '%s'", r.Text))

			err := e.Tts.Speak(r.Text)
			if err != nil {
				cp.WithLogger().Error(fmt.Sprintf("Failed to play text-to-speech: %v", err))
			} else {
				cp.WithLogger().Advise(fmt.Sprintf("Finished playing text-to-speech: '%s'", r.Text))
			}
		case <-ticker.C:
			cp.WithLogger().Debug("Playing keepalive audio to bluetooth speaker")
			cp.WithProducer("audio-player:tts:queue").Push(&qmq.TextToSpeechRequest{
				Text: os.Getenv("BLUETOOTH_SPEAKER_KEEPALIVE_TTS") + " ",
			})
		}
	}
}
