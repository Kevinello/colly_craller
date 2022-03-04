package pkg

import (
	"errors"
	"strings"

	"github.com/thoas/go-funk"
)

// GenerateErrorFromList 聚合抛出error
// Usage: defer GenerateErrorFromList(&err, errList)
// @param errList
// @return err
func GenerateErrorFromList(err *error, errList []error) {
	// 去除nil error
	errList = funk.Filter(errList, func(err error) bool {
		return err != nil
	}).([]error)
	if len(errList) > 0 {
		*err = errors.New(strings.Join(funk.Map(errList, func(err error) string {
			return err.Error()
		}).([]string), "\n"))
	} else {
		*err = nil
	}
	return
}
