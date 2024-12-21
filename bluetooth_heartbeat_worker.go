package main

import (
	"context"
	"time"

	"github.com/rqure/qlib/pkg/app"
	"github.com/rqure/qlib/pkg/log"
	"github.com/rqure/qlib/pkg/signalslots"
)

// Bluetooth Speakers normally sleep after a certain amount of time of inactivity.
// This is a "hack" to keep it alive by playing a silent audio file.
type BluetoothHeartbeatWorker struct {
	lastHeartbeatTime time.Time
	heartbeatInterval time.Duration

	Heartbeat signalslots.Signal
}

func NewBluetoothHeartbeatWorker(heartbeatInterval time.Duration) *BluetoothHeartbeatWorker {
	return &BluetoothHeartbeatWorker{
		heartbeatInterval: heartbeatInterval,
	}
}

func (w *BluetoothHeartbeatWorker) Init(context.Context, app.Handle) {
}

func (w *BluetoothHeartbeatWorker) Deinit(context.Context) {
}

func (w *BluetoothHeartbeatWorker) DoWork(context.Context) {
	if time.Since(w.lastHeartbeatTime) > w.heartbeatInterval {
		w.lastHeartbeatTime = time.Now()

		// Play a silent audio file
		log.Info("Sending heartbeat")
		tts := " "
		w.Heartbeat.Emit(tts)
	}
}
