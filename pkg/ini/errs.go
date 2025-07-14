package ini

import (
	"errors"
	"fmt"
)

// ErrEmptySectionName indicates an error in the behavior header format.
var ErrEmptySectionName = errors.New("empty section name")

// MalformedPropertyError indicates an error in the property format.
type EmptyKeyError struct {
	SectionName string
}

// Error returns a string representation of the EmptyKeyError.
func (e *EmptyKeyError) Error() string {
	return fmt.Sprintf("empty key in section [%s]", e.SectionName)
}
