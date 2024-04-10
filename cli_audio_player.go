package main

import (
	"os/exec"
	"sync"
)

type CliAudioPlayer struct {
	m sync.Mutex
}

func NewAudioPlayer() AudioPlayer {
	return &CliAudioPlayer{}
}

func (a *CliAudioPlayer) Play(filename string) error {
	a.m.Lock()
	defer a.m.Unlock()

	cmd := exec.Command("play", filename)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
