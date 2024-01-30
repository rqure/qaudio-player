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

type AudioPlayer struct {
	oldSampleRate beep.SampleRate
	initialized bool
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		initialized: false,
	}
}

func (a *AudioPlayer) PlayAudio(filename string) error {
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
	
	// Total number of samples in the streamer
	totalSamples := streamer.Len()
	// Sample rate (number of samples per second)
	sampleRate := format.SampleRate
	// Duration in seconds
	durationSeconds := float64(totalSamples) / float64(sampleRate)
	// Convert duration to a time.Duration for easy formatting
	duration := time.Duration(durationSeconds * float64(time.Second))
	
	done := make(chan bool)

	if !a.initialized {
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		a.initialized = true
		
		speaker.Play(beep.Seq(stream, beep.Callback(func() {
			done <- true
		})))
	} else {
		resampled := beep.Resample(4, format.SampleRate, a.oldSampleRate, streamer)
		
		speaker.Play(beep.Seq(resampled, beep.Callback(func() {
			done <- true
		})))
	}
	
	a.oldSampleRate = format.SampleRate

	select {
	case <-done:
		return nil
	case <-time.After(duration + (1 * time.Second)):
		return fmt.Errorf("Timeout occurred")
	}
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

				err := audioPlayer.PlayAudio(request.Filename)
				if err != nil {
					app.Logger().Error(fmt.Sprintf("Failed to play audio: %v", err))
				}
				popped.Ack()
			}
		}
	}
}
