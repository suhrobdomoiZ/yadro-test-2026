package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suhrobdomoiZ/yadro-test-2026/internal/models"
)

const (
	NodeTypeHost   = 1
	NodeTypeSwitch = 2
)

type Topology struct {
	pool *pgxpool.Pool
}

func NewTopology(pool *pgxpool.Pool) *Topology {
	return &Topology{pool}
}

func (r *Topology) Get(ctx context.Context, logID int) (*models.TopologyResponse, error) {
	query := `
        SELECT 
            n.guid, n.description, n.node_type,
            p.guid as port_guid, p.number, p.state, p.link_speed, p.lid
        FROM nodes n
        LEFT JOIN ports p ON n.id = p.node_id
        WHERE n.log_id = $1
        ORDER BY n.node_type, n.description, p.number
    `

	rows, err := r.pool.Query(ctx, query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodesMap := make(map[string]*models.NodeDTO)

	for rows.Next() {
		var (
			nodeGUID string
			nodeDesc string
			nodeType int

			portGUID  sql.NullString
			portNum   sql.NullInt64
			portState sql.NullInt64
			portSpeed sql.NullInt64
			portLID   sql.NullInt64
		)

		err = rows.Scan(
			&nodeGUID, &nodeDesc, &nodeType,
			&portGUID, &portNum, &portState, &portSpeed, &portLID,
		)
		if err != nil {
			return nil, err
		}

		_, exists := nodesMap[nodeGUID]
		if !exists {
			nodesMap[nodeGUID] = &models.NodeDTO{
				GUID:        nodeGUID,
				Description: nodeDesc,
				NodeType:    nodeType,
				Ports:       []models.PortDTO{},
			}
		}

		if portGUID.Valid {
			port := models.PortDTO{
				GUID:   portGUID.String,
				Number: int(portNum.Int64),
				State:  int(portState.Int64),
				Speed:  int(portSpeed.Int64),
				LID:    int(portLID.Int64),
			}
			nodesMap[nodeGUID].Ports = append(nodesMap[nodeGUID].Ports, port)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	response := &models.TopologyResponse{
		LogID:    logID,
		Switches: []models.NodeDTO{},
		Hosts:    []models.NodeDTO{},
	}

	for _, node := range nodesMap {
		switch node.NodeType {
		case NodeTypeHost:
			response.Hosts = append(response.Hosts, *node)
		case NodeTypeSwitch:
			response.Switches = append(response.Switches, *node)
		}
	}

	return response, nil
}
