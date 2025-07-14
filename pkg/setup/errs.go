package setup

import "strconv"

// MalformedBehaviorHeaderError indicates an error in the behavior header format.
type MalformedBehaviorHeaderError struct {
	LineIndex uint64
	Line      string
	Details   *string
}

// Error returns a string representation of the MalformedBehaviorHeaderError.
func (e *MalformedBehaviorHeaderError) Error() string {
	if e.Details != nil {
		return "Malformed behavior header at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line + " - " + *e.Details
	}
	return "Malformed behavior header at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line
}


// MalformedPropertyError indicates an error in the property format.
type MalformedPropertyError struct {
	LineIndex uint64
	Line      string
	Details   *string
}

// Error returns a string representation of the MalformedPropertyError.
func (e *MalformedPropertyError) Error() string {
	if e.Details != nil {
		return "Malformed property at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line + " - " + *e.Details
	}
	return "Malformed property at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line
}
