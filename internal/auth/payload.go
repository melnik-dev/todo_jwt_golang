package auth

type RegisterRequest struct {
	Name     string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Name     string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=3"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
