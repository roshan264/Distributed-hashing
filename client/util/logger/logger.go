package logger

import (
	"os"
	"log"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

func InitLogger(logPath string) logr.Logger {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	stdLogger := log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	baseLogger := stdr.New(stdLogger)
	logger := baseLogger.V(0) 
	return logger

}

