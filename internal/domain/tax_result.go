package domain

type TaxResult struct {
	Tax   *float64 `json:"tax,omitempty"`
	Error string   `json:"error,omitempty"`
}
