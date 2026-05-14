package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/service"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/utils"
)

type parseHandler struct {
	service *service.ParseService
	logger  *slog.Logger
}

func newParseHandler(logger *slog.Logger, pool *pgxpool.Pool) *parseHandler {
	service := service.NewParseService(logger, pool)

	return &parseHandler{service: service, logger: logger}
}

// POST /api/v1/parse.
func (h *parseHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	var req models.ParseRequest

	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)

		return
	}

	if req.Path == "" {
		http.Error(responseWriter, "Log archive path is empty", http.StatusBadRequest)

		return
	}

	cleanPath := filepath.Clean(req.Path)

	if strings.Contains(cleanPath, "..") {
		http.Error(responseWriter, "Invalid path characters", http.StatusBadRequest)

		return
	}

	if !strings.HasSuffix(strings.ToLower(cleanPath), ".zip") {
		http.Error(responseWriter, "File must be a .zip archive", http.StatusBadRequest)

		return
	}

	fullPath := filepath.Join("data", cleanPath)

	_, err = os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.Error(responseWriter, "File not found", http.StatusNotFound)

		return
	}

	unzipRes, err := h.service.OpenArchive(fullPath)
	if err != nil {
		if errors.Is(err, utils.ErrNoFiles) {
			http.Error(
				responseWriter,
				"Required files(.db_csv,.sharpen_an_info) not found",
				http.StatusNotFound,
			)

			return
		}

		if errors.Is(err, utils.ErrIllegalPath) {
			http.Error(responseWriter, "Illegal path", http.StatusBadRequest)

			return
		}

		h.logger.Error("parseHandler.ServeHTTP", "error", err)
		http.Error(responseWriter, "Failed to open archive", http.StatusInternalServerError)

		return
	}

	resp, err := h.service.ParseFiles(request.Context(), filepath.Base(fullPath), unzipRes)
	if err != nil {
		h.logger.Error("parseHandler.ServeHTTP", "error", err)
		http.Error(
			responseWriter,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError,
		)

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)

	err = json.NewEncoder(responseWriter).Encode(resp)
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
