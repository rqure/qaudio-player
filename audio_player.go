package main

import (
	"context"
	"os/exec"

	qdb "github.com/rqure/qdb/src"
)

type IAudioPlayer interface {
	Play(filename string)
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

func (a *AudioPlayer) Play(filename string) {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	cmd := exec.CommandContext(ctx, "play", filename)

	go func() {
		if err := cmd.Run(); err != nil {
			qdb.Error("[AudioPlayer::Play] Failed to play audio")
		}

		a.cancel = nil
	}()
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
