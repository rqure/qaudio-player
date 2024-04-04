package main

import "os/exec"

type IAudioPlayer interface {
	Play(filename string) error
}

type AudioPlayer struct{}

func NewAudioPlayer() IAudioPlayer {
	return &AudioPlayer{}
}

func (a *AudioPlayer) Play(filename string) error {
	cmd := exec.Command("play", filename)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
