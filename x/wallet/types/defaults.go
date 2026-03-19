package types

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		Wallets:        []Wallet{},
		TransferRecords: []TransferRecord{},
		NextTransferId: 1,
	}
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
	return nil
}

// DefaultParams returns default module parameters
func DefaultParams() Params {
	return Params{}
}

// Validate validates params
func (p Params) Validate() error {
	return nil
}
