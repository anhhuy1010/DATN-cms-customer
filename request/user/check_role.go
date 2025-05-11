package user

import "time"

type (
	CheckRoleRequest struct {
		Token string `json:"token"`
	}
	CheckRoleResponse struct {
		UserUuid string     `json:"user_uuid" `
		UserName string     `json:"username" `
		Email    string     `json:"email" `
		StartDay *time.Time `json:"startday" `
		EndDay   *time.Time `json:"endday" `
	}
)
