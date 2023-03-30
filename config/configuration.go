package configuration

import (
	"errors"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/ethclient"
)

var infuraNames = map[string]string{
	"polygon":      "matic",
	"polygon:test": "maticmum",
	"aurora":       "aurora-mainnet",
}

var knownInfuraNames = []string{"mainnet", "ropsten", "rinkeby", "goerli", "kovan", "aurora"}

type ProviderConfiguration struct {
	Name     string `json:"name"`
	Registry string `json:"registry"`
	RPCURL   string `json:"rpcUrl"`
	ChainID  string `json:"chainId"`
}

type MultiProviderConfiguration struct {
	ProviderConfiguration
	Networks []ProviderConfiguration `json:"networks"`
}

type InfuraConfiguration struct {
	InfuraProjectID string `json:"infuraProjectId"`
}

type ConfigurationOptions interface{}

type ConfiguredNetworks map[string]*ethclient.Client

func configureNetworksWithInfura(projectID string) (ConfiguredNetworks, error) {
	if projectID == "" {
		return ConfiguredNetworks{}, nil
	}

	networks := ConfiguredNetworks{}
	for _, n := range knownInfuraNames {
		infuraName, ok := infuraNames[n]
		if !ok {
			infuraName = n
		}
		rpcURL := "https://" + infuraName + ".infura.io/v3/" + projectID
		client, err := ethclient.Dial(rpcURL)
		if err != nil {
			return nil, err
		}
		networks[n] = client

	}
	return networks, nil
}

func getContractForNetwork(conf ProviderConfiguration) (*ethclient.Client, error) {
	client, err := ethclient.Dial(conf.RPCURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func configureNetwork(net ProviderConfiguration) ConfiguredNetworks {
	networks := ConfiguredNetworks{}
	chainID, ok := math.ParseBig256(net.ChainID)
	if ok {
		if net.RPCURL != "" {
			contract, err := getContractForNetwork(net)
			if err != nil {
				return networks
			}
			networks[chainID.String()] = contract
		}
	} else if net.RPCURL != "" {
		contract, err := getContractForNetwork(net)
		if err != nil {
			return networks
		}
		networks[net.Name] = contract
	}

	return networks
}

func configureNetworks(conf MultiProviderConfiguration) ConfiguredNetworks {
	networks := ConfiguredNetworks{}

	for _, net := range conf.Networks {
		for k, v := range configureNetwork(net) {
			networks[k] = v
		}
	}

	return networks
}

func configureResolverWithNetworks(conf ConfigurationOptions) (ConfiguredNetworks, error) {
	networks := ConfiguredNetworks{}

	infuraConf, ok := conf.(InfuraConfiguration)
	if ok {
		infuraNetworks, err := configureNetworksWithInfura(infuraConf.InfuraProjectID)
		if err != nil {
			return nil, err
		}

		for k, v := range infuraNetworks {
			networks[k] = v
		}
	}

	multiProviderConf, ok := conf.(MultiProviderConfiguration)
	if ok {
		for k, v := range configureNetworks(multiProviderConf) {
			networks[k] = v
		}
	}

	if len(networks) == 0 {
		return nil, errors.New("invalid_config: Please make sure to have at least one network")
	}

	return networks, nil
}
