package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runMainWithStdio(t *testing.T, stdin []byte) (string, error) {
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

	mainErr := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				// In case something panics, convert to error for test visibility
				err = errFromRecover(r)
			}
		}()
		main()
		return nil
	}()

	_ = pw.Close()
	out, readErr := io.ReadAll(pr)
	if readErr != nil {
		t.Fatalf("read stdout: %v", readErr)
	}
	_ = pr.Close()

	return string(out), mainErr
}

func errFromRecover(r any) error {
	if e, ok := r.(error); ok {
		return e
	}
	return nil
}

func TestMain_SingleLine_NoTrailingNewline(t *testing.T) {
	input := []byte(`[{"operation":"buy","unit-cost":10.00,"quantity":100},{"operation":"sell","unit-cost":15.00,"quantity":50},{"operation":"sell","unit-cost":15.00,"quantity":50}]`)
	out, err := runMainWithStdio(t, input)
	if err != nil {
		t.Fatalf("main returned error: %v", err)
	}
	lines := normalizeLines(out)
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d: %q", len(lines), out)
	}
	expected := `[{"tax":0},{"tax":0},{"tax":0}]`
	if strings.TrimSpace(lines[0]) != expected {
		t.Fatalf("unexpected output\n got: %s\nwant: %s", lines[0], expected)
	}
}

func TestMain_InvalidJSON_ExitsWithNonZero(t *testing.T) {
	cmd := exec.Command("go", "run", ".")
	cmd.Stdin = bytes.NewBufferString("not a json")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err == nil {
		t.Fatalf("expected non-zero exit when JSON is invalid; stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	// Optionally assert that stderr contains the message (log.Fatal writes to stderr)
	if !strings.Contains(stderr.String(), "erro ao parsear JSON") && !strings.Contains(stderr.String(), "erro ao calcular imposto") {
		// depending on where it fails, one of these messages should appear
		t.Fatalf("unexpected stderr: %q", stderr.String())
	}
}
