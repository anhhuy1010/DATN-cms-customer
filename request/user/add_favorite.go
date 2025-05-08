package user

type (
	AddFavoriteRequest struct {
		PostUuid string `json:"post_uuid" binding:"required"`
		PostType string `json:"post_type" binding:"required"`
	}
)
