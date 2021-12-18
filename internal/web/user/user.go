package user

// CreateRequest create request
type CreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginCredentials login request
type LoginCredentials struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// RegisterRequest register request
type RegisterRequest struct {
	Username        string `json:"username" form:"username"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

// RegisterResponse register response
type RegisterResponse struct {
	ID uint64 `json:"id"`
}

// UpdateReq update request
type UpdateReq struct {
	Status int `json:"status"`
}

// ListRequest list request
type ListRequest struct {
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}
