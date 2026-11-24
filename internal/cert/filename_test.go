package cert

import (
	"strings"
	"testing"
)

func TestFileNamePartDropsPathCharacters(t *testing.T) {
	got := fileNamePart("../../etc/x y")
	if strings.ContainsAny(got, "/ \\\\") {
		t.Errorf("got %q", got)
	}
	if fileNamePart("rds-mysql_8.0") != "rds-mysql_8.0" {
		t.Error("plain value changed")
	}
}
