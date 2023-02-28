package repo

import (
	"context"
	"fmt"
	"github.com/AntonBraer/urlShorter/pkg/postgres"
	"github.com/jackc/pgx/v4"
)

// LinkRepo - repository for link association.
type LinkRepo struct {
	*postgres.Postgres
}

// NewLinkRepo - create new Link repository.
func NewLinkRepo(pg *postgres.Postgres) *LinkRepo {
	return &LinkRepo{pg}
}

// Add create hash<->link association.
func (r *LinkRepo) Add(ctx context.Context, hash, toLink string) (err error) {
	sql := `INSERT INTO links (hash, to_link) VALUES ($1, $2)`

	if _, err = r.Pool.Exec(ctx, sql, hash, toLink); err != nil {
		return fmt.Errorf("LinkRepo - Add - r.Pool.Exec: %w", err)
	}

	return nil
}

// GetLinkByHash get link by hash
func (r *LinkRepo) GetLinkByHash(ctx context.Context, hash string) (toLink string, err error) {
	sql := `SELECT to_link FROM links WHERE hash = $1`

	if err = r.Pool.QueryRow(ctx, sql, hash).Scan(&toLink); err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("LinkRepo - GetLinkByHash - r.Pool.QueryRow: ErrNoRows")
		}
		return "", fmt.Errorf("LinkRepo - GetLinkByHash - r.Pool.QueryRow: %w, sql = %s", err, sql)
	}
	return toLink, nil
}
