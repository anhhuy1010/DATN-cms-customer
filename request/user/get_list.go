package user

type (
	GetListRequest struct {
		Keyword  string  `form:"keyword"`
		Username *string `form:"username"`
		Page     int     `form:"page"`
		Limit    int     `form:"limit"`
		Sort     string  `form:"sort"`
		IsActive *int    `form:"is_active" `
		Role     *string `form:"role"`
	}
	ListResponse struct {
		Uuid     string `json:"uuid" `
		UserName string `json:"username"`
		IsActive int    `json:"is_active"`
		Role     string `json:"role"`
	}
)
