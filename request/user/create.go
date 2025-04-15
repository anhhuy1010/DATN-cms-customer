package user

type (
	CreateUserRequest struct {
		UserName string `json:"username" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}
	User struct {
		Uuid     string `json:"uuid"`
		UserName string `json:"username"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
)
