package user

type (
	GetDetailUri struct {
		Uuid string `uri:"uuid"`
	}
	GetDetailResponse struct {
		Uuid     string `json:"uuid"`
		Email    string `json:"email"`
		Username string `json:"username"`
		IsActive int    `json:"is_active"`
		IsDelete int    `json:"is_delete"`
		Role     string `json:"role"`
	}
)
