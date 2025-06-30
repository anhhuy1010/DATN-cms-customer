package user

import "time"

type (
	GetListRequest struct {
		Keyword  string  `form:"keyword"`
		Username *string `form:"username"`
		Page     int     `form:"page"`
		Limit    int     `form:"limit"`
		Sort     string  `form:"sort"`
		IsActive *int    `form:"is_active" `
		IsDelete *int    `form:"is_delete" `
		StartDay *string `form:"start_day" `
	}
	ListResponse struct {
		Uuid     string     `json:"uuid" `
		UserName string     `json:"username"`
		Email    string     `json:"email"`
		IsActive int        `json:"is_active"`
		IsDelete int        `json:"is_delete"`
		StartDay *time.Time `json:"startday"`
		EndDay   *time.Time `json:"endday"`
	}
)
