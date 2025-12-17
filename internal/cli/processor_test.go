package cli_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	cli "tax-calculator-cli/internal/cli"
)

func runWithIO(t *testing.T, stdinContent string) (string, error) {
	t.Helper()

	origStdin := os.Stdin
	origStdout := os.Stdout
	defer func() {
		os.Stdin = origStdin
		os.Stdout = origStdout
	}()

	rIn, wIn, err := os.Pipe()
	if err != nil {
		t.Fatalf("falha ao criar pipe de stdin: %v", err)
	}
	os.Stdin = rIn
	if _, err := io.Copy(wIn, bytes.NewBufferString(stdinContent)); err != nil {
		t.Fatalf("falha ao escrever no stdin simulado: %v", err)
	}
	_ = wIn.Close()

	pr, pw, err := os.Pipe()
	if err != nil {
		t.Fatalf("falha ao criar pipe: %v", err)
	}
	os.Stdout = pw

	processor := cli.NewProcessor()
	runErr := processor.Run()

	_ = pw.Close()
	outBytes, readErr := io.ReadAll(pr)
	if readErr != nil {
		t.Fatalf("falha ao ler stdout capturado: %v", readErr)
	}
	_ = pr.Close()

	return string(outBytes), runErr
}

func TestProcessor_Run_SingleLine(t *testing.T) {
	input := `[{"operation":"buy","unit-cost":10.00,"quantity":100},{"operation":"sell","unit-cost":15.00,"quantity":50},{"operation":"sell","unit-cost":15.00,"quantity":50}]`
	out, err := runWithIO(t, input)
	if err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}

	expected := `[{"tax":0},{"tax":0},{"tax":0}]`
	got := strings.TrimSpace(out)
	if got != expected {
		t.Fatalf("saida inesperada\nexpected: %s\n     got: %s", expected, got)
	}
}

func TestProcessor_Run_MultipleLines(t *testing.T) {
	line1 := `[{"operation":"buy","unit-cost":10.00,"quantity":100},{"operation":"sell","unit-cost":15.00,"quantity":50},{"operation":"sell","unit-cost":15.00,"quantity":50}]`
	line2 := `[{"operation":"buy","unit-cost":10.00,"quantity":10000},{"operation":"sell","unit-cost":20.00,"quantity":5000},{"operation":"sell","unit-cost":5.00,"quantity":5000}]`
	input := line1 + "\n" + line2

	out, err := runWithIO(t, input)
	if err != nil {
		t.Fatalf("Run retornou erro inesperado: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("esperava 2 linhas de saida, obtido %d: %q", len(lines), out)
	}

	expected1 := `[{"tax":0},{"tax":0},{"tax":0}]`
	expected2 := `[{"tax":0},{"tax":10000},{"tax":0}]`
	if strings.TrimSpace(lines[0]) != expected1 {
		t.Fatalf("linha 1 inesperada\nexpected: %s\n     got: %s", expected1, lines[0])
	}
	if strings.TrimSpace(lines[1]) != expected2 {
		t.Fatalf("linha 2 inesperada\nexpected: %s\n     got: %s", expected2, lines[1])
	}
}

func TestProcessor_Run_InvalidJSON(t *testing.T) {
	input := `not a json`
	_, err := runWithIO(t, input)
	if err == nil {
		t.Fatalf("esperava erro de parse de JSON, mas err == nil")
	}
	if !strings.Contains(err.Error(), "erro ao parsear JSON") {
		t.Fatalf("mensagem de erro inesperada: %v", err)
	}
}

func TestProcessor_Run_EmptyInput(t *testing.T) {
	out, err := runWithIO(t, "")
	if err != nil {
		t.Fatalf("Run retornou erro inesperado com entrada vazia: %v", err)
	}
	if strings.TrimSpace(out) != "" {
		t.Fatalf("esperava nenhuma saida, obtido: %q", out)
	}
}
