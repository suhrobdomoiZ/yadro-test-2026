package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/suhrobdomoiZ/yadro-test-2026/internal/utils"
)

const (
	expectedKVParts = 2
)

func ParseSharpInfo(path string, data *ParsedData) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open sharp info file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	var currentGUID string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "---") {
			continue
		}

		if strings.HasPrefix(line, "SW_GUID=") {
			parts := strings.SplitN(line, "=", expectedKVParts)

			if len(parts) != expectedKVParts {
				currentGUID = ""

				continue
			}

			guidVal := strings.TrimSpace(parts[utils.IdxValue])
			candidateGUID := "0x" + guidVal

			if _, exists := data.Nodes[candidateGUID]; !exists {
				currentGUID = ""

				continue
			}

			currentGUID = candidateGUID

			continue
		}

		if currentGUID == "" {
			continue
		}

		if !strings.Contains(line, "=") {
			continue
		}

		parts := strings.SplitN(line, "=", expectedKVParts)
		if len(parts) != expectedKVParts {
			continue
		}

		key := strings.TrimSpace(parts[utils.IdxKey])
		val := strings.TrimSpace(parts[utils.IdxValue])

		node, exists := data.Nodes[currentGUID]
		if !exists {
			continue
		}

		if node.Info == nil {
			node.Info = &NodeInfo{Attributes: make(map[string]string)}
		}

		if node.Info.Attributes == nil {
			node.Info.Attributes = make(map[string]string)
		}

		node.Info.Attributes[key] = val
	}

	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("scanner error in sharp info: %w", err)
	}

	return nil
}
