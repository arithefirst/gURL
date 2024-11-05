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

	_, err := get(os.Args[1], make(map[string]string))
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}
