package main

type AudioPlayer interface {
	Play(filename string) error
}
