package models

type TopologyResponse struct {
	LogID    int       `json:"log_id"`
	Switches []NodeDTO `json:"switches"`
	Hosts    []NodeDTO `json:"hosts"`
}

type NodeDTO struct {
	GUID        string    `json:"guid"`
	Description string    `json:"description"`
	NodeType    int       `json:"node_type"`
	Ports       []PortDTO `json:"ports"`
}

type PortDTO struct {
	GUID   string `json:"guid"`
	Number int    `json:"number"`
	State  int    `json:"state"`
	Speed  int    `json:"link_speed"`
	LID    int    `json:"lid"`
}
