package main

import (
	"context"
	"os/exec"

	"github.com/rqure/qlib/pkg/log"
)

type IAudioPlayer interface {
	Play(filename string) error
	IsPlaying() bool
	Cancel()
}

type AudioPlayer struct {
	cancel context.CancelFunc
}

func NewAudioPlayer() IAudioPlayer {
	return &AudioPlayer{
		cancel: nil,
	}
}

func (a *AudioPlayer) Play(filename string) error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	cmd := exec.CommandContext(ctx, "play", filename)

	go func() {
		if err := cmd.Run(); err != nil {
			log.Error("Failed to play audio")
		}

		a.cancel = nil
	}()

	return nil
}

func (a *AudioPlayer) IsPlaying() bool {
	return a.cancel != nil
}

func (a *AudioPlayer) Cancel() {
	if a.cancel == nil {
		return
	}

	a.cancel()
	a.cancel = nil
}
