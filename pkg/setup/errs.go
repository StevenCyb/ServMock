package parser

import "strconv"

type MalformedBehaviorHeaderError struct {
	LineIndex uint64
	Line      string
	Details   *string
}

func (e *MalformedBehaviorHeaderError) Error() string {
	if e.Details != nil {
		return "Malformed behavior header at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line + " - " + *e.Details
	}
	return "Malformed behavior header at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line
}

type MalformedPropertyError struct {
	LineIndex uint64
	Line      string
	Details   *string
}

func (e *MalformedPropertyError) Error() string {
	if e.Details != nil {
		return "Malformed property at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line + " - " + *e.Details
	}
	return "Malformed property at line " + strconv.FormatUint(e.LineIndex, 10) + ": " + e.Line
}
