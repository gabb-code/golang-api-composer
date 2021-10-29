package dto

type EndpointDto struct {
	URL         string                 `json:"url"`
	Method      string                 `json:"method"`
	Payload     map[string]interface{} `json:"payload"`
	Headers     map[string]interface{} `json:"headers"`
	ContentType string                 `json:"content_type"`
}
