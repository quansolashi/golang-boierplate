package response

type UserResponse struct {
	User User `json:"user"`
}

type UserResponses []*UserResponse

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Users []*User
