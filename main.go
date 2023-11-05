// Package main ...
package main

import (
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Logger *log.Logger


func initialLog(){
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

	pages = 2
	channel := make(chan Job, 15)
	quit := make(chan int)
	go scrape(domain, pages, channel, quit)

	db := connect()
Exit:
	for {
		select {
		case job := <-channel:
			db.Create(&job)
		case <-quit:
			break Exit
		}
	}

	// Read
	// var product Product
	// db.First(&product, 1) // find product with integer primary key
	// db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	// db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(&product, 1)
	// Launch another browser with the same docker container.
	// ll := launcher.MustNewManaged(rodUrl)

	// You can set different flags for each browser.
	// ll.Set("disable-sync").Delete("disable-sync")

	// anotherBrowser := rod.New().Client(ll.MustClient()).MustConnect()

	// fmt.Println(
	// 	anotherBrowser.MustPage("https://go-rod.github.io").MustEval("() => document.title"),
	// )

	// utils.Pause()
}
