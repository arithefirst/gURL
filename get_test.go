package main

import (
	"regexp"
	"testing"
)

func TestGet(t *testing.T) {

	testFlags := Flags{
		Port:    443,
		Url:     "https://icanhazip.com/",
		Headers: make(map[string]string),
	}

	res, err := Get(&testFlags)
	if err != nil {
		t.Fatal("Test failed with unexpected error: ", err)
	}

	if len(res) < 100 {
		t.Error("Test failed: Response too short (<100 char)")
	}

	matchHttps, err := regexp.MatchString(`[h,t,p,s]{4,5}\:\/\/`, string(res))
	if err != nil {
		t.Fatal("Regexp crashed with unexpected error: ", err)
	} else if matchHttps {
		t.Error("Test failed: Response contains HTTPS URL")
	}

	// Check for protocol version ex: HTTP/1.1
	matchProtVer, err := regexp.MatchString(`HTTP\/\d\.\d \d{3} [A-Z]*`, string(res))
	if err != nil {
		t.Fatal("Regexp crashed with unexpected error: ", err)
	} else if !matchProtVer {
		t.Error("Test failed: Response does not contain protocol version")
	}
}
