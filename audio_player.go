package main

import (
	"context"
	"os/exec"
	"sync"

	"github.com/rqure/qlib/pkg/log"
)

type IAudioPlayer interface {
	Play(filename string) error
	IsPlaying() bool
	Cancel()
}

type AudioPlayer struct {
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewAudioPlayer() IAudioPlayer {
	return &AudioPlayer{
		cancel: nil,
	}
}

func (a *AudioPlayer) Play(filename string) error {
	ctx, cancel := context.WithCancel(context.Background())
	a.mu.Lock()
	a.cancel = cancel
	a.mu.Unlock()

	cmd := exec.CommandContext(ctx, "play", filename)

	go func() {
		if err := cmd.Run(); err != nil {
			log.Error("Failed to play audio")
		}

		a.mu.Lock()
		a.cancel = nil
		a.mu.Unlock()
	}()

	return nil
}

func (a *AudioPlayer) IsPlaying() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.cancel != nil
}

func (a *AudioPlayer) Cancel() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cancel == nil {
		return
	}

	a.cancel()
	a.cancel = nil
}
