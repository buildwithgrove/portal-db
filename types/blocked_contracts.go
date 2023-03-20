package types

type (
	// Represents global blocked addresses across the entire Portal
	GlobalBlockedContracts struct {
		BlockedAddresses map[BlockedAddress]struct{} `json:"blockedAddresses"`
	}
)

func (g *GlobalBlockedContracts) IsContractBlocked(contract BlockedAddress) bool {
	_, ok := g.BlockedAddresses[contract]
	return ok
}
