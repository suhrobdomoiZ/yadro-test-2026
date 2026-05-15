package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
)

type Node struct {
	pool *pgxpool.Pool
}

func NewNode(pool *pgxpool.Pool) *Node {
	return &Node{pool}
}

func (r *Node) Get(ctx context.Context, nodeID int) (*models.NodeDetailsResponse, error) {
	query := `
		SELECT n.id, n.log_id, n.guid, 
			   COALESCE(n.description, '') AS description,
			   n.node_type,
			   COALESCE(n.system_image_guid, '') AS system_image_guid,
			   COALESCE(n.base_version, 0) AS base_version,
			   COALESCE(n.class_version, 0) AS class_version,
			   ni.serial_number, ni.part_number, ni.revision, ni.product_name
		FROM nodes n
		LEFT JOIN nodes_info ni ON n.id = ni.node_id
		WHERE n.id = $1
	`

	var (
		node                                            models.NodeDetailsResponse
		serialNumber, partNumber, revision, productName sql.NullString
	)

	err := r.pool.QueryRow(ctx, query, nodeID).Scan(
		&node.ID, &node.LogID, &node.GUID, &node.Description, &node.NodeType,
		&node.SystemImageGUID, &node.BaseVersion, &node.ClassVersion,
		&serialNumber, &partNumber, &revision, &productName,
	)
	if err != nil {
		return nil, fmt.Errorf("get node by id: %w", err)
	}

	if serialNumber.Valid || partNumber.Valid || revision.Valid || productName.Valid {
		node.Info = &models.NodeInfoDTO{
			SerialNumber: serialNumber.String,
			PartNumber:   partNumber.String,
			Revision:     revision.String,
			ProductName:  productName.String,
		}
	}

	return &node, nil
}
