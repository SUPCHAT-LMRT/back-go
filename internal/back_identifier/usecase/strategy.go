package usecase

import "context"

type BackIdentifierStrategy interface {
	Handle(ctx context.Context) (string, error)
}
