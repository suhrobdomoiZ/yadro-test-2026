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

type infoHandler struct {
	service *service.InfoService
	logger  *slog.Logger
}

func newInfoHandler(logger *slog.Logger, pool *pgxpool.Pool) *infoHandler {
	service := service.NewInfoService(logger, pool)

	return &infoHandler{service: service, logger: logger}
}

// GET /api/v1/log/{log_id}.
func (h *infoHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	logIDStr := request.PathValue("log_id")

	logID, err := strconv.Atoi(logIDStr)
	if err != nil {
		http.Error(responseWriter, "Invalid log_id", http.StatusBadRequest)

		return
	}

	info, err := h.service.GetLogInfo(request.Context(), logID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(responseWriter, "Info not found", http.StatusNotFound)

			return
		}

		http.Error(responseWriter, "Failed to find info", http.StatusInternalServerError)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	err = json.NewEncoder(responseWriter).Encode(info)
	if err != nil {
		h.logger.Error("infoHandler.ServeHTTP", "error", err)
		http.Error(
			responseWriter,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}
}
