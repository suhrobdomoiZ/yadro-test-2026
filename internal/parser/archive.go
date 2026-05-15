package parser

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/suhrobdomoiZ/yadro-test-2026/internal/utils"
)

const (
	maxUncompressedSize = 1 << 30
	dirPerm             = 0o750
	dstPerm             = 0o600
)

func Unzip(path string) (*UnzipResult, error) {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	defer func() {
		err := zipReader.Close()
		if err != nil {
			slog.Warn("unzip: error closing zip", "error", err)
		}
	}()

	tempDir, err := os.MkdirTemp("", "unziped_logs")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}

	result := &UnzipResult{TempDir: tempDir}

	for _, file := range zipReader.File {
		resultPath := filepath.Join(tempDir, file.Name)

		_, err = filepath.Rel(tempDir, resultPath)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", utils.ErrIllegalPath, file.Name)
		}

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(resultPath, dirPerm)
			if err != nil {
				return nil, fmt.Errorf("mkdir: %w", err)
			}

			continue
		}

		err = extractFile(file, resultPath)
		if err != nil {
			return nil, fmt.Errorf("extract file (%s): %w", file.Name, err)
		}

		switch filepath.Ext(file.Name) {
		case ".db_csv":
			result.CSVPath = resultPath
		case ".sharp_an_info":
			result.InfoPath = resultPath
		}
	}

	if result.CSVPath == "" || result.InfoPath == "" {
		return nil, utils.ErrNoFiles
	}

	return result, nil
}

func extractFile(file *zip.File, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}

	defer func() {
		err := src.Close()
		if err != nil {
			slog.Warn("extractFile: error closing src", "error", err)
		}
	}()

	err = os.MkdirAll(filepath.Dir(destPath), dirPerm)
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, dstPerm)
	if err != nil {
		return err
	}

	defer func() {
		err := dst.Close()
		if err != nil {
			slog.Warn("extractFile: error closing dst", "error", err)
		}
	}()

	limitedSrc := io.LimitReader(src, maxUncompressedSize)

	_, err = io.Copy(dst, limitedSrc)
	if err != nil {
		return err
	}

	return nil
}
