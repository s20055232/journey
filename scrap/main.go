// Package main ...
package main

import (
	"fmt"
)



func main() {
	initialLog()
	loadDotEnv()
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
