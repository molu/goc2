package db

import (
	"database/sql"
	"log"
)

type DBCommand struct {
	Id        string
	AgentID   string
	Command   string
	Result    string
	CreatedAt string
	UpdatedAt string
}

func CreateCommand(dbSession *sql.DB, cmd *DBCommand) (res sql.Result, err error) {
	stmt, err := dbSession.Prepare(`
		INSERT INTO commands (id, agent_id, command, result, created_at, updated_at) VALUES ( ?, ?, ?, ?, ?, ? )
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err = stmt.Exec(&cmd.Id, &cmd.AgentID, &cmd.Command, &cmd.Result, &cmd.CreatedAt, &cmd.UpdatedAt)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func GetCommand(dbSession *sql.DB, agentID string) (cmd *DBCommand, err error) {
	cmd = &DBCommand{}

	stmt, err := dbSession.Prepare("SELECT * FROM commands WHERE result = '' AND agent_id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(agentID)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err = rows.Scan(&cmd.Id, &cmd.AgentID, &cmd.Command, &cmd.Result, &cmd.CreatedAt, &cmd.UpdatedAt); err != nil {
			log.Fatal(err)
		}
	}

	return
}

func UpdateCommand(dbSession *sql.DB, cmd *DBCommand) {
	stmt, err := dbSession.Prepare(`
		UPDATE commands SET command = ?, result = ?, created_at = ?, updated_at = ? WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(&cmd.Command, &cmd.Result, &cmd.CreatedAt, &cmd.UpdatedAt, &cmd.Id); err != nil {
		log.Fatal(err)
	}
}

func DeleteCommand(dbSession *sql.DB, cmdID string) {
	stmt, err := dbSession.Prepare(`
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
