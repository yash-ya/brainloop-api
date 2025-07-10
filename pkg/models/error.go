package models

type ErrorResponse struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"status_code"`
	Error      struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
