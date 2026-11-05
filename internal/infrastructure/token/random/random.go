package random

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/drew138/abacus/internal/domain"
)

type Tokens struct{}

func New() Tokens {
	return Tokens{}
}

func (Tokens) New(_ context.Context) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("reading random bytes: %w: %w", domain.ErrInternal, err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
