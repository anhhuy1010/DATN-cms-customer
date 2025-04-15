package user

type (
	UpdateUri struct {
		Uuid string `uri:"uuid"`
	}
	UpdateRequest struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	}
)
