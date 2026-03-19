package types

// DefaultGenesis 返回默认创世状态
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate 验证创世状态
func (gs GenesisState) Validate() error {
	return nil
}

// NewGenesisState 创建创世状态
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}
