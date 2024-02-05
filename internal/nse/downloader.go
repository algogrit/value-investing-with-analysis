package nse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"time"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/value_analysis/entities"
)

var mainSite = "https://www.nseindia.com/companies-listing/corporate-filings-annual-reports"
var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=%s"

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

const retryLimit = 5

type Downloader struct {
	destinationDir string
	client         *http.Client
	cookies        []*http.Cookie
	retryTracker   map[string]int
	retryMut       sync.Mutex
}

func (d *Downloader) loadCookie() {
	req, err := http.NewRequest("GET", mainSite, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	req.Header.Set("user-agent", UserAgent)

	resp, err := d.client.Do(req)

	if err != nil {
		log.Fatalf("Unable to fetch cookie! status: %d | error: %s", resp.StatusCode, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch cookie! status: %d | error: %s", resp.StatusCode, err)
	}

	d.cookies = resp.Cookies()

	log.Debug("Cookies:", d.cookies)
}

func (d *Downloader) prepareRequest(req *http.Request, ignoreCookie bool) {
	req.Header.Set("user-agent", UserAgent)

	if ignoreCookie {
		return
	}

	for _, cookie := range d.cookies {
		req.AddCookie(cookie)
	}
}

func (d *Downloader) Nifty50List() []*entities.Script {
	req, err := http.NewRequest("GET", scriptsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	d.prepareRequest(req, false)

	resp, err := d.client.Do(req)

	if err != nil {
		log.Fatalf("Unable to fetch scripts! status: %d | error: %s", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
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

func (d *Downloader) PopulateStatementsList(s *entities.Script) {
	statementsLink := fmt.Sprintf(annualReports, s.NSECode)

	req, err := http.NewRequest("GET", statementsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	d.prepareRequest(req, false)

	resp, err := d.client.Do(req)

	if err != nil {
		log.Fatalf("Unable to fetch statements! status: %d | error: %s", resp.StatusCode, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch statements! status: %d | error: %s", resp.StatusCode, err)
	}
	err = json.NewDecoder(resp.Body).Decode(&s.StatementsList)

	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("Unable to decode:", err)
	}
}

func (d *Downloader) downloadFile(ctx context.Context, destinationDir, fileName, fileURL string) (err error) {
	defer func() {
		d.retryMut.Lock()

		if err == nil {
			log.Info("Completed Downloading :", fileURL)
		}

		if err == nil || d.retryTracker[fileURL] >= retryLimit {
			d.retryMut.Unlock()
			return // Naked return
		}

		d.retryTracker[fileURL]++
		d.retryMut.Unlock()

		log.Warnf("Retrying %s for %d time cause of %s...", fileURL, d.retryTracker[fileURL], err)
		err = d.downloadFile(ctx, destinationDir, fileName, fileURL)
	}()

	d.retryMut.Lock()
	currentRetryCount := d.retryTracker[fileURL]
	d.retryMut.Unlock()
	if currentRetryCount > 0 {
		sleepDur := math.Pow(2, float64(currentRetryCount))
		time.Sleep(time.Duration(sleepDur) * time.Second)
	}

	log.Info("Downloading :", destinationDir, fileName, fileURL)

	req, err := http.NewRequest("GET", fileURL, nil)

	if err != nil {
		return
	}

	d.prepareRequest(req, true)

	resp, err := d.client.Do(req)

	if err != nil {
		if resp != nil {
			respBody, err := io.ReadAll(resp.Body)
			log.Debug("Response Body:", string(respBody), err)
			log.Debugf("Response Header: %#v", resp.Header)
		}

		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		log.Debug("Response Body:", string(respBody), err)
		log.Debugf("Response Header: %#v", resp.Header)

		return fmt.Errorf("Unable to download; got status: %d for %s", resp.StatusCode, fileURL)
	}

	f, err := os.OpenFile(filepath.Join(destinationDir, fileName), os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	return
}

func (d *Downloader) DownloadFiles(ctx context.Context, s *entities.Script) {
	scriptDir := filepath.Join(d.destinationDir, s.NSECode)

	err := os.Mkdir(scriptDir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal("Unable to create script dir:", err)
	}

	for _, stmt := range s.StatementsList.Data {
		segments := strings.Split(stmt.FileLink, ".")
		extension := segments[len(segments)-1]

		fileName := fmt.Sprintf("%s-%s.%s", stmt.FromYear, stmt.ToYear, extension)
		err := d.downloadFile(ctx, scriptDir, fileName, stmt.FileLink)

		if err != nil {
			log.Errorf("Unable to download: %#v got %s", stmt, err)
		}
	}
}

func NewDownloader(destinationDir string) *Downloader {
	err := os.Mkdir(destinationDir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	downloader := &Downloader{destinationDir: destinationDir, client: client, retryTracker: make(map[string]int)}

	downloader.loadCookie()

	return downloader
}
