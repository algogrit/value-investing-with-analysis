package unzipper

import "fmt"

type NotAZipFileErr struct {
	FilePath       string
	IsZipExtension bool
	Op             string
}

func (nazf *NotAZipFileErr) Error() string {
	return fmt.Sprintf("unable to %s, cause file: %s is not a zip file", nazf.Op, nazf.FilePath)
}
