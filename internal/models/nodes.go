package models

type NodeDetailsResponse struct {
	ID              int          `json:"id"`
	LogID           int          `json:"log_id"`
	GUID            string       `json:"guid"`
	Description     string       `json:"description"`
	NodeType        int          `json:"node_type"`
	SystemImageGUID string       `json:"system_image_guid"`
	BaseVersion     int          `json:"base_version"`
	ClassVersion    int          `json:"class_version"`
	Info            *NodeInfoDTO `json:"info,omitempty"`
}

type NodeInfoDTO struct {
	SerialNumber string `json:"serial_number"`
	PartNumber   string `json:"part_number"`
	Revision     string `json:"revision"`
	ProductName  string `json:"product_name"`
}
