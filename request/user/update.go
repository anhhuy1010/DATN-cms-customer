package user

type (
	UpdateUri struct {
		Uuid string `uri:"uuid"`
	}
	UpdateRequest struct {
		UserName  string `json:"username"`
		Email     string `json:"email"`
		Image     string `json:"image"`
		Password  string `json:"password"`
		Introduce string `json:"introduce"`
	}
)
