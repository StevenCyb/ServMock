package ini

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errReader struct{}

func (e *errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func TestParseBasic(t *testing.T) {
	raw := `
key_global = global

[server]
port = 8080
host = localhost

[database]
user = root
password = secret
`

	sections, err := Parse(strings.NewReader(raw), false)
	assert.NoError(t, err)
	assert.Len(t, sections, 3)
	assert.Equal(t, "default", sections[0].Name)
	assert.Len(t, sections[0].Keys, 1)
	assert.Equal(t, "key_global", sections[0].Keys[0].Name)
	assert.Equal(t, "global", sections[0].Keys[0].Value)
	srv := sections[1]
	assert.Equal(t, "server", srv.Name)
	assert.Len(t, srv.Keys, 2)
	assert.Equal(t, "port", srv.Keys[0].Name)
	assert.Equal(t, "8080", srv.Keys[0].Value)
	assert.Equal(t, "host", srv.Keys[1].Name)
	assert.Equal(t, "localhost", srv.Keys[1].Value)
	db := sections[2]
	assert.Equal(t, "database", db.Name)
	assert.Len(t, db.Keys, 2)
	assert.Equal(t, "user", db.Keys[0].Name)
	assert.Equal(t, "root", db.Keys[0].Value)
	assert.Equal(t, "password", db.Keys[1].Name)
	assert.Equal(t, "secret", db.Keys[1].Value)
}

func TestGlobalOnly(t *testing.T) {
	raw := `
foo=bar
baz = qux
`
	sections, err := Parse(strings.NewReader(raw), false)
	assert.NoError(t, err)
	assert.Len(t, sections, 1)
	assert.Equal(t, "default", sections[0].Name)
	assert.Len(t, sections[0].Keys, 2)
	assert.Equal(t, "foo", sections[0].Keys[0].Name)
	assert.Equal(t, "bar", sections[0].Keys[0].Value)
	assert.Equal(t, "baz", sections[0].Keys[1].Name)
	assert.Equal(t, "qux", sections[0].Keys[1].Value)
}

func TestCommentsAndBlank(t *testing.T) {
	raw := `
; comment1
# comment2

[sec]
key=value

; trailing comment
`
	sections, err := Parse(strings.NewReader(raw), false)
	assert.NoError(t, err)
	assert.Len(t, sections, 2)
	assert.Equal(t, "default", sections[0].Name)
	assert.Len(t, sections[0].Keys, 0)
	assert.Equal(t, "sec", sections[1].Name)
	assert.Len(t, sections[1].Keys, 1)
	assert.Equal(t, "key", sections[1].Keys[0].Name)
	assert.Equal(t, "value", sections[1].Keys[0].Value)
}

func TestEmptySectionNameError(t *testing.T) {
	raw := `[]
key=val
`
	_, err := Parse(strings.NewReader(raw), false)
	assert.Equal(t, EmptySectionNameError, err)
}

func TestEmptyKeyError(t *testing.T) {
	raw := `
[sec]
 =value
`
	_, err := Parse(strings.NewReader(raw), false)
	assert.Error(t, err)
	assert.EqualError(t, err, "empty key in section [sec]")
}

func TestScannerError(t *testing.T) {
	_, err := Parse(&errReader{}, false)
	assert.Error(t, err)
	assert.EqualError(t, err, "read error")
}

func TestDuplicateSection(t *testing.T) {
	raw := `
[dup]
a=1
[dup]
b=2
`
	sections, err := Parse(strings.NewReader(raw), false)
	assert.NoError(t, err)
	assert.Len(t, sections, 2)
	dup := sections[1]
	assert.Equal(t, "dup", dup.Name)
	assert.Len(t, dup.Keys, 2)
	assert.Equal(t, "a", dup.Keys[0].Name)
	assert.Equal(t, "1", dup.Keys[0].Value)
	assert.Equal(t, "b", dup.Keys[1].Name)
	assert.Equal(t, "2", dup.Keys[1].Value)
}

func TestDuplicateSectionKeepAsDuplicate(t *testing.T) {
	raw := `
[dup]
a=1
[dup]
b=2
`
	sections, err := Parse(strings.NewReader(raw), true)
	assert.NoError(t, err)
	assert.Len(t, sections, 3)
	dup := sections[1]
	dup2 := sections[2]
	assert.Equal(t, "dup", dup.Name)
	assert.Equal(t, "dup", dup2.Name)
	assert.Len(t, dup.Keys, 1)
	assert.Len(t, dup2.Keys, 1)
	assert.Equal(t, "a", dup.Keys[0].Name)
	assert.Equal(t, "1", dup.Keys[0].Value)
	assert.Equal(t, "b", dup2.Keys[0].Name)
	assert.Equal(t, "2", dup2.Keys[0].Value)
}

func TestIndexing(t *testing.T) {
	raw := `;default section
x=y

[dup]
a=1

[dup]
b=2
`
	sections, err := Parse(strings.NewReader(raw), true)
	assert.NoError(t, err)
	assert.Len(t, sections, 3)
	assert.Equal(t, 0, int(sections[0].LineIndex))
	assert.Equal(t, 4, int(sections[1].LineIndex))
	assert.Equal(t, 7, int(sections[2].LineIndex))
}
