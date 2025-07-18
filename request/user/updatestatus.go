package user

type (
	UpdateStatusUri struct {
		Uuid string `uri:"uuid"`
	}
	UpdateStatusRequest struct {
		IsActive *int `json:"is_active"`
		IsDelete *int `json:"is_delete"`
	}
)
