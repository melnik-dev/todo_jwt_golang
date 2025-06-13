package auth

type RegisterRequest struct {
	Name     string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Name     string `json:"username" binding:"required,min=3,max=50,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=50"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
