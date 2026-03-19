package types

const (
	ModuleName = "miner"
	StoreKey   = ModuleName
	RouterKey  = ModuleName

	ParamsKey     = "Params-value-"
	MinerKey      = "Miner-value-"
	MinerCountKey = "MinerCount-value-"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GovModuleName is the gov module name
const GovModuleName = "gov"

// DefaultParams returns default params
func DefaultParams() Params {
	return Params{}
}
