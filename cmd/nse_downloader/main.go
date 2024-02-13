package main

import (
	"context"
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/value_analysis/internal/nse"
)

// TODOs
// 1. Download not only Nifty50 but as much as possible
// 2. Download quarterly statements as well as annual statements
// 3. Build a similar tool for BSE

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	downloader := nse.NewDownloader("./statements")

	downloader.Nifty50List()

	downloader.PopulateAllStatementsList() // Blocking for all goroutines

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	for _, script := range downloader.Scripts {
		downloader.DownloadAndUnzip(ctx, script)
	}

	json.NewEncoder(os.Stdout).Encode(downloader.Scripts)
}
