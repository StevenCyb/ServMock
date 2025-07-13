package ini

// Property represents a key-value pair in a section
type Property struct {
	Name      string
	Value     string
	LineIndex uint64
}

// Section represents a section in the INI file
type Section struct {
	Name      string
	Keys      []Property
	LineIndex uint64
}
