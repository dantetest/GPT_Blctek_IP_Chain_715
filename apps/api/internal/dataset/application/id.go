package application

import (
	"crypto/rand"
	"encoding/hex"
)

type IDGenerator interface {
	New(prefix string) string
}

type RandomIDGenerator struct{}

func (RandomIDGenerator) New(prefix string) string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		panic("cryptographic random source unavailable: " + err.Error())
	}
	return prefix + hex.EncodeToString(value[:])
}
