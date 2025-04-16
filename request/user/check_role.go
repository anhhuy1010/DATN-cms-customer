package user

import "time"

type (
	CheckRoleRequest struct {
		Token string `json:"token"`
	}
	CheckRoleResponse struct {
		UserUuid string     `json:"user_uuid" `
		StartDay *time.Time `json:"startday" `
		EndDay   *time.Time `json:"endday" `
	}
)
