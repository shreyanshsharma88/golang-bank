package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	GenerateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

type Payload struct {
	Username  string    `json:"username"`
	ExpiresAt int64     `json:"expires_at"`
	ID        uuid.UUID `json:"id"`
	IssuedAt  int64     `json:"issued_at"`
}

var ErrorExpiredToken = errors.New("token is expired")

func NewPayload(username string, duration time.Duration) (*Payload, error) {

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		Username:  username,
		ExpiresAt: time.Now().Add(duration).Unix(),
		ID:        tokenId,
		IssuedAt:  time.Now().Unix(),
	}
	return payload, nil

}


func (p *Payload) Validate() error {
	if time.Now().Unix() > p.ExpiresAt {
		return ErrorExpiredToken
	}
	return nil
}