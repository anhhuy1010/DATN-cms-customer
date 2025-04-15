package user

type (
	LoginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" bson:"password" binding:"required"`
	}

	LoginResponse struct {
		Token string `json:"token"`
		Role  string `json:"role"`
	}
)
