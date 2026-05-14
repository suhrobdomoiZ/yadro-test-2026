package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/parser"
)

type Parse struct {
	pool *pgxpool.Pool
}

func NewParse(pool *pgxpool.Pool) *Parse {
	return &Parse{pool}
}

func (r *Parse) SaveParsedData(
	ctx context.Context,
	filename string,
	parsedData *parser.ParsedData,
) (int, error) {
	trx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin trx: %w", err)
	}

	defer func() {
		_ = trx.Rollback(ctx)
	}()

	var logID int

	logsQuery := `
        INSERT INTO logs (filename, status) 
        VALUES ($1, 'processing') 
        RETURNING id
    `

	err = trx.QueryRow(ctx, logsQuery, filename).Scan(&logID)
	if err != nil {
		return 0, fmt.Errorf("insert log: %w", err)
	}

	for _, node := range parsedData.Nodes {
		var nodeID int

		nodeQuery := `
            INSERT INTO nodes (log_id, guid, description, node_type, system_image_guid, base_version, class_version)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id
        `

		err = trx.QueryRow(ctx, nodeQuery,
			logID,
			node.GUID,
			node.Description,
			node.NodeType,
			node.SystemImageGUID,
			node.BaseVersion,
			node.ClassVersion,
		).Scan(&nodeID)
		if err != nil {
			return 0, fmt.Errorf("insert node %s: %w", node.GUID, err)
		}

		if node.Info != nil {
			nodeInfoQuery := `
                INSERT INTO nodes_info (node_id, serial_number, part_number, revision, product_name)
                VALUES ($1, $2, $3, $4, $5)
            `

			_, err = trx.Exec(ctx,
				nodeInfoQuery,
				nodeID,
				node.Info.SerialNumber,
				node.Info.PartNumber,
				node.Info.Revision,
				node.Info.ProductName,
			)
			if err != nil {
				return 0, fmt.Errorf("insert node_info %s: %w", node.GUID, err)
			}
		}

		for _, port := range node.Ports {
			portQuery := `
                INSERT INTO ports (log_id, node_id, number, lid, state, physical_state, link_speed, link_width)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
            `

			_, err = trx.Exec(ctx, portQuery,
				logID,
				nodeID,
				port.Number,
				port.LID,
				port.State,
				port.PhysicalState,
				port.LinkSpeed,
				port.LinkWidth,
			)
			if err != nil {
				return 0, fmt.Errorf("insert port %d for node %s: %w", port.Number, node.GUID, err)
			}
		}
	}

	updateLogStatusQuery := `UPDATE logs SET status = 'completed' WHERE id = $1`

	_, err = trx.Exec(ctx, updateLogStatusQuery, logID)
	if err != nil {
		return 0, fmt.Errorf("update log status: %w", err)
	}

	err = trx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit trx: %w", err)
	}

	return logID, nil
}
