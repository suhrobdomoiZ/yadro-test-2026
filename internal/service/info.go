package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/repository"
)

type InfoService struct {
	repository repository.IInfoRepository
	logger     *slog.Logger
}

func NewInfoService(logger *slog.Logger, pool *pgxpool.Pool) *InfoService {
	return &InfoService{logger: logger, repository: repository.NewInfo(pool)}
}

func (s *InfoService) GetLogInfo(ctx context.Context, logID int) (*models.InfoResponse, error) {
	response, err := s.repository.Get(ctx, logID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
