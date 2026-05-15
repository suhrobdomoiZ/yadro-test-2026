package repository

import (
	"context"

	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/parser"
)

type IParseRepository interface {
	SaveParsedData(
		ctx context.Context,
		filename string,
		parsedData *parser.ParsedData,
	) (int, error)
}

type ITopologyRepository interface {
	Get(ctx context.Context, logID int) (*models.TopologyResponse, error)
}
