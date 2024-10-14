package main

import (
	"os"
	"time"

	qdb "github.com/rqure/qdb/src"
)

func getDatabaseAddress() string {
	addr := os.Getenv("QDB_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}

	return addr
}

func main() {
	db := qdb.NewRedisDatabase(qdb.RedisDatabaseConfig{
		Address: getDatabaseAddress(),
	})

	dbWorker := qdb.NewDatabaseWorker(db)
	leaderElectionWorker := qdb.NewLeaderElectionWorker(db)
	audioFileRequestHandler := NewAudioFileRequestHandler(db)
	textToSpeechRequestHandler := NewTextToSpeechRequestHandler(db)
	bluetoothHeartbeatWorker := NewBluetoothHeartbeatWorker(10 * time.Minute)
	audioPlayerWorker := NewAudioPlayerWorker()
	schemaValidator := qdb.NewSchemaValidator(db)

	schemaValidator.AddEntity("Root", "SchemaUpdateTrigger")
	schemaValidator.AddEntity("AudioController", "AudioFile", "TextToSpeech")
	schemaValidator.AddEntity("MP3File", "Content", "Description")

	dbWorker.Signals.SchemaUpdated.Connect(qdb.Slot(schemaValidator.ValidationRequired))
	dbWorker.Signals.Connected.Connect(qdb.Slot(schemaValidator.ValidationRequired))
	leaderElectionWorker.AddAvailabilityCriteria(func() bool {
		return dbWorker.IsConnected() && schemaValidator.IsValid()
	})

	dbWorker.Signals.Connected.Connect(qdb.Slot(leaderElectionWorker.OnDatabaseConnected))
	dbWorker.Signals.Disconnected.Connect(qdb.Slot(leaderElectionWorker.OnDatabaseDisconnected))

	leaderElectionWorker.Signals.BecameLeader.Connect(qdb.Slot(audioFileRequestHandler.OnBecameLeader))
	leaderElectionWorker.Signals.LosingLeadership.Connect(qdb.Slot(audioFileRequestHandler.OnLostLeadership))

	leaderElectionWorker.Signals.BecameLeader.Connect(qdb.Slot(textToSpeechRequestHandler.OnBecameLeader))
	leaderElectionWorker.Signals.LosingLeadership.Connect(qdb.Slot(textToSpeechRequestHandler.OnLostLeadership))

	audioFileRequestHandler.Signals.NewRequest.Connect(qdb.SlotWithArgs(audioPlayerWorker.OnAddAudioFileToQueue))
	textToSpeechRequestHandler.Signals.NewRequest.Connect(qdb.SlotWithArgs(audioPlayerWorker.OnAddTtsToQueue))
	bluetoothHeartbeatWorker.Signals.Heartbeat.Connect(qdb.SlotWithArgs(audioPlayerWorker.OnAddTtsToQueue))

	// Create a new application configuration
	config := qdb.ApplicationConfig{
		Name: "audio-player",
		Workers: []qdb.IWorker{
			dbWorker,
			leaderElectionWorker,
			audioPlayerWorker,
			audioFileRequestHandler,
			textToSpeechRequestHandler,
			bluetoothHeartbeatWorker,
		},
	}

	// Create a new application
	app := qdb.NewApplication(config)

	// Execute the application
	app.Execute()
}
