package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Flags struct {
	Url     string
	Version bool
	Headers map[string]string
}

func CliFlags() Flags {
	var returnFlags Flags

	// Map the cli flags to the struct
	flag.StringVar(&returnFlags.Url, "u", "", "URL to Request")
	flag.BoolVar(&returnFlags.Version, "v", false, "Display version information")
	returnFlags.Headers = make(map[string]string)

	// Parse flags
	flag.Parse()

	return returnFlags
}

func main() {
	// Read CLI Flags
	cliFlags := CliFlags()

	// Print version info if -v is set
	if cliFlags.Version {
		versionData := "Go URL by arithefirst\n"
		versionData += "gURL Version beta+0.1\n"
		versionData += "---------------------\n"
		versionData += "github.com/arithefirst/gurl\n"
		fmt.Print(versionData)
		return
	}

	// Check to see if the URl is empty
	if cliFlags.Url == "" {
		fmt.Println("Usage: gurl -u <target URL>")
		os.Exit(0)
	}

	// Pass to Get → get.go
	res, err := Get(&cliFlags)
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	// Print the resposne
	fmt.Printf("Response:\n%v", string(res))
}
