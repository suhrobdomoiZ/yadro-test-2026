package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/repository"
)

type PortService struct {
	repository repository.IPortsRepository
	logger     *slog.Logger
}

func NewPortService(logger *slog.Logger, pool *pgxpool.Pool) *PortService {
	return &PortService{logger: logger, repository: repository.NewPorts(pool)}
}

func (s *PortService) GetPortsByID(ctx context.Context, nodeID int) (models.PortResponse, error) {
	ports, err := s.repository.Get(ctx, nodeID)
	if err != nil {
		return models.PortResponse{}, err
	}

	return models.PortResponse{Ports: ports}, nil
}
