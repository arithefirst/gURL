package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
)

type Flags struct {
	Url            string
	Version        bool
	ShowResHeaders bool
	KeepAlive      bool
	Headers        map[string]string
}

// CliFlags Parses CLI Flags
func CliFlags() Flags {
	var returnFlags Flags

	// Short Flags
	flag.StringVar(&returnFlags.Url, "u", "", "URL to Request")
	flag.BoolVar(&returnFlags.Version, "v", false, "Display version information")
	flag.BoolVar(&returnFlags.ShowResHeaders, "sh", false, "Show response headers")
	flag.BoolVar(&returnFlags.KeepAlive, "k", false, "Set connection to \"keep-alive\"")

	// Long flags
	flag.StringVar(&returnFlags.Url, "url", "", "URL to Request")
	flag.BoolVar(&returnFlags.Version, "version", false, "Display version information")
	flag.BoolVar(&returnFlags.ShowResHeaders, "show-headers", false, "Show response headers")
	flag.BoolVar(&returnFlags.KeepAlive, "keepalive", false, "Set connection to \"keep-alive\"")
	returnFlags.Headers = make(map[string]string)

	// Custom help message for -h/--help
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  -h, --help            Display this Message\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  -u, --url             URL to Request\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  -v, --version         Display version information\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  -sh, --show-headers   Display version information\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "  -k, --keepalive       Set connection to \"keep-alive\"\n")
	}

	flag.Parse()
	return returnFlags
}

// SetupRequest General setup that all request types will use
func SetupRequest(flags *Flags) (net.Conn, string, error) {
	// Parse URL
	parsedURL, err := url.Parse(flags.Url)
	if err != nil {
		return nil, "", err
	}

	// If the URL was parsed w/o a protocol prepend "http://"
	if parsedURL.Host == "" {
		parsedURL, err = url.Parse("http://" + flags.Url)
		if err != nil {
			return nil, "", err
		}
	}

	var host string
	if parsedURL.Scheme == "http" && parsedURL.Port() == "" {
		// HTTP on Port 80
		host = parsedURL.Hostname()
		host += ":80"
	} else if parsedURL.Scheme == "https" && parsedURL.Port() == "" {
		// HTTPS on Port 443
		host = parsedURL.Hostname()
		host += ":443"
	} else {
		// Otherwise use user specified port
		host = parsedURL.Hostname() + ":" + parsedURL.Port()
	}

	// Make a TCP Connection
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, "", err
	}

	// Use TLS if the request is HTTPS
	var client net.Conn
	if parsedURL.Scheme == "https" {
		client = tls.Client(conn, &tls.Config{
			ServerName: parsedURL.Hostname(),
		})
	} else {
		client = conn
	}

	return client, host, nil
}

func main() {
	// Read CLI Flags
	cliFlags := CliFlags()

	// Print version info if -v/--version is set
	if cliFlags.Version {
		versionData := "Go URL by arithefirst\n"
		versionData += "gURL Version beta+0.2\n"
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

	if !cliFlags.ShowResHeaders {
		// Print the response w/o headers
		resStr := string(res)
		fmt.Print(resStr[strings.Index(resStr, "\r\n\r\n")+4:])
	} else {
		// Print the response w/ headers
		fmt.Print(string(res))
	}
}
