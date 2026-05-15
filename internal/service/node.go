package service

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/repository"
)

type NodeService struct {
	repository repository.INodeRepository
	logger     *slog.Logger
}

func NewNodeService(logger *slog.Logger, pool *pgxpool.Pool) *NodeService {
	return &NodeService{logger: logger, repository: repository.NewNode(pool)}
}

func (s *NodeService) GetNodeByID(
	ctx context.Context,
	nodeID int,
) (*models.NodeDetailsResponse, error) {
	response, err := s.repository.Get(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
