package response

import "time"

type LoginResponse struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	AccessToken    string    `json:"access_token"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}
