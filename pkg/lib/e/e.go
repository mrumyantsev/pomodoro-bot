package e

import (
	"fmt"
	"strings"
)

func Wrap(desc string, err error) error {
	return fmt.Errorf("%s: %w", desc, err)
}

func ToPrettyString(err error) string {
	if err == nil {
		return ""
	}
	if err.Error() == "" {
		return err.Error()
	}

	firstUp := strings.ToUpper(err.Error()[0:1])

	return fmt.Sprintf("%s%s.", firstUp, err.Error()[1:])
}
