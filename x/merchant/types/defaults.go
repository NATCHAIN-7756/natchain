package types

func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:         DefaultParams(),
		Merchants:      []Merchant{},
		PaymentRecords: []PaymentRecord{},
		NextPaymentId:  1,
	}
}

func (gs GenesisState) Validate() error {
	return nil
}

func DefaultParams() Params {
	return Params{
		MaxMerchants:       1000,
		MinCommissionRate:  0,
		MaxCommissionRate:  10000,
	}
}

func (p Params) Validate() error {
	return nil
}
