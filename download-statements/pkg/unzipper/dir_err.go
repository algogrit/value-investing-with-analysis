package unzipper

import "fmt"

type DirErr struct {
	Path string
	Op   string
	Errs []error
}

func (de *DirErr) Error() string {
	return fmt.Sprintf("Encountered %d error(s) in directory: %s while %s", len(de.Errs), de.Path, de.Op)
}

func (de *DirErr) Unwrap() []error {
	return de.Errs
}
