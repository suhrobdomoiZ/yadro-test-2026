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

type portHandler struct {
	service *service.PortService
	logger  *slog.Logger
}

func newPortHandler(logger *slog.Logger, pool *pgxpool.Pool) *portHandler {
	service := service.NewPortService(logger, pool)

	return &portHandler{service: service, logger: logger}
}

// GET /api/v1/port/{node_id}.
func (h *portHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	nodeIDStr := request.PathValue("node_id")

	nodeID, err := strconv.Atoi(nodeIDStr)
	if err != nil {
		http.Error(responseWriter, "Invalid node_id", http.StatusBadRequest)

		return
	}

	ports, err := h.service.GetPortsByID(request.Context(), nodeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(responseWriter, "Ports not found", http.StatusNotFound)

			return
		}

		http.Error(responseWriter, "Failed to find ports", http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	err = json.NewEncoder(responseWriter).Encode(ports)
	if err != nil {
		h.logger.Error("portsHandler.ServeHTTP", "error", err)
		http.Error(
			responseWriter,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}
}
