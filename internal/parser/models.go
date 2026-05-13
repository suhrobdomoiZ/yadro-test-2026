package parser

type UnzipResult struct {
	TempDir  string
	CSVPath  string
	InfoPath string
}

type Node struct {
	GUID            string
	Description     string
	NodeType        int
	SystemImageGUID string
	BaseVersion     int
	ClassVersion    int

	Ports []Port
	Info  *NodeInfo
}

type Port struct {
	Number        int
	LID           int
	State         int
	PhysicalState int
	LinkSpeed     int
	LinkWidth     int
}

type NodeInfo struct {
	SerialNumber string
	PartNumber   string
	Revision     string
	ProductName  string
}

type ParsedData struct {
	Nodes map[string]*Node
}
