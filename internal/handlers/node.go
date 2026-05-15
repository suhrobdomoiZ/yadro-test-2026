package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/service"
)

type nodeHandler struct {
	service *service.NodeService
	logger  *slog.Logger
}

func newNodeHandler(logger *slog.Logger, pool *pgxpool.Pool) *nodeHandler {
	service := service.NewNodeService(logger, pool)

	return &nodeHandler{service: service, logger: logger}
}

// GET /api/v1/node/{node_id}.
func (h *nodeHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	nodeIDStr := request.PathValue("node_id")

	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		http.Error(responseWriter, "Invalid node_id", http.StatusBadRequest)

		return
	}

	node, err := h.service.GetNodeByID(request.Context(), nodeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(responseWriter, "Node not found", http.StatusNotFound)

			return
		}

		http.Error(responseWriter, "Failed to find node", http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	err = json.NewEncoder(responseWriter).Encode(node)
	if err != nil {
		h.logger.Error("nodeHandler.ServeHTTP", "error", err)
		http.Error(
			responseWriter,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}
}
