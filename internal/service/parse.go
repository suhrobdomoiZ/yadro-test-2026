package service

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/parser"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/repository"
)

type ParseService struct {
	repository repository.IParseRepository
	logger     *slog.Logger
}

func NewParseService(logger *slog.Logger, pool *pgxpool.Pool) *ParseService {
	return &ParseService{logger: logger, repository: repository.NewParse(pool)}
}

func (s *ParseService) OpenArchive(path string) (*parser.UnzipResult, error) {
	result, err := parser.Unzip(path)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ParseService) ParseFiles(
	ctx context.Context,
	filename string,
	unzipResult *parser.UnzipResult,
) (*models.ParseResponse, error) {
	data, err := parser.ParseFiles(unzipResult)
	if err != nil {
		return nil, err
	}

	logID, err := s.repository.SaveParsedData(ctx, filename, data)
	if err != nil {
		return nil, err
	}

	stringLogID := strconv.Itoa(logID)

	return &models.ParseResponse{LogID: stringLogID}, nil
}
