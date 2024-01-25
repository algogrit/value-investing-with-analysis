package nse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"

	"codermana.com/go/pkg/asdl/entities"
)

var mainSite = "https://www.nseindia.com/companies-listing/corporate-filings-annual-reports"
var scriptsLink = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var annualReports = "https://www.nseindia.com/api/annual-reports?index=equities&symbol=%s"

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"

const downloadThrottleFactor = 1
const cooldownPeriod = 100 * time.Millisecond

type Downloader struct {
	destinationDir string
	client         *http.Client
	cookies        []*http.Cookie
	s              *semaphore.Weighted
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

func (d *Downloader) PopulateStatementsList(s *entities.Script) {
	statementsLink := fmt.Sprintf(annualReports, s.NSECode)

	req, err := http.NewRequest("GET", statementsLink, nil)

	if err != nil {
		log.Fatal("Unable to construct request:", err)
	}

	d.prepareRequest(req)

	resp, err := d.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		log.Fatalf("Unable to fetch statements! status: %d | error: %s", resp.StatusCode, err)
	}

	err = json.NewDecoder(resp.Body).Decode(&s.StatementsList)

	if err != nil {
		io.Copy(os.Stderr, resp.Body)
		log.Fatal("Unable to decode:", err)
	}
}

// TODO: Add expontential backoff
// TODO: Decompress file automatically
func (d *Downloader) downloadFile(ctx context.Context, destinationDir, fileName, fileURL string) error {
	err := d.s.Acquire(ctx, 1)

	if err != nil {
		return err
	}

	defer d.s.Release(1)
	defer time.Sleep(cooldownPeriod)

	log.Info("Downloading :", destinationDir, fileName, fileURL)

	defer log.Info("Completed Downloading :", destinationDir, fileName, fileURL)

	req, err := http.NewRequest("GET", fileURL, nil)

	if err != nil {
		return err
	}

	d.prepareRequest(req)

	resp, err := d.client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unable to download; got status: %d for %s", resp.StatusCode, fileURL)
	}

	f, err := os.OpenFile(filepath.Join(destinationDir, fileName), os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)

	return err
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

	s := semaphore.NewWeighted(downloadThrottleFactor)

	downloader := &Downloader{destinationDir: destinationDir, client: client, s: s}

	downloader.loadCookie()

	return downloader
}
