package main

import (
	"os"
	"time"

	"github.com/rqure/qlib/pkg/app"
	"github.com/rqure/qlib/pkg/app/workers"
	"github.com/rqure/qlib/pkg/data/store"
)

func getStoreAddress() string {
	addr := os.Getenv("Q_ADDR")
	if addr == "" {
		addr = "ws://webgateway:20000/ws"
	}

	return addr
}

func main() {
	s := store.NewWeb(store.WebConfig{
		Address: getStoreAddress(),
	})

	storeWorker := workers.NewStore(s)
	leadershipWorker := workers.NewLeadership(s)
	audioFileRequestHandler := NewAudioFileRequestHandler(s)
	textToSpeechRequestHandler := NewTextToSpeechRequestHandler(s)
	bluetoothHeartbeatWorker := NewBluetoothHeartbeatWorker(10 * time.Minute)
	audioPlayerWorker := NewAudioPlayerWorker()
	schemaValidator := leadershipWorker.GetEntityFieldValidator()

	schemaValidator.RegisterEntityFields("Root", "SchemaUpdateTrigger")
	schemaValidator.RegisterEntityFields("AudioController", "AudioFile", "TextToSpeech", "TTSLanguage")
	schemaValidator.RegisterEntityFields("MP3File", "Content", "Description")

	storeWorker.Connected.Connect(leadershipWorker.OnStoreConnected)
	storeWorker.Disconnected.Connect(leadershipWorker.OnStoreDisconnected)

	leadershipWorker.BecameLeader().Connect(audioFileRequestHandler.OnBecameLeader)
	leadershipWorker.LosingLeadership().Connect(audioFileRequestHandler.OnLostLeadership)

	leadershipWorker.BecameLeader().Connect(textToSpeechRequestHandler.OnBecameLeader)
	leadershipWorker.LosingLeadership().Connect(textToSpeechRequestHandler.OnLostLeadership)

	audioFileRequestHandler.NewRequest.Connect(audioPlayerWorker.OnAddAudioFileToQueue)
	textToSpeechRequestHandler.NewRequest.Connect(audioPlayerWorker.OnAddTtsToQueue)
	bluetoothHeartbeatWorker.Heartbeat.Connect(audioPlayerWorker.OnAddTtsToQueue)

	a := app.NewApplication("audioplayer")
	a.AddWorker(storeWorker)
	a.AddWorker(leadershipWorker)
	a.AddWorker(audioPlayerWorker)
	a.AddWorker(audioFileRequestHandler)
	a.AddWorker(textToSpeechRequestHandler)
	a.AddWorker(bluetoothHeartbeatWorker)
	a.Execute()
}
