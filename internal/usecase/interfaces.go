// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	// Task -.
	Link interface {
		Add(ctx context.Context, toLink string) (string, error)
		GetLink(ctx context.Context, hash string) (string, error)
	}

	// LinkRepo -.
	LinkRepo interface {
		Add(ctx context.Context, hash, toLink string) error
		GetLinkByHash(ctx context.Context, hash string) (string, error)
	}
)
