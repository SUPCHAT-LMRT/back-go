package usecase

import (
	"context"
	"os"
)

type HostnameBackIdentifierStrategy struct {
	hostname string
}

func NewHostnameBackIdentifierStrategy() (BackIdentifierStrategy, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &HostnameBackIdentifierStrategy{hostname: hostname}, nil
}

func (h HostnameBackIdentifierStrategy) Handle(_ context.Context) (string, error) {
	return h.hostname, nil
}
