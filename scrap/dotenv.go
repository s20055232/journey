package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var rodURL string
var jobCat string
var domain string

func loadDotEnv(path ...string) {
	var err error
	if len(path) == 0 {
		err = godotenv.Load()	
	}else if len(path) > 1 {
		log.Fatal("Too many .env files")
	}else {
		err = godotenv.Load(path[0])
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rodURL = os.Getenv("ROD_URL")
	log.Printf("rod url: %s", rodURL)
	jobCat = os.Getenv("JOBCAT")
	log.Printf("job category: %s", jobCat)
	domain = os.Getenv("DOMAIN")
	log.Printf("domain: %s", domain)
}
