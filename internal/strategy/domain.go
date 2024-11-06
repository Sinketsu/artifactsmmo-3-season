package strategy

import "context"

type Strategy interface {
	Name() string
	Do(ctx context.Context) error
}
