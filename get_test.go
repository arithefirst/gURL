package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"
)

func TestGet(t *testing.T) {
	// Start testing webserver in a goroutine
	go func() {
		// Normal test
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, `{"string":"Testing String","int":69420,"bool":true}`)
			if err != nil {
				log.Fatal(err)
			}
		})

		// Header test
		http.HandleFunc("/headers", func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, r.Header.Get("test-header"))
			if err != nil {
				log.Fatal(err)
			}
		})

		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	testFlags := Flags{
		Url:     "http://localhost:8080/",
		Headers: nil,
	}

	res, err := Get(&testFlags)
	if err != nil {
		t.Fatal("Test failed with unexpected error: ", err)
	}

	if len(res) < 100 {
		t.Error("Test failed: Response too short (<100 char)")
	}

	resStr := string(res)
	matchHttps, err := regexp.MatchString(`[h,t,p,s]{4,5}\:\/\/`, resStr)
	if err != nil {
		t.Fatal("Regexp crashed with unexpected error: ", err)
	} else if matchHttps {
		t.Error("Test failed: Response contains HTTPS URL")
	}

	// Check for protocol version ex: HTTP/1.1
	matchProtVer, err := regexp.MatchString(`HTTP\/\d\.\d \d{3} [A-Z]*`, resStr)
	if err != nil {
		t.Fatal("Regexp crashed with unexpected error: ", err)
	} else if !matchProtVer {
		t.Error("Test failed: Response does not contain protocol version")
	}

	matchResponse, err := regexp.MatchString(`\{"string":"Testing String","int":69420,"bool":true\}`, resStr)
	if err != nil {
		t.Fatal("Regexp crashed with unexpected error: ", err)
	} else if !matchResponse {
		t.Error("Test failed: Different response than expected (", resStr, ")")
	}

	headerTestFlags := Flags{
		Url:     "http://localhost:8080/headers",
		Headers: []string{"test-header: gurl-test-header"},
	}

	headerCheck, err := Get(&headerTestFlags)
	matchHeader, err := regexp.MatchString(`gurl-test-header`, string(headerCheck))
	if err != nil {
		t.Fatal("Regexp crashed with unexpected error: ", err)
	} else if !matchHeader {
		t.Error("Test failed: Different header than expected (", resStr, ")")
	}

}
