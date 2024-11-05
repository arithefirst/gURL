package main

import (
	"log"
)

func main() {}

// Function to make error handling less repetitive
func handle(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
