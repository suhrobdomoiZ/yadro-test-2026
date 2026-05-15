package parser

import "fmt"

func ParseFiles(res *UnzipResult) (*ParsedData, error) {
	data, err := ParseCSV(res.CSVPath)
	if err != nil {
		return nil, fmt.Errorf("parse CSV failed: %w", err)
	}

	err = ParseSharpInfo(res.InfoPath, data)
	if err != nil {
		return nil, fmt.Errorf("parse sharp info failed: %w", err)
	}

	return data, nil
}
