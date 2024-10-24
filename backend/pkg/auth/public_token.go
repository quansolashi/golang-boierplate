package auth

import (
	"time"

	"github.com/o1egl/paseto"
)

type PublicCustomClaims struct {
	ID    uint64           `json:"id"`
	Token paseto.JSONToken `json:"token"`
}

type PublicClient interface {
	Create(id uint64, expireDuration time.Duration) (string, error)
	Verify(token string) (*PublicCustomClaims, error)
}

type pclient struct {
	version    *paseto.V2
	publicKey  []byte
	privateKey []byte
}

func NewPublicClient() PublicClient {
	return &pclient{
		version:    paseto.NewV2(),
		publicKey:  []byte{},
		privateKey: []byte{},
	}
}

func (c *pclient) Create(id uint64, expireDuration time.Duration) (string, error) {
	payload := &PublicCustomClaims{
		ID: id,
		Token: paseto.JSONToken{
			IssuedAt:   time.Now(),
			Expiration: time.Now().Add(expireDuration),
		},
	}

	token, err := c.version.Sign(c.privateKey, payload, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *pclient) Verify(token string) (*PublicCustomClaims, error) {
	var claims PublicCustomClaims

	err := c.version.Verify(token, c.publicKey, &claims, nil)
	if err != nil {
		return nil, err
	}
	return &claims, nil
}
