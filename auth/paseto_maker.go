package auth

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, size should be %d", chacha20poly1305.KeySize)
	}
	pasetoMaker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return pasetoMaker, nil
}

func (pm *PasetoMaker) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return pm.paseto.Encrypt(pm.symmetricKey, payload, nil)
}

func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := pm.paseto.Decrypt(token, pm.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorExpiredToken
	}

	if err := payload.Validate(); err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return payload, nil
}
