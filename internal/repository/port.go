package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
)

type Ports struct {
	pool *pgxpool.Pool
}

func NewPorts(pool *pgxpool.Pool) *Ports {
	return &Ports{pool}
}

func (r *Ports) Get(ctx context.Context, nodeID int) ([]models.PortInstance, error) {
	query := `
		SELECT id, log_id, node_id, guid, number, lid, state, physical_state, link_speed, link_width
		FROM ports
		WHERE node_id = $1
		ORDER BY number
	`

	rows, err := r.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("query ports by node id: %w", err)
	}

	defer rows.Close()

	var ports []models.PortInstance

	for rows.Next() {
		var (
			tmpPort                                                         models.PortInstance
			portGUID                                                        sql.NullString
			portLID, portState, portPhysState, portLinkSpeed, portLinkWidth sql.NullInt64
		)

		err = rows.Scan(
			&tmpPort.ID, &tmpPort.LogID, &tmpPort.NodeID, &portGUID, &tmpPort.Number,
			&portLID, &portState, &portPhysState, &portLinkSpeed, &portLinkWidth,
		)
		if err != nil {
			return nil, fmt.Errorf("scan port row: %w", err)
		}

		tmpPort.GUID = portGUID.String
		tmpPort.LID = int(portLID.Int64)
		tmpPort.State = int(portState.Int64)
		tmpPort.PhysicalState = int(portPhysState.Int64)
		tmpPort.LinkSpeed = int(portLinkSpeed.Int64)
		tmpPort.LinkWidth = int(portLinkWidth.Int64)

		ports = append(ports, tmpPort)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("iterate ports rows: %w", err)
	}

	if ports == nil {
		ports = []models.PortInstance{}
	}

	return ports, nil
}
