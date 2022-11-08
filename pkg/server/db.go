package server

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName string = "db.sqlite3"

func getDB() (db *sql.DB) {
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	createTables(db)
	// defer db.Close()
	return
}

func createTables(db *sql.DB) (err error) {
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS agents (id TEXT NOT NULL UNIQUE, created_at TEXT NOT NULL, updated_at TEXT)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}

	stmt, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS 
			commands (agent_id TEXT NOT NULL UNIQUE, command TEXT, result TEXT, created_at TEXT NOT NULL, updated_at TEXT)
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	if _, err := stmt.Exec(); err != nil {
		log.Fatal(err)
	}
	return
}

func createRandomCmds(db *sql.DB, agentID string, count int) {
	randomCmds := []string{"whoami", "id", "dir", "uname", "cat /etc/passwd"}
	for i := 0; i < count; i++ {
		cmdID := uuid.NewString()
		cmd := randomCmds[i]
		randomCmd := DBCommand{
			id:        cmdID,
			agentID:   agentID,
			command:   cmd,
			result:    "",
			createdAt: "now",
			updatedAt: "none",
		}
		createAgentCmd(db, randomCmd)
	}

}

func createAgent(db *sql.DB, agent DBAgent) (err error) {
	stmt, err := db.Prepare("INSERT INTO agents (id, created_at, updated_at) VALUES ( ?, ?, ? )")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(agent.id, agent.createdAt, agent.updatedAt); err != nil {
		log.Fatal(err)
	}
	return
}

func agentExists(db *sql.DB, agentID string) bool {
	if _, err := getAgent(db, agentID); err != nil {
		return false
	}
	return true
}

func getAgent(db *sql.DB, agentID string) (agent DBAgent, err error) {
	stmt, err := db.Prepare("SELECT * FROM agents WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(agentID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&agent.id, &agent.createdAt, &agent.updatedAt)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func updateAgent(db *sql.DB, agent DBAgent) {
	stmt, err := db.Prepare(`
		UPDATE agents 
		SET created_at = ?,
			updated_at = ?,
		WHERE agent_id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(agent.createdAt, agent.updatedAt, agent.id); err != nil {
		log.Fatal(err)
	}
}

func deleteAgent(db *sql.DB, agentID string) {
	stmt, err := db.Prepare(`
		DELETE FROM agents WHERE id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(agentID); err != nil {
		log.Fatal(err)
	}
}

func createAgentCmd(db *sql.DB, cmd DBCommand) {
	dbCmd := DBCommand{
		id:        cmd.id,
		agentID:   cmd.agentID,
		command:   cmd.command,
		result:    cmd.result,
		createdAt: cmd.createdAt,
		updatedAt: cmd.updatedAt,
	}
	stmt, err := db.Prepare(`
		INSERT INTO commands (id, agent_id, cmd, result, created_at, updated_at) VALUES ( ?, ?, ?, ?, ? )
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(dbCmd.agentID, dbCmd.command, dbCmd.result, dbCmd.createdAt, dbCmd.updatedAt); err != nil {
		log.Fatal(err)
	}
}

func getAgentCmd(agentID string) (cmd string) {
	log.Printf("get cmd for agent %s", agentID)
	cmd = "whoami"
	return
}

func updateAgentCmd(db *sql.DB, dbCmd DBCommand) {
	stmt, err := db.Prepare(`
		UPDATE commands 
		SET cmd = ?,
			result = ?,
			created_at = ?,
			updated_at = ?,
		WHERE agent_id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(dbCmd.command, dbCmd.result, dbCmd.createdAt, dbCmd.updatedAt, dbCmd.agentID); err != nil {
		log.Fatal(err)
	}
}

func deleteAgentCmd(db *sql.DB, cmdID string) {
	stmt, err := db.Prepare(`
		DELETE FROM commands WHERE id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(cmdID); err != nil {
		log.Fatal(err)
	}
}