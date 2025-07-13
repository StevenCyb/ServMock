package ini

import (
	"bufio"
	"io"
	"strings"
)

// Parse reads INI data from the reader and returns sections in order of appearance.
// If `allowDuplicated=true` allows multiple sections with the same name
// else duplicate section headers merge into the first occurrence.
func Parse(r io.Reader, allowDuplicated bool) ([]Section, error) {
	scanner := bufio.NewScanner(r)

	// Initialize with global (default) section
	sections := []Section{{Name: "default", Keys: nil}}
	current := &sections[0]

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			name := strings.TrimSpace(line[1 : len(line)-1])
			if name == "" {
				return nil, EmptySectionNameError
			}
			// Find existing or append new
			found := false
			if !allowDuplicated {
				for i := range sections {
					if sections[i].Name == name {
						current = &sections[i]
						found = true
						break
					}
				}
			}

			if !found {
				sections = append(sections, Section{Name: name, Keys: nil})
				current = &sections[len(sections)-1]
			}
			continue
		}

		// Key=Value line
		if idx := strings.Index(line, "="); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			if key == "" {
				return nil, &EmptyKeyError{SectionName: current.Name}
			}
			current.Keys = append(current.Keys, Key{Name: key, Value: val})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sections, nil
}
