package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	qmq "github.com/rqure/qmq/src"
)

func PlayAudio(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	defer speaker.Clear()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

func main() {
	app := qmq.NewQMQApplication("audio-player")
	app.Initialize()
	defer app.Deinitialize()

	app.AddConsumer("audio-player:queue")

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
			request := &qmq.QMQAudioRequest{}
			popped := app.Consumer("audio-player:queue").Pop(request)

			if popped != nil {
				app.Logger().Advise(fmt.Sprintf("Playing audio file: %s", request.Filename))

				err := PlayAudio(request.Filename)
				if err != nil {
					app.Logger().Error(fmt.Sprintf("Failed to play audio: %v", err))
				}
				popped.Ack()
			}
		}
	}
}
