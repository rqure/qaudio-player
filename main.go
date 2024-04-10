package main

import qmq "github.com/rqure/qmq/src"

type NameProvider struct{}

func (np *NameProvider) Get() string {
	return "clock"
}

type TransformerProviderFactory struct{}

func (t *TransformerProviderFactory) Create(components qmq.EngineComponentProvider) qmq.TransformerProvider {
	transformerProvider := qmq.NewDefaultTransformerProvider()
	transformerProvider.Set("consumer:audio-player:file:exchange", []qmq.Transformer{
		qmq.NewMessageToAnyTransformer(components.WithLogger()),
		NewAnyToAudioTransformer(components.WithLogger()),
	})
	transformerProvider.Set("consumer:audio-player:tts:exchange", []qmq.Transformer{
		qmq.NewMessageToAnyTransformer(components.WithLogger()),
		NewAnyToTtsTransformer(components.WithLogger()),
	})
	return transformerProvider
}

func main() {
	engine := qmq.NewDefaultEngine(qmq.DefaultEngineConfig{
		NameProvider:               &NameProvider{},
		TransformerProviderFactory: &TransformerProviderFactory{},
		EngineProcessor:            &EngineProcessor{},
	})
	engine.Run()
}
