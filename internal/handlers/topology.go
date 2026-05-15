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

type topologyHandler struct {
	service *service.TopologyService
	logger  *slog.Logger
}

func newTopologyHandler(logger *slog.Logger, pool *pgxpool.Pool) *topologyHandler {
	return &topologyHandler{service.NewTopologyService(logger, pool), logger}
}

func (h *topologyHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	logIDStr := request.PathValue("log_id")

	logID, err := strconv.Atoi(logIDStr)
	if err != nil {
		http.Error(responseWriter, "Invalid log_id", http.StatusBadRequest)

		return
	}

	topology, err := h.service.GetTopology(request.Context(), logID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(responseWriter, err.Error(), http.StatusNotFound)

			return
		}

		http.Error(responseWriter, "Failed to find topology", http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	err = json.NewEncoder(responseWriter).Encode(topology)
	if err != nil {
		h.logger.Error("parseHandler.ServeHTTP", "error", err)
		http.Error(
			responseWriter,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}
}
