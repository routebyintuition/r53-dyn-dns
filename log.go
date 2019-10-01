package main

import (
	"io"
	"log"
	"os"
)

var (
	// Info is used for informative log messages like, "IP updated"
	Info *log.Logger
	// Warning is used to warn about non-fatal issues like, "could not update IP, will try again"
	Warning *log.Logger
	// Error is used for issues that could/should be fatal
	Error *log.Logger
)

func (conf *Config) logInit() {
	logFileLocation := conf.LogDirectory + "messages.txt"

	logFile, err := os.OpenFile(logFileLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Could not create/write to log file: ", logFileLocation, ": ", err)
	}

	infoHandle := io.MultiWriter(logFile, os.Stdout)
	warningHandle := io.MultiWriter(logFile, os.Stdout)
	errorHandle := io.MultiWriter(logFile, os.Stdout)

	Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime)
	Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime)
	Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime)
}
