package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/value_analysis/entities"
	"codermana.com/go/pkg/value_analysis/internal/nse"
	"codermana.com/go/pkg/value_analysis/pkg/unzipper"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

func renameToPDF(filePath string) {
	log.Debug("Renaming:", filePath)

	dirPath := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	fileExtension := filepath.Ext(filePath)
	fileTitle := strings.Split(fileName, fileExtension)[0]

	pdfFileName := fmt.Sprintf("%s/%s.pdf", dirPath, fileTitle)
	os.Rename(filePath, pdfFileName)
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

	var unzipWG sync.WaitGroup

	for _, script := range scripts {
		downloader := nse.NewDownloader("./statements")

		downloader.DownloadFiles(ctx, script)

		unzipWG.Add(1)
		go func(symbol string) {
			defer unzipWG.Done()

			dirPath := fmt.Sprintf("./statements/%s", symbol)
			err := unzipper.UnzipDir(dirPath, false)

			if dirErr, ok := err.(*unzipper.DirErr); ok {
				fileErrs := dirErr.Unwrap()

				for _, fileErr := range fileErrs {
					if nazf, ok := fileErr.(*unzipper.NotAZipFileErr); ok && nazf.IsZipExtension {
						renameToPDF(nazf.FilePath)
					}
				}
			}

			if err != nil {
				fmt.Errorf("unable to archive for symbol: %s cause of %s", symbol, err)
			}
		}(script.NSECode)
	}

	unzipWG.Wait()

	json.NewEncoder(os.Stdout).Encode(scripts)
}
