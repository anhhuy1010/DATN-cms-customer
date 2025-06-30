package user

type (
	GetInfolUri struct {
		Uuid string `uri:"uuid"`
	}
	GetInfoResponse struct {
		Username string `json:"username"`
		Image    string `json:"image"`
		Email    string `json:"email"`
	}
)
