package workerpool

import (
	"context"
)

type WP chan struct{}

func NewWP(limit int) WP {
	return make(chan struct{}, limit)
}

func (wp WP) Exec(ctx context.Context, work func(ctx context.Context)) error {
	select {
	case <-ctx.Done():
	case wp <- struct{}{}:
		go func() {
			work(ctx)
			select {
			case <-ctx.Done():
			case <-wp:
			}
		}()
	}
	return ctx.Err()
}
