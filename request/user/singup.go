package user

type (
	SignUpRequest struct {
		Password string `json:"password" binding:"required"`
		UserName string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
)
