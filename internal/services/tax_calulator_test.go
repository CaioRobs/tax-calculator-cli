package services_test

import (
	"testing"

	"tax-calculator-cli/internal/domain"
	"tax-calculator-cli/internal/services"
)

func floatEq(a, b float64) bool {
	const eps = 1e-6
	if a == b {
		return true
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < eps
}

func TestCalculate_InvalidOperationReturnsError(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 10, UnitCost: 10.0},
		{Operation: "hold", Quantity: 5, UnitCost: 10.0}, // invalid
	}

	var results []domain.TaxResult
	err := calculator.Calculate(ops, &results)
	if err == nil {
		t.Fatalf("expected error for invalid operation, got nil")
	}
}

func TestCalculate_SellNotTaxableAtThreshold_NoTax(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 1000, UnitCost: 10.0},  // WAP = 10
		{Operation: "sell", Quantity: 1000, UnitCost: 20.0}, // total = 20_000 => NÃO tributável (condição é > 20000)
	}

	var results []domain.TaxResult
	if err := calculator.Calculate(ops, &results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if !floatEq(results[0].Tax, 0.0) {
		t.Errorf("expected first op tax 0.0, got %v", results[0].Tax)
	}

	if !floatEq(results[1].Tax, 0.0) {
		t.Errorf("sell at threshold should not be taxed, got %v", results[1].Tax)
	}
}

func TestCalculate_SellNotTaxableBelowThreshold_NoTax(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 1000, UnitCost: 10.0},  // WAP = 10
		{Operation: "sell", Quantity: 1000, UnitCost: 19.0}, // total = 19_000 => NÃO tributável
	}

	var results []domain.TaxResult
	if err := calculator.Calculate(ops, &results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if !floatEq(results[1].Tax, 0.0) {
		t.Errorf("sell below threshold should not be taxed, got %v", results[1].Tax)
	}
}

func TestCalculate_NonTaxableProfit_DoesNotReduceAccumulatedLoss(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 1000, UnitCost: 10.0},  // WAP = 10
		{Operation: "sell", Quantity: 500, UnitCost: 5.0},   // loss = (5-10)*500 = -2500 -> accumulatedLoss = 2500
		{Operation: "sell", Quantity: 50, UnitCost: 100.0},  // lucro = (100-10)*50 = 4500 > accumulatedLoss; total=5000 <= 20000 -> NÃO tributa e NÃO altera accumulatedLoss (permanece 2500)
		{Operation: "sell", Quantity: 300, UnitCost: 100.0}, // total=30000 > 20000 -> tributa; lucro=90*300=27000; netProfit=27000-2500=24500 -> tax=4900
	}

	var results []domain.TaxResult
	if err := calculator.Calculate(ops, &results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(ops) {
		t.Fatalf("expected %d results, got %d", len(ops), len(results))
	}

	// venda com prejuízo
	if !floatEq(results[1].Tax, 0.0) {
		t.Errorf("expected loss sell tax 0.0, got %v", results[1].Tax)
	}
	// lucro não tributável (<=20k): imposto 0 e accumulatedLoss não reduzido (validado pela próxima venda)
	if !floatEq(results[2].Tax, 0.0) {
		t.Errorf("expected non-taxable profit sell tax 0.0, got %v", results[2].Tax)
	}
	// próxima venda tributável deve abater o mesmo accumulatedLoss (2500), produzindo tax=4900
	expectedTax := 4900.0
	if !floatEq(results[3].Tax, expectedTax) {
		t.Errorf("expected taxable sell tax %v, got %v", expectedTax, results[3].Tax)
	}
}

func TestCalculate_TaxableSellProducesExpectedTax(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 100, UnitCost: 100.0}, // weighted avg = 100
		{Operation: "sell", Quantity: 50, UnitCost: 500.0}, // operationTotal = 25_000 > 20_000, profit = (500-100)*50 = 20_000 -> tax = 4_000
	}

	var results []domain.TaxResult
	if err := calculator.Calculate(ops, &results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if !floatEq(results[0].Tax, 0.0) {
		t.Errorf("expected first op tax 0.0, got %v", results[0].Tax)
	}

	expectedTax := 4000.0
	if !floatEq(results[1].Tax, expectedTax) {
		t.Errorf("expected second op tax %v, got %v", expectedTax, results[1].Tax)
	}
}

func TestCalculate_AccumulatedLossOffsetsFutureTaxes(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 100, UnitCost: 100.0}, // weighted avg = 100
		{Operation: "sell", Quantity: 50, UnitCost: 50.0},  // loss = (50-100)*50 = -2500 -> accumulatedLoss = 2500
		{Operation: "sell", Quantity: 50, UnitCost: 500.0}, // operationTotal = 25_000 > 20_000, profit = 20_000 -> netProfit = 20_000 - 2500 = 17_500 -> tax = 3500
	}

	var results []domain.TaxResult
	if err := calculator.Calculate(ops, &results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if !floatEq(results[0].Tax, 0.0) {
		t.Errorf("expected first op tax 0.0, got %v", results[0].Tax)
	}

	if !floatEq(results[1].Tax, 0.0) {
		t.Errorf("expected second op (loss) tax 0.0, got %v", results[1].Tax)
	}

	expectedTax := 3500.0
	if !floatEq(results[2].Tax, expectedTax) {
		t.Errorf("expected third op tax %v, got %v", expectedTax, results[2].Tax)
	}
}

func TestCalculate_ProfitLessThanAccumulatedLoss_NoTaxAndLossReduces(t *testing.T) {
	calculator := services.NewTaxCalculator()

	ops := []domain.Operation{
		{Operation: "buy", Quantity: 100, UnitCost: 100.0}, // weighted avg = 100
		{Operation: "sell", Quantity: 50, UnitCost: 50.0},  // loss = 2500 -> accumulatedLoss = 2500
		{Operation: "sell", Quantity: 10, UnitCost: 200.0}, // operationProfit = (200-100)*10 = 1000 <= accumulatedLoss -> no tax, accumulatedLoss reduced
	}

	var results []domain.TaxResult
	if err := calculator.Calculate(ops, &results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	if !floatEq(results[2].Tax, 0.0) {
		t.Errorf("expected third op tax 0.0, got %v", results[2].Tax)
	}
}
