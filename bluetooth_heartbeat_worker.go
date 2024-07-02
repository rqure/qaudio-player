package main

import (
	"time"

	qdb "github.com/rqure/qdb/src"
)

type BluetoothHeartbeatWorkerSignals struct {
	Heartbeat qdb.Signal
}

// Bluetooth Speakers normally sleep after a certain amount of time of inactivity.
// This is a "hack" to keep it alive by playing a silent audio file.
type BluetoothHeartbeatWorker struct {
	lastHeartbeatTime time.Time
	heartbeatInterval time.Duration
	Signals           BluetoothHeartbeatWorkerSignals
}

func NewBluetoothHeartbeatWorker(heartbeatInterval time.Duration) *BluetoothHeartbeatWorker {
	return &BluetoothHeartbeatWorker{
		heartbeatInterval: heartbeatInterval,
	}
}

func (w *BluetoothHeartbeatWorker) Init() {
}

func (w *BluetoothHeartbeatWorker) Deinit() {
}

func (w *BluetoothHeartbeatWorker) DoWork() {
	if time.Since(w.lastHeartbeatTime) > w.heartbeatInterval {
		w.lastHeartbeatTime = time.Now()

		// Play a silent audio file
		qdb.Info("[BluetoothHeartbeatWorker::DoWork] Sending heartbeat")
		tts := " "
		w.Signals.Heartbeat.Emit(tts)
	}
}
