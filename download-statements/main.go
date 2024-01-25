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
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	downloader := nse.NewDownloader()

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
	json.NewEncoder(os.Stdout).Encode(scripts)
}
