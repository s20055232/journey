package main

import (
	"io"
	"log"
	"os"
)

var Logger *log.Logger

func initialLog() {
	stdOut := os.Stdout
	logFile, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	Logger = log.New(io.MultiWriter(stdOut, logFile), "", log.Lshortfile|log.LstdFlags)
	defer logFile.Close()
	Logger.Println("Initialize logger success.")
}