package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// TODO: Refactor into NSE entity
var mainSite = "https://www.nseindia.com/companies-listing/corporate-filings-annual-reports"
var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=%s"

var cookies []*http.Cookie

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type StatementMetadata struct {
	FromYear string `json:"fromYr"`
	ToYear   string `json:"toYr"`
	FileLink string `json:"fileName"`
}

type StatementsList struct {
	Data []StatementMetadata
}

type Script struct {
	Name           string         `json:"companyName"`
	NSECode        string         `json:"symbol"`
	StatementsList StatementsList `json:"statements"`
}

func (s *Script) PopulateStatementsList() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	statementsLink := fmt.Sprintf(annualReports, s.NSECode)

	req, err := http.NewRequest("GET", statementsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	req.Header.Set("user-agent", UserAgent)

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	log.Debug("Request Cookie:", req.Header.Get("cookie"))

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch statements! status: %d | error: %s", resp.StatusCode, err)
	}

	err = json.NewDecoder(resp.Body).Decode(&s.StatementsList)

	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("Unable to decode:", err)
	}
}

type Nifty50Data struct {
	Priority int64  `json:"priority"`
	Meta     Script `json:"meta,omitempty"`
}

type Nifty50Resp struct {
	Data []Nifty50Data `json:"data"`
}

func getNifty50List() []*Script {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", scriptsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	req.Header.Set("user-agent", UserAgent)

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch scripts! status: %d | error: %s", resp.StatusCode, err)
	}

	var data Nifty50Resp

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("Unable to decode:", err)
	}

	var output []*Script

	for _, datum := range data.Data {
		if datum.Priority == 1 {
			continue
		}
		script := datum.Meta
		output = append(output, &script)
	}

	return output
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	refreshCookie()
}

func refreshCookie() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", mainSite, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	req.Header.Set("user-agent", UserAgent)

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch cookie! status: %d | error: %s", resp.StatusCode, err)
	}

	cookies = resp.Cookies()

	log.Debug("Cookies:", cookies)
}

// Runs in it's own goroutine
func main() {
	scripts := getNifty50List()

	var wg sync.WaitGroup // Counter: 0
	for _, script := range scripts {
		wg.Add(1)

		go func(script *Script) { // Closure Property
			defer wg.Done()
			script.PopulateStatementsList() // Runs independently
		}(script)
	}

	wg.Wait()
	json.NewEncoder(os.Stdout).Encode(scripts)
}
