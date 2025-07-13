package ini

// Property represents a key-value pair in a section
type Property struct {
	Key       string
	Value     string
	LineIndex uint64
}

// Section represents a section in the INI file
type Section struct {
	Name       string
	Properties []Property
	LineIndex  uint64
}
