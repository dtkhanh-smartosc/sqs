package cosmwasmpool

import (
	"github.com/osmosis-labs/osmosis/osmomath"
)

const (
	ALLOY_TRANSMUTER_CONTRACT_NAME               = "crates.io:transmuter"
	ALLOY_TRANSMUTER_MIN_CONTRACT_VERSION        = "3.0.0"
	ALLOY_TRANSMUTER_CONTRACT_VERSION_CONSTRAINT = ">= " + ALLOY_TRANSMUTER_MIN_CONTRACT_VERSION
)

func (model *CosmWasmPoolModel) IsAlloyTransmuter() bool {
	return model.ContractInfo.Matches(
		ALLOY_TRANSMUTER_CONTRACT_NAME,
		mustParseSemverConstraint(ALLOY_TRANSMUTER_CONTRACT_VERSION_CONSTRAINT),
	)
}

// Tranmuter Alloyed Data, since v3.0.0
// AssetConfigs is a list of denom and normalization factor pairs including the alloyed denom.
type AlloyTransmuterData struct {
	AlloyedDenom string                  `json:"alloyed_denom"`
	AssetConfigs []TransmuterAssetConfig `json:"asset_configs"`
}

// Configuration for each asset in the transmuter pool
type TransmuterAssetConfig struct {
	// Denom of the asset
	Denom string `json:"denom"`

	// Normalization factor for the asset.
	// [more info](https://github.com/osmosis-labs/transmuter/tree/v3.0.0?tab=readme-ov-file#normalization-factors)
	NormalizationFactor osmomath.Int `json:"normalization_factor"`
}