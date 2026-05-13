package parser

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Unzip(path string) (*UnzipResult, error) {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("parser.Unzip: %w", err) // Используем %w для обертки ошибки
	}
	defer zipReader.Close()

	tempDir, err := os.MkdirTemp("", "unziped_logs")
	if err != nil {
		return nil, fmt.Errorf("parser.Unzip: %w", err)
	}

	result := &UnzipResult{TempDir: tempDir}

	for _, file := range zipReader.File {
		resultPath := filepath.Join(tempDir, file.Name)

		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(resultPath, os.ModePerm); err != nil {
				return nil, fmt.Errorf("parser.Unzip: mkdir: %w", err)
			}
			continue
		}

		err = extractFile(file, resultPath)
		if err != nil {
			return nil, fmt.Errorf("parser.extractFile(%s): %w", file.Name, err)
		}

		switch filepath.Ext(file.Name) {
		case ".db_csv":
			result.CSVPath = resultPath
		case ".sharp_an_info":
			result.InfoPath = resultPath
		}
	}

	if result.CSVPath == "" || result.InfoPath == "" {
		return nil, fmt.Errorf("parser.Unzip: archive does not contain required files (.db_csv, .sharp_an_info)")
	}

	return result, nil
}

func extractFile(file *zip.File, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
