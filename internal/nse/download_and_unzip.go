package nse

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"codermana.com/go/pkg/value_analysis/entities"
	"codermana.com/go/pkg/value_analysis/pkg/unzipper"
)

func renameToPDF(filePath string) {
	log.Debug("Renaming:", filePath)

	dirPath := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	fileExtension := filepath.Ext(filePath)
	fileTitle := strings.Split(fileName, fileExtension)[0]

	pdfFileName := fmt.Sprintf("%s/%s.pdf", dirPath, fileTitle)
	os.Rename(filePath, pdfFileName)
}

func (d *Downloader) DownloadAndUnzip(ctx context.Context, script *entities.Script) {
	symbol := script.NSECode
	d.DownloadFiles(ctx, script)

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
}
