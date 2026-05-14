package handlers

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppHandler struct {
	parseHandler *parseHandler
}

func NewAppHandlerHandler(logger *slog.Logger, pool *pgxpool.Pool) *AppHandler {
	return &AppHandler{
		parseHandler: newParseHandler(logger, pool),
	}
}

// GET /api/v1/topology/{log_id}
// GET /api/v1/node/{node_id}
// GET /api/v1/port/{node_id}
// GET /api/v1/log/{log_id}
