// Package main ...
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Logger *log.Logger

func initialLog() {
	stdOut := os.Stdout
	logFile, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	Logger = log.New(io.MultiWriter(stdOut, logFile), "", log.Lshortfile|log.LstdFlags)
	Logger.Println("Initialize logger success.")
}

func main() {
	initialLog()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rodURL := os.Getenv("ROD_URL")
	jobCat := os.Getenv("JOBCAT")
	domain := os.Getenv("DOMAIN")
	pages := getPages(rodURL, domain, jobCat)

	Logger.Printf("get pages: %v\n", pages)

	channel := make(chan Job, 15)
	quit := make(chan int)
	go scrape(domain, pages, channel, quit)

	db := connect()
Exit:
	for {
		select {
		case job := <-channel:
			connectRabbitMQAndSend(fmt.Sprintf("%v\n%v", job.Name, job.Description))
			db.Create(&job)
		case <-quit:
			break Exit
		}
	}
}
