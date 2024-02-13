package nse

import "strings"

type ErrorList []error

func (el ErrorList) Error() string {
	errStrings := make([]string, 0, len(el))

	for _, err := range el {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, "\n-------\n")
}

func NewErrorList(errs []error) error {
	if errs != nil && len(errs) > 0 {
		return ErrorList(errs)
	}
	return nil
}
