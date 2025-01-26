package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/rydwhelchel/spacetraders/api/openapi"
)

/**
 * Spacetrader service serves as a middleware for requests via openapi.
 */

type TraderService struct {
	ctx       context.Context
	apiClient *openapi.APIClient
	Data      *TraderData // Consider only allowing access to this via Gets, so we can refresh if needed
	db        *sql.DB
}

func NewTraderService(db *sql.DB) *TraderService {
	config := openapi.NewConfiguration()
	apiClient := openapi.NewAPIClient(config)
	token := os.Getenv("SHIP_TRADER_TOKEN")
	if len(token) == 0 {
		panic("failed to find token")
	}

	ctx := context.WithValue(context.Background(), openapi.ContextAccessToken, token)

	ts := &TraderService{
		ctx:       ctx,
		apiClient: apiClient,
		db:        db,
		Data:      newTraderData(),
	}
	// Prime the ""cache""
	ts.RetrieveData()

	// TODO: Think about whether this needs to be refreshed @ start -- it should theoretically be the same as when server shut down
	// Get most up to date agent data
	ts.GetAgentData()

	// TODO: Determine what calls to make at refresh & make helper for it
	return ts
}

/* TraderData is a sort of "cache" which contains the "state" of the game & my agent/ships
This is used in place of a cache because there are different ways to get data
I.E. I can get ship data by requesting ship data, but I also get ship data by sending a task to a ship. In both cases we should update the ship */
// TODO: Persist trader data so restarting the service does not start from scratch
type TraderData struct {
	// TODO: Consider expiration times
	Agent   *openapi.Agent          `json:"agent"`
	Fleet   map[string]openapi.Ship `json:"fleet"` // Map of ShipSymbol -> Ship
	Systems []openapi.System        `json:"systems"`
}

func newTraderData() *TraderData {
	return &TraderData{
		Fleet:   make(map[string]openapi.Ship),
		Systems: []openapi.System{},
	}
}

type requestFunc int

const (
	getAgentData requestFunc = iota
)

// String representations of each requestFunc enum
var requestString = map[requestFunc]string{
	getAgentData: "GET AgentData",
}

func (ts *TraderService) PersistData() {
	countQuery := `SELECT COUNT(*) FROM spacetrader`
	var count int
	err := ts.db.QueryRow(countQuery).Scan(&count)
	if err != nil {
		log.Fatalf("error while counting - %v", err.Error())
	}

	// Prepare data to be persisted
	jsonData, err := json.Marshal(ts.Data)
	if err != nil {
		log.Fatalf("failed to marshal json data (%+v) - %v", ts.Data, err.Error())
	}

	// We have already inserted a row, so we should update
	if count > 0 {
		updateQuery := `UPDATE spacetrader SET data = ?`
		statement, err := ts.db.Prepare(updateQuery)
		if err != nil {
			log.Fatalf("failed to prep update query - %v", err.Error())
		}
		statement.Exec(jsonData)

	} else {
		insertQuery := `INSERT INTO spacetrader(data) VALUES (?)`
		statement, err := ts.db.Prepare(insertQuery)
		if err != nil {
			log.Fatalf("failed to prep insert query - %v", err.Error())
		}
		statement.Exec(jsonData)
	}
}

func (ts *TraderService) RetrieveData() {
	retrieveQuery := `SELECT * FROM spacetrader`
	var data []byte
	err := ts.db.QueryRow(retrieveQuery).Scan(&data)
	if err != nil {
		log.Fatalf("error while retrieving data - %v", err.Error())
	}
	var td TraderData
	json.Unmarshal(data, &td)
	ts.Data = &td
}

// TODO: Since we're limited to 2 requests per second, we should implement a priority queue system for requests

// GetAgentData panics if an error is encountered
func (ts *TraderService) GetAgentData() openapi.Agent {
	a, _, err := ts.apiClient.AgentsAPI.GetMyAgent(ts.ctx).Execute()
	if err != nil {
		// TODO: Figure out a better way to handle this
		log.Panicf("failed to get agent - %v", err.Error())
	}

	ts.Data.Agent = &a.Data

	return a.Data
}
