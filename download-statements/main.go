package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/asdl/internal/nse"
)

// // TODO: Refactor into NSE entity

// func (s *Script) PopulateStatementsList() {
// 	client := &http.Client{
// 		Timeout: 5 * time.Second,
// 	}

// 	statementsLink := fmt.Sprintf(annualReports, s.NSECode)

// 	req, err := http.NewRequest("GET", statementsLink, nil)

// 	if err != nil {
// 		log.Fatal("Unable to construct request:", err)
// 	}

// 	req.Header.Set("user-agent", UserAgent)

// 	for _, cookie := range cookies {
// 		req.AddCookie(cookie)
// 	}

// 	log.Debug("Request Cookie:", req.Header.Get("cookie"))

// 	resp, err := client.Do(req)

// 	if err != nil || resp.StatusCode != http.StatusOK {
// 		log.Fatalf("Unable to fetch statements! status: %d | error: %s", resp.StatusCode, err)
// 	}

// 	err = json.NewDecoder(resp.Body).Decode(&s.StatementsList)

// 	if err != nil {
// 		io.Copy(os.Stderr, resp.Body)
// 		log.Fatal("Unable to decode:", err)
// 	}
// }

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	downloader := nse.NewDownloader()

	scripts := downloader.Nifty50List()

	// scripts := getNifty50List()

	// var wg sync.WaitGroup // Counter: 0
	// for _, script := range scripts {
	// 	wg.Add(1)

	// 	go func(script *Script) { // Closure Property
	// 		defer wg.Done()
	// 		script.PopulateStatementsList() // Runs independently
	// 	}(script)
	// }

	// wg.Wait()
	json.NewEncoder(os.Stdout).Encode(scripts)
}
