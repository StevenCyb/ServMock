package ini

// Key represents a key-value pair in a section
type Key struct {
	Name  string
	Value string
}

// Section represents a section in the INI file
type Section struct {
	Name string
	Keys []Key
}
