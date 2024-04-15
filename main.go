package main

import qmq "github.com/rqure/qmq/src"

type NameProvider struct{}

func (np *NameProvider) Get() string {
	return "clock"
}

type TransformerProviderFactory struct{}

func (t *TransformerProviderFactory) Create(components qmq.EngineComponentProvider) qmq.TransformerProvider {
	transformerProvider := qmq.NewDefaultTransformerProvider()
	transformerProvider.Set("consumer:audio-player:file:queue", []qmq.Transformer{
		qmq.NewMessageToAnyTransformer(components.WithLogger()),
		NewAnyToAudioTransformer(components.WithLogger()),
	})
	transformerProvider.Set("consumer:audio-player:tts:queue", []qmq.Transformer{
		qmq.NewMessageToAnyTransformer(components.WithLogger()),
		NewAnyToTtsTransformer(components.WithLogger()),
	})
	transformerProvider.Set("producer:audio-player:tts:queue", []qmq.Transformer{
		qmq.NewProtoToAnyTransformer(components.WithLogger()),
		qmq.NewAnyToMessageTransformer(components.WithLogger()),
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
