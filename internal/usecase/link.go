package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/url"
)

// LinkUseCase - UseCase для работы со ссылками
type LinkUseCase struct {
	repo LinkRepo
}

// NewLinkUseCase - создание нового UseCase
func NewLinkUseCase(r LinkRepo) *LinkUseCase {
	return &LinkUseCase{
		repo: r,
	}
}

// Add - Добавление новой ссылки в бд.
func (uc *LinkUseCase) Add(ctx context.Context, toLink string) (hash string, err error) {
	if err = validateURL(toLink); err != nil {
		return "", fmt.Errorf("LinkUseCase - Add: toLink in incorrect format %w", err)

	}

	for hash == "" {
		hash = generateRandomString()
		var link string
		if link, _ = uc.repo.GetLinkByHash(ctx, hash); link != "" {
			hash = ""
		}
	}

	err = uc.repo.Add(ctx, hash, toLink)
	if err != nil {
		return "", fmt.Errorf("LinkUseCase - Add - uc.repo.Add: %w", err)
	}

	return hash, nil
}

// GetLink - получение ссылки из бд.
func (uc *LinkUseCase) GetLink(ctx context.Context, hash string) (toLink string, err error) {
	if hash == "" {
		return "", errors.New("LinkUseCase - GetLink: hash is empty")
	}
	if toLink, err = uc.repo.GetLinkByHash(ctx, hash); err != nil {
		return "", fmt.Errorf("LinkUseCase - GetLinkByHash - uc.repo.GetLinkByHash: %w", err)
	}
	return toLink, nil
}

// generateRandomString - генерации случайного хеша функции
func generateRandomString() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	const urlLength = 8
	b := make([]byte, urlLength)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		b[i] = letterBytes[n.Int64()]
	}
	return string(b)
}

// validateURL - валидации url
func validateURL(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", u.Scheme)
	}

	return nil
}
