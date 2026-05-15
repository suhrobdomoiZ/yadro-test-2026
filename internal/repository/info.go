package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
)

type Info struct {
	pool *pgxpool.Pool
}

func NewInfo(pool *pgxpool.Pool) *Info {
	return &Info{pool}
}

func (r *Info) Get(ctx context.Context, logID int) (*models.InfoResponse, error) {
	query := `
		SELECT l.id, l.filename, l.status, l.created_at,
			   (SELECT COUNT(*) FROM nodes WHERE log_id = l.id) AS node_count,
			   (SELECT COUNT(*) FROM ports WHERE log_id = l.id) AS port_count
		FROM logs l
		WHERE l.id = $1
	`

	var info models.InfoResponse

	err := r.pool.QueryRow(ctx, query, logID).Scan(
		&info.ID, &info.Filename, &info.Status, &info.CreatedAt,
		&info.NodeCount, &info.PortCount,
	)
	if err != nil {
		return nil, fmt.Errorf("get log info by id: %w", err)
	}

	return &info, nil
}
