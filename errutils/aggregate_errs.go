package errutils

import (
	"context"
)

func AggregateErrs(ctx context.Context, dest chan error, src <-chan error, srcInfo string) {
	for {
		select {
		case err, ok := <-src:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case dest <- Wrapf(err, srcInfo):
			}
		case <-ctx.Done():
			return
		}
	}
}
