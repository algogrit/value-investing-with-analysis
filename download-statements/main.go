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

// Links from: https://www.nseindia.com/companies-listing/corporate-filings-annual-reports

var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=%s"

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

	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	req.Header.Set("cookie", "defaultLang=en; _ga=GA1.1.142010354.1705473941; AKA_A2=A; bm_mi=1E39062981E8DE359AA2D9953C8D4B10~YAAQ4oMsMSp02z2NAQAAPAesPxbd8GN6KSAGrZ8onPZp4HguBvE0fMHkiLpffBhXHG4crrK63fIjLCezxH0GuR8XVrHuvUNSAYBTBNZOYdaRHEbqB6JxBRnwf1DVO2tnSJk4phiV/vvvEty2jvKU7/xJPWXt3QEDPfgSKssmywFD9HI1IXJQAy1zxtzufH72YhIpCVbfFn0n9Kwzoja700M1n7ShcWcXfXkNSSZ8ebwRMzTMyZCRh9WDXc7aOZD58cK4TmlCyg+vPAI4H3pfYzb6lUol7lhCB1YMOxl4KEhSwjEBE5YRu2MEQoNCwanCqABhoS3mdAYqFrOL+1cUmisZ+Uvt/X3Wy9PSpuCyO5Xy~1; ak_bmsc=D2FE2B9CBCFA15002F4AAE1A2962A246~000000000000000000000000000000~YAAQ4oMsMVZ02z2NAQAAeAusPxaMtRNNGk48Bke7Ks4O1/wxBllG2ZgQ3mSi/yIama2vWNjZYuz0iH1hqvV2wczzBLKhXAOLHUpMhG66uDhbDLZ8+wRHHQ7SzQcCgQ/Fa7s/HgOO+F6vxYUF32plxrFcHeRC9K9unkwDyPkSyQF+5ziPKKaRy8OVgd0z6gbs9twB62Isnrgx2YIum2+946mKj/4PizUUmIaFXfe2Gm9LAt0HXWvtdmWSNBDgNneqgiXcVGpYnFqGBQ9kEEdkpnFP5/me3HF4KRJOYH2YCD3IWUiitbYMXksUFo7ZR6unVTKu8QST/krhDLIHscM1weZamATSBgIrtrB9t4ceEC5Bu9RH4JvG4S7SFFbqWcsi3lFhNUh8TRN/JoicLB1JkKwmmJU+8f6V6L+/mljzhtl50NizNuh0V1ZcZPfqU9eqtFjn+vQUI/X2GsGdybxE8nzHlqDGvxZVQMJuAvjE+MNeZob/0BZ0vA1QL7//GeFvE3Yjx3KXeePEhYbrwoId2tM7RQWSrpKK4Cg3h+XcYNQE37DPA8o+DcOm0KvoG7kAO9nYQhI=; _ga_QJZ4447QD3=GS1.1.1706170256.3.0.1706170272.0.0.0; nsit=MWa7OD3eUu8zTwZGA3zFqcTN; nseappid=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcGkubnNlIiwiYXVkIjoiYXBpLm5zZSIsImlhdCI6MTcwNjE3MDM4NCwiZXhwIjoxNzA2MTc3NTg0fQ.fJP8STllkobe0lUdodAlKAzKNLWNCxmZ_j5oJAXESPM; _ga_87M7PJ3R97=GS1.1.1706170256.5.1.1706170385.0.0.0; bm_sv=8580ADA0C0442BE9E40BF3EFBB656B9A~YAAQ4oMsMdWL2z2NAQAAhnuuPxa3WFI+kkjy8H7Ndm7i/g1df8TJ9iDbiW3WCvmpGIZBmYw+zPKlwn8kJXPNz14cw5EsJpt46t0bTgO54jjFpSDDm4WfX5shSlF1P11TGno0iuMAcSvkXbCBK91MzP5kxLkESzCycBtblyDvqeT1xtMhQ6Dletj8J2M6nL2215AZ3biKO2C83IggjgBHo6H8X9ShBZIhp1Ka/W80mFuL8HDElNN6ro+x4GmmhV9oRsGe~1")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Unable to fetch scripts:", err)
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

	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Unable to fetch scripts:", err)
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
