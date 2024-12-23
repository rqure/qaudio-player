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
	heartbeatTicker *time.Ticker

	Heartbeat signalslots.Signal
}

func NewBluetoothHeartbeatWorker(heartbeatInterval time.Duration) *BluetoothHeartbeatWorker {
	return &BluetoothHeartbeatWorker{
		heartbeatTicker: time.NewTicker(heartbeatInterval),
	}
}

func (w *BluetoothHeartbeatWorker) Init(context.Context, app.Handle) {
}

func (w *BluetoothHeartbeatWorker) Deinit(context.Context) {
	w.heartbeatTicker.Stop()
}

func (w *BluetoothHeartbeatWorker) DoWork(ctx context.Context) {
	select {
	case <-w.heartbeatTicker.C:
		log.Info("Sending heartbeat")
		tts := " "
		w.Heartbeat.Emit(ctx, tts, "en", "MALE")
	default:
	}
}
