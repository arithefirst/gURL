package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestPost(t *testing.T) {
	go func() {
		// Returns the same thing posted to it
		http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
			// Check if the request method is POST
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Convert the body to a string
			buf := new(strings.Builder)
			_, err := io.Copy(buf, r.Body)
			if err != nil {
				log.Fatal(err)
			}

			_, err = fmt.Fprintf(w, buf.String())
			if err != nil {
				log.Fatal(err)
			}
		})

		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	testFlags := Flags{
		Url:      "http://localhost:8000/post",
		PostBody: `{"string":"Testing String","int":69420,"bool":true}`,
		Headers:  nil,
	}

	res, err := Post(&testFlags)
	if err != nil {
		t.Fatal(err)
	} else {
		resStr := string(res)
		if resStr[strings.Index(resStr, "\r\n\r\n")+4:] != `{"string":"Testing String","int":69420,"bool":true}` {
			t.Fatal("Server resposne != POST body")
		}
	}
}
