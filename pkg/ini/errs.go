package ini

import (
	"errors"
	"fmt"
)

var EmptySectionNameError = errors.New("empty section name")

type EmptyKeyError struct {
	SectionName string
}

func (e *EmptyKeyError) Error() string {
	return fmt.Sprintf("empty key in section [%s]", e.SectionName)
}
