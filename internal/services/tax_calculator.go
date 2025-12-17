package services

import (
	"fmt"
	"tax-calculator-cli/internal/domain"
)

type TaxCalculator struct{}

func NewTaxCalculator() *TaxCalculator {
	return &TaxCalculator{}
}

func buyOperation(
	totalSharesQuantity *float64,
	weightedAveragePrice *float64,
	operation domain.Operation,
) domain.TaxResult {
	standardTax := 0.0

	unitCost := operation.UnitCost
	quantity := float64(operation.Quantity)

	newTotalShares := *totalSharesQuantity + quantity
	newWeightedAveragePrice :=
		Round2(((*totalSharesQuantity * *weightedAveragePrice) + (quantity * unitCost)) / newTotalShares)

	*totalSharesQuantity = newTotalShares
	*weightedAveragePrice = newWeightedAveragePrice

	return domain.TaxResult{Tax: &standardTax}
}

func validateStocksQuantity(
	quantity float64,
	totalSharesQuantity float64,
) error {
	if quantity > totalSharesQuantity {
		return fmt.Errorf("can't sell more stocks than you have")
	}
	return nil
}

func sellOperation(
	totalSharesQuantity *float64,
	weightedAveragePrice *float64,
	accumulatedLoss *float64,
	operation domain.Operation,
) domain.TaxResult {
	quantity := float64(operation.Quantity)
	if err := validateStocksQuantity(quantity, *totalSharesQuantity); err != nil {
		return domain.TaxResult{Error: "Can't sell more stocks than you have"}
	}

	minOperationValueForTaxing := 20000.0
	taxRate := 0.20
	unitCost := operation.UnitCost

	operationTotal := unitCost * quantity
	operationProfit := (unitCost - *weightedAveragePrice) * quantity

	tax := 0.0
	isTaxable := operationTotal > minOperationValueForTaxing

	switch {
	case operationProfit < 0:
		*accumulatedLoss += -operationProfit

	case operationProfit <= *accumulatedLoss:
		*accumulatedLoss -= operationProfit

	default:
		if !isTaxable {
			break
		}
		netProfit := operationProfit - *accumulatedLoss
		tax = Round2(netProfit * taxRate)
		*accumulatedLoss = 0
	}

	*totalSharesQuantity -= quantity
	roundedTax := Round2(tax)

	return domain.TaxResult{Tax: &roundedTax}
}

func (t *TaxCalculator) Calculate(operations []domain.Operation, results *[]domain.TaxResult) error {
	tempResults := make([]domain.TaxResult, len(operations))

	totalSharesQuantity := 0.0
	weightedAveragePrice := 0.0
	accumulatedLoss := 0.0

	for i, operation := range operations {
		switch operation.Operation {
		case "buy":
			tempResults[i] = buyOperation(&totalSharesQuantity, &weightedAveragePrice, operation)
		case "sell":
			tempResults[i] = sellOperation(&totalSharesQuantity, &weightedAveragePrice, &accumulatedLoss, operation)
		default:
			return fmt.Errorf("operação invalida: %q", operation.Operation)
		}
	}

	*results = tempResults
	return nil
}
