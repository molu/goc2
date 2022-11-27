package logger

import (
	"log"
	"os"
)

func SetFileLogger(logFilePath string) (logFile *os.File, err error) {
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("cannot create the logfile: %v", err)
	}
	// defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	return
}
