package user

type (
	CheckRoleRequest struct {
		Token string `json:"token"`
	}
	CheckRoleResponse struct {
		UserUuid string `json:"user_uuid" `
		Userrole string `json:"userrole"`
	}
)
