package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rydwhelchel/spacetraders/api"
	"github.com/rydwhelchel/spacetraders/view"
)

const spaceTraderDBFile = "spacetrader.db"

func main() {
	/* Setup logging to file */
	year, month, day := time.Now().Date()
	// Currently, each log file will be named by the day the process was started
	// TODO: Should enhance to write to a new file when day rolls over if service is still running
	// FIXME: Make this make the directories if they do not already exist
	logName := fmt.Sprintf("./tmp/logs/%v/%v-%v.log", year, month, day)
	f, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		tmp, err := os.Create(logName)
		if err != nil {
			log.Fatalf("error creating file: %v", err)
		}
		f = tmp
	}
	defer f.Close()
	log.SetOutput(f)

	/* Load env */
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env")
	}

	db := setupDatabase()
	defer db.Close()

	traderService := api.NewTraderService(db)
	// Persist the ""cache"" before shutting down
	defer traderService.PersistData()

	tui := view.NewModel(traderService)

	program := tea.NewProgram(tui)

	// For now ignoring the returned model
	_, err = program.Run()
	if err != nil {
		log.Fatalf("error while running bubbletea program - err: %v", err.Error())
	}
}

// TODO: Clean up some of these panics and clean up logic here, bit messy
func setupDatabase() *sql.DB {
	// If the database already exists
	if _, err := os.ReadFile(spaceTraderDBFile); err == nil {
		sqliteDB, err := sql.Open("sqlite3", spaceTraderDBFile)
		if err != nil {
			log.Fatalf("failed to open DB - %v", err)
		}
		// Creates the table if it doesn't already exist
		createTable(sqliteDB)

		return sqliteDB
	}
	// Else create database & prepare
	file, err := os.Create(spaceTraderDBFile)
	if err != nil {
		log.Fatalf("failed to create DB - %v", err)
	}
	file.Close()
	sqliteDB, err := sql.Open("sqlite3", spaceTraderDBFile)
	if err != nil {
		log.Fatalf("failed to open DB after creation - %v", err)
	}

	createTable(sqliteDB)

	return sqliteDB
}

func createTable(db *sql.DB) {
	// The table we're creating is a fairly stupid json blob
	// we don't need anything overly complex because we are basically using this as a persistent cache of just the TraderData object
	createSpacetraderTable := `CREATE TABLE IF NOT EXISTS spacetrader (
		"data" TEXT
	  );`

	statement, err := db.Prepare(createSpacetraderTable)
	if err != nil {
		log.Fatalf("failed to prepare table statement - %v", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("failed to execute table statement - %v", err)
	}
}
