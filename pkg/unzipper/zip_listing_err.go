package unzipper

import "fmt"

type ZipListingErr struct {
	Path string
	Op   string
	Errs []error
}

func (zle *ZipListingErr) Error() string {
	return fmt.Sprintf("Encountered %d error(s) in zip listing: %s while %s", len(zle.Errs), zle.Path, zle.Op)
}

func (zle *ZipListingErr) Unwrap() []error {
	return zle.Errs
}
