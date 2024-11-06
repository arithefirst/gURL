package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Flags struct {
	Url     string
	Headers map[string]string
}

func CliFlags() Flags {
	var returnFlags Flags

	// Map the cli flags to the struct
	flag.StringVar(&returnFlags.Url, "u", "", "URL to Request")
	returnFlags.Headers = make(map[string]string)

	// Parse flags
	flag.Parse()

	// Check to see if the URl is empty
	if returnFlags.Url == "" {
		fmt.Println("Usage: gurl -u <target URL>")
		os.Exit(0)
	}

	return returnFlags
}

func main() {
	// Read CLI Flags
	cliFlags := CliFlags()

	// Pass to Get → get.go
	res, err := Get(&cliFlags)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	// Print the resposne
	fmt.Printf("Response:\n%v", string(res))
}
