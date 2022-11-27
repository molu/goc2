package db

import (
	"database/sql"
	"fmt"
	"log"
)

type DBAgent struct {
	Id        string
	CreatedAt string
	UpdatedAt string
}

func CreateAgent(dbSession *sql.DB, agent *DBAgent) (err error) {
	stmt, err := dbSession.Prepare("INSERT INTO agents (id, created_at, updated_at) VALUES ( ?, ?, ? )")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(&agent.Id, &agent.CreatedAt, &agent.UpdatedAt); err != nil {
		log.Fatal(err)
	}
	return
}

func AgentExists(dbSession *sql.DB, agentID string) bool {
	if _, err := GetAgent(dbSession, agentID); err != nil {
		return false
	}
	fmt.Print("After GetAgent")
	return true
}

func GetAgent(dbSession *sql.DB, agentID string) (agent *DBAgent, err error) {
	agent = &DBAgent{}

	stmt, err := dbSession.Prepare("SELECT * FROM agents WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(agentID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err = rows.Scan(&agent.Id, &agent.CreatedAt, &agent.UpdatedAt); err != nil {
			log.Fatal(err)
		}
	}

	return
}

func UpdateAgent(dbSession *sql.DB, agent *DBAgent) {
	stmt, err := dbSession.Prepare(`
		UPDATE agents 
		SET created_at = ?,
			updated_at = ?,
		WHERE agent_id = ?
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(&agent.CreatedAt, &agent.UpdatedAt, &agent.Id); err != nil {
		log.Fatal(err)
	}
}

func DeleteAgent(dbSession *sql.DB, agentID string) {
	stmt, err := dbSession.Prepare(`
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
