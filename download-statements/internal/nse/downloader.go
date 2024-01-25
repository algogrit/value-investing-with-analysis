package nse

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/asdl/entities"
)

var mainSite = "https://www.nseindia.com/companies-listing/corporate-filings-annual-reports"
var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=%s"

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

type Downloader struct {
	client  *http.Client
	cookies []*http.Cookie
}

func (d *Downloader) prepareRequest(req *http.Request) {
	req.Header.Set("user-agent", UserAgent)

	for _, cookie := range d.cookies {
		req.AddCookie(cookie)
	}
}

func (d *Downloader) Nifty50List() []*entities.Script {
	req, err := http.NewRequest("GET", scriptsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	d.prepareRequest(req)

	resp, err := d.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch scripts! status: %d | error: %s", resp.StatusCode, err)
	}

	var data Nifty50Resp

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("Unable to decode:", err)
	}

	var output []*entities.Script

	for _, datum := range data.Data {
		if datum.Priority == 1 {
			continue
		}
		script := datum.Meta
		output = append(output, &script)
	}

	return output
}

func (d *Downloader) loadCookie() {
	req, err := http.NewRequest("GET", mainSite, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	req.Header.Set("user-agent", UserAgent)

	resp, err := d.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch cookie! status: %d | error: %s", resp.StatusCode, err)
	}

	d.cookies = resp.Cookies()

	log.Debug("Cookies:", d.cookies)
}

func NewDownloader() *Downloader {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	downloader := &Downloader{client: client}

	downloader.loadCookie()

	return downloader
}
