package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
)

func Get(flags *Flags) ([]byte, error) {
	// Parse URL
	parsedURL, err := url.Parse(flags.Url)
	if err != nil {
		return []byte(nil), err
	}

	// If the URL was parsed w/o a protocol prepend "http://"
	if parsedURL.Host == "" {
		parsedURL, err = url.Parse("http://" + flags.Url)
		if err != nil {
			return []byte(nil), err
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
		return []byte(nil), err
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

	// Defer closing the connection
	defer client.Close()

	// Request Format:
	// Protocol / HTTP/ver
	// Host
	// Headers
	// Connection Close/Keep Alive

	req := "GET / HTTP/1.1\r\n"
	req += fmt.Sprintf("Host: %s\r\n", host)

	for k, v := range flags.Headers {
		req += fmt.Sprintf("%v: %v\r\n", k, v)
	}

	req += "Connection: close \r\n"
	req += "\r\n"

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

	return res, nil
}
