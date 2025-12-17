package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	cli "tax-calculator-cli/internal/cli"
)

func runProcessorWithInput(t *testing.T, stdin []byte) (string, error) {
	t.Helper()
	origIn, origOut := os.Stdin, os.Stdout
	defer func() { os.Stdin, os.Stdout = origIn, origOut }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdin: %v", err)
	}
	os.Stdin = r
	if _, err := w.Write(stdin); err != nil {
		t.Fatalf("write stdin: %v", err)
	}
	_ = w.Close()

	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = pw

	p := cli.NewProcessor()
	runErr := p.Run()

	_ = pw.Close()
	out, readErr := io.ReadAll(pr)
	if readErr != nil {
		t.Fatalf("read stdout: %v", readErr)
	}
	_ = pr.Close()

	return string(out), runErr
}

func readRepoFile(t *testing.T, rel string) []byte {
	t.Helper()
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Dir(filename)
	path := filepath.Join(root, rel)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file %s: %v", path, err)
	}
	return b
}

func normalizeLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return lines
}

func Test_EndToEnd_InputMatchesExpectedOutput(t *testing.T) {
	in := readRepoFile(t, "io/input.txt")
	exp := readRepoFile(t, "io/expected_output.txt")

	gotStr, err := runProcessorWithInput(t, in)
	if err != nil {
		t.Fatalf("processor returned error: %v", err)
	}

	gotLines := normalizeLines(gotStr)
	expLines := normalizeLines(string(exp))

	if len(gotLines) != len(expLines) {
		t.Fatalf("line count mismatch: got %d, want %d\nGOT:\n%s\nWANT:\n%s", len(gotLines), len(expLines), gotStr, string(exp))
	}
	for i := range expLines {
		if strings.TrimSpace(gotLines[i]) != strings.TrimSpace(expLines[i]) {
			t.Fatalf("line %d mismatch\n got: %s\nwant: %s", i+1, gotLines[i], expLines[i])
		}
	}
}

func Test_EndToEnd_SingleLine_NoTrailingNewline(t *testing.T) {
	single := `[{"operation":"buy","unit-cost":10.00,"quantity":100},{"operation":"sell","unit-cost":15.00,"quantity":50},{"operation":"sell","unit-cost":15.00,"quantity":50}]`

	gotStr, err := runProcessorWithInput(t, []byte(single))
	if err != nil {
		t.Fatalf("processor returned error: %v", err)
	}

	lines := normalizeLines(gotStr)
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d: %q", len(lines), gotStr)
	}
	expected := `[{"tax":0},{"tax":0},{"tax":0}]`
	if strings.TrimSpace(lines[0]) != expected {
		t.Fatalf("unexpected output\n got: %s\nwant: %s", lines[0], expected)
	}
}
