package models

import "time"

type InfoResponse struct {
	ID        int       `json:"id"`
	Filename  string    `json:"filename"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	NodeCount int       `json:"node_count"`
	PortCount int       `json:"port_count"`
}
