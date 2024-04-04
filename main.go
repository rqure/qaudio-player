package main

import (
	"os"
	"os/signal"
	"strconv"
	"time"

	qmq "github.com/rqure/qmq/src"
)

func main() {
	app := qmq.NewQMQApplication("audio-player")
	app.Initialize()
	defer app.Deinitialize()

	audioPlayer := NewAudioPlayer()
	app.AddConsumer("audio-player:file:queue").Initialize()
	app.AddConsumer("audio-player:tts:queue").Initialize()

	audioRequestProcessor := NewAudioRequestProcessor(app, audioPlayer)
	textToSpeechProcessor := NewTextToSpeechProcessor(app, audioPlayer)

	tickRateMs, err := strconv.Atoi(os.Getenv("TICK_RATE_MS"))
	if err != nil {
		tickRateMs = 100
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	ticker := time.NewTicker(time.Duration(tickRateMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-sigint:
			app.Logger().Advise("SIGINT received")
			return
		case <-ticker.C:
			audioRequestProcessor.Tick()
			textToSpeechProcessor.Tick()
		}
	}
}
