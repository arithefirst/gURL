package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Make sure the URL is inputted
	if len(os.Args) == 1 {
		fmt.Println("Usage: gurl <target URL>")
		os.Exit(0)
	}

	get(os.Args[1], headers)
}

// Function to make error handling less repetitive
func handle(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
