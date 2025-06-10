package usecase

import "context"

type GetBackIdentifierUseCase struct {
	strategy BackIdentifierStrategy
}

func NewGetBackIdentifierUseCase(
	strategy BackIdentifierStrategy,
) *GetBackIdentifierUseCase {
	return &GetBackIdentifierUseCase{strategy: strategy}
}

func (u *GetBackIdentifierUseCase) Execute(ctx context.Context) (string, error) {
	return u.strategy.Handle(ctx)
}
