package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/value_analysis/entities"
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

	scripts := downloader.Nifty50List()

	var wg sync.WaitGroup
	for _, script := range scripts {
		wg.Add(1)

		go func(script *entities.Script) {
			defer wg.Done()
			downloader.PopulateStatementsList(script)

		}(script)
	}
	wg.Wait()

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	for _, script := range scripts {
		downloader.DownloadAndUnzip(ctx, script)
	}

	json.NewEncoder(os.Stdout).Encode(scripts)
}
