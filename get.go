package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
)

func get(addr string, headers map[string]string) ([]byte, error) {
	// Parse URL
	parsedURL, err := url.Parse(addr)
	if err != nil {
		return []byte(nil), err
	}

	// Make a TCP Connection on port 443
	conn, err := net.Dial("tcp", parsedURL.Host+":443")
	if err != nil {
		return []byte(nil), err
	}

	// Create TLS connection
	client := tls.Client(conn, &tls.Config{
		ServerName: parsedURL.Hostname(),
	})

	// Defer closing the connection
	defer client.Close()

	// Request Format:
	// Protocol / HTTP/ver
	// Host
	// Headers
	// Connection Close/Keep Alive

	req := "GET / HTTP/1.1\r\n"
	req += fmt.Sprintf("Host: %s\r\n", parsedURL.Host)

	for k, v := range headers {
		req += fmt.Sprintf("%v: %v\r\n", k, v)
	}

	req += "Connection: close \r\n"
	req += "\r\n"

	fmt.Printf("Request:\n%vResponse:\n", req)

	// Send the request to the host
	_, err = client.Write([]byte(req))
	if err != nil {
		return []byte(nil), err
	}

	// Read the response and return to the user
	res, err := io.ReadAll(client)
	if err != nil {
		return []byte(nil), err
	}

	fmt.Println(string(res))
	return res, nil
}
