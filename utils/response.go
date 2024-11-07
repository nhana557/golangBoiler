package utils

type ApiResponse struct {
	StatusCode int `json:"status_code"`
	Body interface{} `json:"body"`
}
type Body struct {
	Message string `json:"message"`
	Data interface{} `json:"data"`
}


func SuccessResponse(statusCode int, data interface{}) ApiResponse {
	response := ApiResponse{
		StatusCode: statusCode,
		Body:       Body{Message: "success", Data: data},
	}

	return response
}

func ErrorResponse(statusCode int, message error) ApiResponse {
	response := ApiResponse{
		StatusCode: statusCode,
		Body:       Body{Message: message.Error()},
	}

	return response
}