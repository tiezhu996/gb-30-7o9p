package dto

// RegisterRequest registers a user.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
}

// LoginRequest logs a user in.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest updates a profile.
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	Avatar   string `json:"avatar" binding:"omitempty,max=255"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Bio      string `json:"bio" binding:"omitempty,max=512"`
}

// LoginResponse carries token + user.
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
