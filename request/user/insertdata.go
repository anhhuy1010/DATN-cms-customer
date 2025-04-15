package user

type (
	GetInsertRequest struct {
		Uuid     string `json:"uuid"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email"  binding:"required"`
		Username string `json:"username" binding:"required"`
	}
	InsertResponse struct {
		Uuid string `json:"uuid" `
	}
)
