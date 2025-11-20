package pricing

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed pricing_data.json
var bundledPricingData []byte

// LoadBundled loads the bundled pricing data
func LoadBundled() (*Data, error) {
	var data Data
	if err := json.Unmarshal(bundledPricingData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse bundled pricing data: %w", err)
	}
	return &data, nil
}
