package strategy

import (
	"context"
	"time"
)

type empty struct{}

func Empty() *empty {
	return &empty{}
}

func (s *empty) Name() string {
	return "empty"
}

func (s *empty) Do(ctx context.Context) error {
	time.Sleep(1 * time.Second)
	return nil
}
