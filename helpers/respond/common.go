package respond

func Success(data interface{}, message string) Respond {
	return Respond{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

func MissingParams() Respond {
	return Respond{
		Code:    1001,
		Message: "Missing params",
		Data:    nil,
	}
}

func CreatedFail() Respond {
	return Respond{
		Code:    1002,
		Message: "Created fail!",
		Data:    nil,
	}
}

func UpdatedFail() Respond {
	return Respond{
		Code:    1003,
		Message: "Updated fail!",
		Data:    nil,
	}
}

func Unauthorized() Respond {
	return Respond{
		Code:    1004,
		Message: "Unauthorized",
		Data:    nil,
	}
}

func Forbidden() Respond {
	return Respond{
		Code:    1005,
		Message: "Forbidden",
		Data:    nil,
	}
}

func ManyRequest() Respond {
	return Respond{
		Code:    1006,
		Message: "Too many request",
		Data:    nil,
	}
}

func NotFound() Respond {
	return Respond{
		Code:    1007,
		Message: "Not found",
		Data:    nil,
	}
}

func MissingHeader() Respond {
	return Respond{
		Code:    1008,
		Message: "Missing request header",
		Data:    nil,
	}
}

func InValidParams() Respond {
	return Respond{
		Code:    1009,
		Message: "Invalid params",
		Data:    nil,
	}
}

func ErrorResponse(message string) Respond {
	return Respond{
		Code:    1010,
		Message: message,
		Data:    nil,
	}
}
