package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
)

func get(addr string) {
	// Parse URL
	parsedURL, err := url.Parse(addr)
	handle(err)

	// Make a TCP Connection on port 443
	conn, err := net.Dial("tcp", parsedURL.Host+":443")

	// Create TLS connection
	client := tls.Client(conn, &tls.Config{
		ServerName: parsedURL.Hostname(),
	})

	// Defer closing the connection
	defer func(client *tls.Conn) {
		err := client.Close()
		handle(err)
	}(client)

	// Request Format:
	// Protocol / HTTP/ver
	// Host
	// Headers
	// Connection Close/Keep Alive

	req := "GET / HTTP/1.1\r\n"
	req += fmt.Sprintf("Host: %s\r\n", parsedURL.Host)
	req += "Connection: close \r\n"
	req += "\r\n"

	fmt.Printf("Request:\n%vResponse:\n", req)

	// Send the request to the host
	_, err = client.Write([]byte(req))
	handle(err)

	// Read the response and return to the user
	res, err := io.ReadAll(client)
	fmt.Println(string(res))
}
