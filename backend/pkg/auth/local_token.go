package auth

import (
	"time"

	"github.com/o1egl/paseto"
)

type LocalCustomClaims struct {
	ID    uint64           `json:"id"`
	Token paseto.JSONToken `json:"token"`
}

type LocalClient interface {
	CreateAccessToken(id uint64) (string, error)
	VerifyToken(token string) (*LocalCustomClaims, error)
}

type lclient struct {
	version  *paseto.V2
	localKey []byte
}

func NewLocalClient(secret string) LocalClient {
	return &lclient{
		version:  paseto.NewV2(),
		localKey: []byte(secret),
	}
}

func (c *lclient) create(id uint64, expireDuration time.Duration) (string, error) {
	payload := &LocalCustomClaims{
		ID: id,
		Token: paseto.JSONToken{
			IssuedAt:   time.Now(),
			Expiration: time.Now().Add(expireDuration),
		},
	}

	token, err := c.version.Encrypt(c.localKey, payload, "")
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *lclient) CreateAccessToken(id uint64) (string, error) {
	return c.create(id, 24*time.Hour)
}

func (c *lclient) VerifyToken(token string) (*LocalCustomClaims, error) {
	var claims LocalCustomClaims

	err := c.version.Decrypt(token, c.localKey, &claims, nil)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}
