// Package main ...
package main

import (
	"fmt"
	"log"
	"os"
)

// start scraping
func main() {
	initialLog()
	loadDotEnv()
	pages := getPages(rodURL, domain, jobCat)
	Logger.Printf("get pages: %v\n", pages)

	channel := make(chan Job, 15)
	quit := make(chan int)
	go scrape(domain, pages, channel, quit)

	// db := connect()
	num := 0
Exit:
	for {
		select {
		case job := <-channel:
			num += 1
			go writeFile(fmt.Sprint(num), job)

			// connectRabbitMQAndSend(fmt.Sprintf("%v\n%v", job.Name, job.Description))
			// db.Create(&job)
		case <-quit:
			break Exit
		}
	}
}

func writeFile(fileName string, job Job) {
	file, err := os.Create(fmt.Sprintf("data/%v.txt", fileName))
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.WriteString(job.Description + "\n" + job.Company.Description)
	if err != nil {
		log.Fatal(err)
	}
}
