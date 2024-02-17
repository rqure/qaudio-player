package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	qmq "github.com/rqure/qmq/src"
)

type AudioPlayer struct {
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{}
}

func (a *AudioPlayer) Play(filename string) error {
	cmd := exec.Command("mpg123", filename)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	app := qmq.NewQMQApplication("audio-player")
	app.Initialize()
	defer app.Deinitialize()

	audioPlayer := NewAudioPlayer()

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

				go func() {
					<-time.After(5 * time.Second)
					popped.Ack()
				}()

				err := audioPlayer.Play(request.Filename)

				if err != nil {
					app.Logger().Panic(fmt.Sprintf("Failed to play audio: %v", err))
					os.Exit(1)
				} else {
					app.Logger().Advise(fmt.Sprintf("Finished playing audio file %s", request.Filename))
				}
			}
		}
	}
}
