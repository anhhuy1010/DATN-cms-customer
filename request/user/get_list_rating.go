package user

type (
	ExpertUriParam struct {
		ExpertUuid string `uri:"expert_uuid" binding:"required"`
	}

	GetListRatingRequest struct {
		Keyword string `form:"keyword"`
		Page    int    `form:"page"`
		Limit   int    `form:"limit"`
		Sort    string `form:"sort"`
		Rating  *int   `form:"rating"`
	}
	ListRatingResponse struct {
		Uuid         string `json:"uuid" `
		CustomerUuid string `json:"customer_uuid"`
		CustomerName string `json:"customer_name"`
		Rating       int    `json:"rating"`
		Comment      string `json:"comment"`
	}
)
