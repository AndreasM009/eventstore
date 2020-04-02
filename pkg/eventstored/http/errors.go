package http

// ErrorResponse for holding error info
type ErrorResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// NewErrorResponse for holding error info
func NewErrorResponse(errorCode, errorMessage string) ErrorResponse {
	return ErrorResponse{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
}
