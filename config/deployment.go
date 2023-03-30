package configuration

// EthrDidRegistryDeployment represents metadata for a deployment of the ERC1056 registry contract.
type EthrDidRegistryDeployment struct {
	ChainID     int    `json:"chainId"`
	Registry    string `json:"registry"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	RPCURL      string `json:"rpcUrl,omitempty"`
	LegacyNonce bool   `json:"legacyNonce"`
	Additional  map[string]interface{}
}

// Deployments represents the known deployments of the ERC1056 registry contract.
var Deployments = []EthrDidRegistryDeployment{
	{ChainID: 1, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "mainnet", LegacyNonce: true},
	{ChainID: 3, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "ropsten", LegacyNonce: true},
	{ChainID: 4, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "rinkeby", LegacyNonce: true},
	{ChainID: 5, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "goerli", LegacyNonce: true},
	{ChainID: 42, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "kovan", LegacyNonce: true},
	{ChainID: 30, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "rsk", LegacyNonce: true},
	{ChainID: 31, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "rsk:testnet", LegacyNonce: true},
	{ChainID: 246, Registry: "0xE29672f34e92b56C9169f9D485fFc8b9A136BCE4", Name: "ewc", Description: "energy web chain", LegacyNonce: false},
	{ChainID: 73799, Registry: "0xC15D5A57A8Eb0e1dCBE5D88B8f9a82017e5Cc4AF", Name: "volta", LegacyNonce: true},
	{ChainID: 246785, Registry: "0xdCa7EF03e98e0DC2B855bE647C39ABe984fcF21B", Name: "artis:tau1", LegacyNonce: true},
	{ChainID: 246529, Registry: "0xdCa7EF03e98e0DC2B855bE647C39ABe984fcF21B", Name: "artis:sigma1", LegacyNonce: true},
	{ChainID: 137, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "polygon", LegacyNonce: true},
	{ChainID: 80001, Registry: "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b", Name: "polygon:test", LegacyNonce: true},
	{ChainID: 1313161554, Registry: "0x63eD58B671EeD12Bc1652845ba5b2CDfBff198e0", Name: "aurora", LegacyNonce: true},
}
