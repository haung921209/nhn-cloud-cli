package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = old
	b, _ := io.ReadAll(r)
	return string(b)
}

func TestPrintWriteResultFollowsOutputFlag(t *testing.T) {
	oldOut, oldQ := output, query
	defer func() { output, query = oldOut, oldQ }()
	resp := map[string]string{"jobId": "j-1"}

	output, query = "json", ""
	got := captureStdout(t, func() { printWriteResult(resp, "Job ID: %s\n", "j-1") })
	if !strings.Contains(got, `"jobId": "j-1"`) || strings.Contains(got, "Job ID:") {
		t.Errorf("-o json printed %q", got)
	}
	output = "table"
	got = captureStdout(t, func() { printWriteResult(resp, "Job ID: %s\n", "j-1") })
	if got != "Job ID: j-1\n" {
		t.Errorf("text output %q", got)
	}
	got = captureStdout(t, func() { progressf("Fetching...\n") })
	if got != "" {
		t.Errorf("progress went to stdout: %q", got)
	}
}
