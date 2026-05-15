package models

type ParseRequest struct {
	Path string `json:"path"`
}

type ParseResponse struct {
	LogID string `json:"log_id"`
}
