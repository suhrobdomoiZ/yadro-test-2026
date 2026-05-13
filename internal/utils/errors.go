package utils

import "errors"

type AppError error

var (
	ErrNoFiles AppError = errors.New(
		"archive does not contain required files (.db_csv, .sharp_an_info)",
	)
	ErrIllegalPath      AppError = errors.New("illegal file path (zip slip)")
	ErrCreateTempDir    AppError = errors.New("failed to create temp dir")
	ErrInvalidCSVFormat AppError = errors.New("invalid .db_csv format")
	ErrInvalidMode      AppError = errors.New("invalid mode")
)
