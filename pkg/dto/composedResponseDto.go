package dto

type ComposedResponseDto struct {
	StatusCode   int         `json:"status_code"`
	StatusDesc   string      `json:"status_desc"`
	Success      bool        `json:"success"`
	ComposedData interface{} `json:"composed_data"`
}
