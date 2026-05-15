package models

type PortResponse struct {
	Ports []PortInstance `json:"ports"`
}

type PortInstance struct {
	ID            int    `json:"id"`
	LogID         int    `json:"log_id"`
	NodeID        int    `json:"node_id"`
	GUID          string `json:"guid"`
	Number        int    `json:"number"`
	LID           int    `json:"lid"`
	State         int    `json:"state"`
	PhysicalState int    `json:"physical_state"`
	LinkSpeed     int    `json:"link_speed"`
	LinkWidth     int    `json:"link_width"`
}
