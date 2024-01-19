package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=ASIANPAINT"

type Script struct {
	Name    string `json:"companyName"`
	NSECode string `json:"symbol"`
}

type Nifty50Data struct {
	Priority int64  `json:"priority"`
	Meta     Script `json:"meta,omitempty"`
}

type Nifty50Resp struct {
	Data []Nifty50Data `json:"data"`
}

func getNifty50List() []Script {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", scriptsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Unable to fetch scripts:", err)
	}

	var data Nifty50Resp

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		log.Fatal("Unable to decode:", err)
	}

	var output []Script

	for _, datum := range data.Data {
		if datum.Priority == 1 {
			continue
		}
		output = append(output, datum.Meta)
	}

	return output
}

func main() {
	scripts := getNifty50List()

	json.NewEncoder(os.Stdout).Encode(scripts)
}
