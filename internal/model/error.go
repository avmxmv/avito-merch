package model

type ErrorResponse struct {
	Message     string `json:"error"`
	Code        int    `json:"code"`
	Description string `json:"description,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
