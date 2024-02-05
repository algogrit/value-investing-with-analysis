package unzipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"codermana.com/go/pkg/value_analysis/pkg/mathext"
)

func UnzipFile(filePath string, preserveFile bool) error {
	dirPath := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	fileExtension := filepath.Ext(filePath)
	fileTitle := strings.Split(fileName, fileExtension)[0]

	zipListing, err := zip.OpenReader(filePath)

	if err != nil {
		return &NotAZipFileErr{FilePath: filePath, Op: "Unzip", IsZipExtension: fileExtension == ".zip"}
	}
	defer zipListing.Close()

	zFileCount := len(zipListing.File)

	var zFileErrs []error

	for zFileIdx, zFile := range zipListing.File {
		zFileReader, err := zFile.Open()
		zFileExt := filepath.Ext(zFile.Name)

		if err != nil {
			zFileErrs = append(zFileErrs, err)
		}
		defer zFileReader.Close()

		fileNamePad := ""

		if zFileCount > 1 {
			padWidth := mathext.DigitCount(zFileCount)
			fileNamePad = fmt.Sprintf(fmt.Sprintf("-%%0%dd", padWidth), zFileIdx)
		}

		saveLocation := fmt.Sprintf("%s/%s%s%s", dirPath, fileTitle, fileNamePad, zFileExt)

		saveFile, err := os.OpenFile(saveLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			zFileErrs = append(zFileErrs, err)
		}

		_, err = io.Copy(saveFile, zFileReader)

		if err != nil {
			zFileErrs = append(zFileErrs, err)
		}
	}

	if zFileErrs != nil || len(zFileErrs) != 0 {
		return &ZipListingErr{Path: dirPath, Op: "Unzip", Errs: zFileErrs}
	}

	if !preserveFile {
		return os.Remove(filePath)
	}

	return nil
}

func UnzipDir(dirPath string, preserveFile bool) error {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		return err
	}

	var fileErrs []error

	for _, file := range files {
		fileName := file.Name()

		filePath := fmt.Sprintf("%s/%s", dirPath, fileName)

		err := UnzipFile(filePath, preserveFile)

		if err != nil {
			fileErrs = append(fileErrs, err)
		}
	}

	if fileErrs != nil || len(fileErrs) != 0 {
		return &DirErr{Path: dirPath, Op: "Unzip", Errs: fileErrs}
	}
	return nil
}
