package main

import qmq "github.com/rqure/qmq/src"

type NameProvider struct{}

func (np *NameProvider) Get() string {
	return "audio-player"
}

type TransformerProviderFactory struct{}

func (t *TransformerProviderFactory) Create(components qmq.EngineComponentProvider) qmq.TransformerProvider {
	transformerProvider := qmq.NewDefaultTransformerProvider()
	transformerProvider.Set("consumer:audio-player:cmd:play-file", []qmq.Transformer{
		qmq.NewTracePopTransformer(components.WithLogger()),
		qmq.NewMessageToAnyTransformer(components.WithLogger()),
		NewAnyToAudioTransformer(components.WithLogger()),
	})
	transformerProvider.Set("consumer:audio-player:cmd:play-tts", []qmq.Transformer{
		qmq.NewTracePopTransformer(components.WithLogger()),
		qmq.NewMessageToAnyTransformer(components.WithLogger()),
		NewAnyToTtsTransformer(components.WithLogger()),
	})
	transformerProvider.Set("producer:audio-player:cmd:play-tts", []qmq.Transformer{
		qmq.NewProtoToAnyTransformer(components.WithLogger()),
		qmq.NewAnyToMessageTransformer(components.WithLogger(), qmq.AnyToMessageTransformerConfig{
			SourceProvider: components.WithNameProvider(),
		}),
		qmq.NewTracePushTransformer(components.WithLogger()),
	})
	return transformerProvider
}

func main() {
	engine := qmq.NewDefaultEngine(qmq.DefaultEngineConfig{
		NameProvider:               &NameProvider{},
		TransformerProviderFactory: &TransformerProviderFactory{},
		EngineProcessor:            NewEngineProcessor(NewAudioPlayer()),
	})
	engine.Run()
}
