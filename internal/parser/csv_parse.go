package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/suhrobdomoiZ/yadro-test-2026/internal/utils"
)

const (
	ModeNone = iota
	ModeNodes
	ModePorts
	ModeInfo
)

const (
	minNodeFields = 8
	minPortFields = 22
	minInfoFields = 5
)

func ParseCSV(path string) (*ParsedData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	data := &ParsedData{
		Nodes: make(map[string]*Node),
	}

	scanner := bufio.NewScanner(file)
	mode := ModeNone

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "START_NODES") {
			mode = ModeNodes
			_ = scanner.Scan()

			continue
		}

		if strings.HasPrefix(line, "END_NODES") {
			mode = ModeNone

			continue
		}

		if strings.HasPrefix(line, "START_PORTS") {
			mode = ModePorts
			_ = scanner.Scan()

			continue
		}

		if strings.HasPrefix(line, "END_PORTS") {
			mode = ModeNone

			continue
		}

		if strings.HasPrefix(line, "START_SYSTEM_GENERAL_INFORMATION") {
			mode = ModeInfo
			_ = scanner.Scan()

			continue
		}

		if strings.HasPrefix(line, "END_SYSTEM_GENERAL_INFORMATION") {
			mode = ModeNone

			continue
		}

		switch mode {
		case ModeNodes:
			node, pErr := parseNodeLine(line)
			if pErr != nil {
				return nil, fmt.Errorf("parse node line: %w", pErr)
			}

			data.Nodes[node.GUID] = node

		case ModePorts:
			port, nodeGUID, pErr := parsePortLine(line)
			if pErr != nil {
				return nil, fmt.Errorf("parse port line: %w", pErr)
			}

			if node, exists := data.Nodes[nodeGUID]; exists {
				node.Ports = append(node.Ports, *port)
			}

		case ModeInfo:
			guid, info, pErr := parseInfoLine(line)
			if pErr != nil {
				return nil, fmt.Errorf("parse info line: %w", pErr)
			}

			if node, exists := data.Nodes[guid]; exists {
				node.Info = info
			}

		case ModeNone:

		default:
			return nil, fmt.Errorf("%w: %d", utils.ErrInvalidMode, mode)
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return data, nil
}

func getField(fields []string, index int) (string, error) {
	if index >= len(fields) {
		return "", utils.ErrInvalidCSVFormat
	}

	return strings.Trim(fields[index], "\""), nil
}

func parseInt(val string) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}

	return res
}

func parseNodeLine(line string) (*Node, error) {
	fields := strings.Split(line, ",")
	if len(fields) < minNodeFields {
		return nil, utils.ErrInvalidCSVFormat
	}

	desc, _ := getField(fields, utils.IdxNodeDesc)
	guid, _ := getField(fields, utils.IdxNodeGUID)
	sysImg, _ := getField(fields, utils.IdxNodeSysImg)
	nTypeStr, _ := getField(fields, utils.IdxNodeType)
	baseStr, _ := getField(fields, utils.IdxNodeBaseVer)
	classStr, _ := getField(fields, utils.IdxNodeClassVer)

	return &Node{
		GUID:            guid,
		Description:     desc,
		NodeType:        parseInt(nTypeStr),
		SystemImageGUID: sysImg,
		BaseVersion:     parseInt(baseStr),
		ClassVersion:    parseInt(classStr),
	}, nil
}

func parsePortLine(line string) (*Port, string, error) {
	fields := strings.Split(line, ",")
	if len(fields) < minPortFields {
		return nil, "", utils.ErrInvalidCSVFormat
	}

	nodeGUID, gErr := getField(fields, utils.IdxPortNodeGUID)
	if gErr != nil {
		return nil, "", gErr
	}

	pNumStr, _ := getField(fields, utils.IdxPortPortNum)
	lidStr, _ := getField(fields, utils.IdxPortLID)
	stateStr, _ := getField(fields, utils.IdxPortState)
	phyStr, _ := getField(fields, utils.IdxPortPhyState)
	speedStr, _ := getField(fields, utils.IdxPortSpeed)
	widthStr, _ := getField(fields, utils.IdxPortWidth)

	port := &Port{
		Number:        parseInt(pNumStr),
		LID:           parseInt(lidStr),
		State:         parseInt(stateStr),
		PhysicalState: parseInt(phyStr),
		LinkSpeed:     parseInt(speedStr),
		LinkWidth:     parseInt(widthStr),
	}

	return port, nodeGUID, nil
}

func parseInfoLine(line string) (string, *NodeInfo, error) {
	fields := strings.Split(line, ",")
	if len(fields) < minInfoFields {
		return "", nil, utils.ErrInvalidCSVFormat
	}

	guid, _ := getField(fields, utils.IdxInfoGUID)
	serial, _ := getField(fields, utils.IdxInfoSerial)
	part, _ := getField(fields, utils.IdxInfoPart)
	rev, _ := getField(fields, utils.IdxInfoRev)
	prod, _ := getField(fields, utils.IdxInfoProd)

	info := &NodeInfo{
		SerialNumber: serial,
		PartNumber:   part,
		Revision:     rev,
		ProductName:  prod,
	}

	return guid, info, nil
}
