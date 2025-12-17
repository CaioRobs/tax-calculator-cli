package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"tax-calculator-cli/internal/domain"
	"tax-calculator-cli/internal/services"
)

type Processor struct {
	taxCalculator *services.TaxCalculator
}

func NewProcessor() *Processor {
	return &Processor{
		taxCalculator: services.NewTaxCalculator(),
	}
}

func (processor *Processor) Run() error {
	reader := bufio.NewReader(os.Stdin)

	// for each line in stdin
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return fmt.Errorf("erro ao ler entrada: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		var operations []domain.Operation
		if err := json.Unmarshal([]byte(line), &operations); err != nil {
			return fmt.Errorf("erro ao parsear JSON: %w", err)
		}

		var results []domain.TaxResult
		if err := processor.taxCalculator.Calculate(operations, &results); err != nil {
			return fmt.Errorf("erro ao calcular imposto: %w", err)
		}

		output, err := json.Marshal(results)
		if err != nil {
			return fmt.Errorf("erro ao gerar JSON de saída: %w", err)
		}

		fmt.Println(string(output))

		if err == io.EOF {
			break
		}
	}

	return nil
}
