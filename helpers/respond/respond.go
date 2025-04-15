package respond

type Respond struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InternalServerError() Respond {
	return Respond{
		Code:    1010,
		Message: "Internal error",
		Data:    nil,
	}
}

func ErrorCommon(message string) interface{} {
	return Respond{
		Code:    1010,
		Message: message,
		Data:    nil,
	}
}

type PaginationResponse struct {
	Limit int         `json:"limit"   swaggertype:"primitive,string"`
	Page  int         `json:"page"`
	Pages int         `json:"pages"`
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}

func SuccessPagination(data interface{}, page int, limit int, pages int, total int64) PaginationResponse {
	return PaginationResponse{
		Limit: limit,
		Page:  page,
		Pages: pages,
		Total: total,
		Items: data,
	}
}
