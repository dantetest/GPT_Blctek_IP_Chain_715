package idempotency

import (
	"errors"
	"strings"
)

var ErrInvalidKey = errors.New("idempotency key must contain 16-128 safe ASCII characters")

type Key string

func Parse(value string) (Key, error) {
	value = strings.TrimSpace(value)
	if len(value) < 16 || len(value) > 128 {
		return "", ErrInvalidKey
	}
	for _, r := range value {
		if !(r == '-' || r == '_' || r == '.' || r == ':' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
			return "", ErrInvalidKey
		}
	}
	return Key(value), nil
}
