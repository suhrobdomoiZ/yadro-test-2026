package handlers

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppHandler struct {
	ParseHandler    *parseHandler
	TopologyHandler *topologyHandler
	NodeHandler     *nodeHandler
	PortHandler     *portHandler
}

func NewAppHandlerHandler(logger *slog.Logger, pool *pgxpool.Pool) *AppHandler {
	return &AppHandler{
		ParseHandler:    newParseHandler(logger, pool),
		TopologyHandler: newTopologyHandler(logger, pool),
		NodeHandler:     newNodeHandler(logger, pool),
		PortHandler:     newPortHandler(logger, pool),
	}
}

// GET /api/v1/log/{log_id}
