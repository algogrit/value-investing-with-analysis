package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=ASIANPAINT"

type Script struct {
	Name    string `json:"companyName"`
	NSECode string `json:"symbol"`
}

type Nifty50Resp struct {
	Data []struct {
		meta Script
	} `json:"data"`
}

func getNifty50List() []Script {
	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
	// 	},
	// 	Timeout: 5 * time.Second,
	// }

	// os.Setenv("GODEBUG", "http2client=0")

	// req, err := http.NewRequest("GET", scriptsLink, nil)

	// if err != nil {
	// 	log.Fatal("Unable to construct request:", err)
	// }

	// resp, err := client.Do(req)

	// TODO: Get over the `stream error: stream ID 1; INTERNAL_ERROR; received from peer` - Server likely has a bug in Http 2 implementation - Figure out how browser/postman are able to deal with it!
	resp, err := http.Get(scriptsLink)

	if err != nil {
		log.Fatal("Unable to fetch scripts:", err)
	}

	var data Nifty50Resp

	json.NewDecoder(resp.Body).Decode(&data)

	var output []Script

	for _, datum := range data.Data {
		output = append(output, datum.meta)
	}

	return output
}

func main() {
	scripts := getNifty50List()

	json.NewEncoder(os.Stdout).Encode(scripts)
}
