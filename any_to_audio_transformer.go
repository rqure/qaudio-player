package main

import (
	"fmt"

	qmq "github.com/rqure/qmq/src"
	"google.golang.org/protobuf/types/known/anypb"
)

type AnyToAudioTransformer struct {
	logger qmq.Logger
}

func NewAnyToAudioTransformer(logger qmq.Logger) qmq.Transformer {
	return &AnyToAudioTransformer{logger: logger}
}

func (t *AnyToAudioTransformer) Transform(i interface{}) interface{} {
	a, ok := i.(*anypb.Any)
	if !ok {
		t.logger.Error(fmt.Sprintf("AnyToAudioTransformer: invalid input type %T", i))
		return nil
	}

	m := &qmq.AudioRequest{}
	if err := a.UnmarshalTo(m); err != nil {
		t.logger.Error(fmt.Sprintf("AnyToAudioTransformer: failed to unmarshal %v", err))
		return nil
	}

	return m
}
