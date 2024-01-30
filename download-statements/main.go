package main

import (
	"encoding/json"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/asdl/entities"
	"codermana.com/go/pkg/asdl/internal/nse"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	// ctx, cancelFn := context.WithCancel(context.Background())
	// defer cancelFn()
	downloader := nse.NewDownloader("./statements")
	// downloader := nse.NewDownloader("%%TEMP%%/asdl") // For Windows

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
	// for _, script := range scripts {
	// 	downloader := nse.NewDownloader("./statements")

	// 	downloader.DownloadFiles(ctx, script)
	// }

	json.NewEncoder(os.Stdout).Encode(scripts)
}
