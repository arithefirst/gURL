package main

import (
	"fmt"
	"io"
)

func Post(flags *Flags) ([]byte, error) {
	// Use SetupRequest() to setup the connection
	client, host, path, err := SetupRequest(flags)
	if err != nil {
		return nil, err
	}

	// Defer closing the connection
	defer client.Close()

	// Request Format:
	// Protocol / HTTP/ver
	// Host
	// Headers
	// Connection Close/Keep Alive
	// Body

	req := "POST " + path + " HTTP/1.1\r\n"
	req += fmt.Sprintf("Host: %s\r\n", host)
	req += fmt.Sprintf("content-length: %d\r\n", len(flags.PostBody))
	for _, v := range flags.Headers {
		req += fmt.Sprintf("%v\r\n", v)
	}

	if flags.KeepAlive {
		req += "Connection: keep-alive \r\n"
	} else {
		req += "Connection: close \r\n"
	}

	req += "\r\n"

	// Append body to the request
	req += flags.PostBody
	req += "\r\n"

	// Send the request to the host
	_, err = client.Write([]byte(req))
	if err != nil {
		return nil, err
	}

	// Read the response and return to the user
	res, err := io.ReadAll(client)
	if err != nil {
		return nil, err
	}

	return res, nil
}
