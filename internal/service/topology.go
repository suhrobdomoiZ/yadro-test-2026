package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/repository"
)

type TopologyService struct {
	repository repository.ITopologyRepository
	logger     *slog.Logger
}

func NewTopologyService(logger *slog.Logger, pool *pgxpool.Pool) *TopologyService {
	return &TopologyService{
		repository: repository.NewTopology(pool),
		logger:     logger,
	}
}

func (s *TopologyService) GetTopology(
	ctx context.Context,
	logID int,
) (*models.TopologyResponse, error) {
	response, err := s.repository.Get(ctx, logID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
