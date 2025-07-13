package ini

// Key represents a key-value pair in a section
type Key struct {
	Name      string
	Value     string
	LineIndex uint64
}

// Section represents a section in the INI file
type Section struct {
	Name      string
	Keys      []Key
	LineIndex uint64
}
