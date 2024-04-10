package main

import (
	"fmt"

	qmq "github.com/rqure/qmq/src"
	"google.golang.org/protobuf/types/known/anypb"
)

type AnyToTtsTransformer struct {
	logger qmq.Logger
}

func NewAnyToTtsTransformer(logger qmq.Logger) qmq.Transformer {
	return &AnyToTtsTransformer{logger: logger}
}

func (t *AnyToTtsTransformer) Transform(i interface{}) interface{} {
	a, ok := i.(*anypb.Any)
	if !ok {
		t.logger.Error(fmt.Sprintf("AnyToTtsTransformer: invalid input type %T", i))
		return nil
	}

	m := &qmq.TextToSpeechRequest{}
	if err := a.UnmarshalTo(m); err != nil {
		t.logger.Error(fmt.Sprintf("AnyToTtsTransformer: failed to unmarshal %v", err))
		return nil
	}

	return m
}
