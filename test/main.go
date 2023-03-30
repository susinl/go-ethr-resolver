package test

// package ethrdidresolver

// import (
// 	"context"
// 	"encoding/base64"
// 	"encoding/hex"
// 	"fmt"
// 	"log"
// 	"math"
// 	"math/big"
// 	"net/url"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/path/to/ethr"
// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/mr-tron/base58"
// )

// type ConfigurationOptions struct {
// 	Networks map[string]NetworkConfiguration
// }

// type NetworkConfiguration struct {
// 	RegistryAddress common.Address
// 	Provider        *ethclient.Client
// }

// type DIDResolver struct {
// 	Networks map[string]NetworkConfiguration
// }

// type DIDResolutionOptions struct {
// 	BlockTag string
// }

// type DIDResolutionResult struct {
// 	ResolutionMetadata DIDResolutionMetadata
// 	DocumentMetadata   DIDDocumentMetadata
// 	Document           *DIDDocument
// }

// type DIDResolutionMetadata struct {
// 	ContentType string
// 	Error       string
// 	Message     string
// }

// type DIDDocumentMetadata struct {
// 	Deactivated   bool
// 	VersionID     string
// 	Updated       string
// 	NextVersionID string
// 	NextUpdate    string
// }

// type DIDDocument struct {
// 	Context           []string
// 	ID                string
// 	VerificationMethod []VerificationMethod
// 	Authentication     []string
// 	AssertionMethod    []string
// 	KeyAgreement       []string
// 	Service            []Service
// }

// type VerificationMethod struct {
// 	ID                string
// 	Type              string
// 	Controller         string
// 	BlockchainAccountID string
// 	PublicKeyHex       string
// 	PublicKeyBase64    string
// 	PublicKeyBase58    string
// 	PublicKeyPem       string
// 	Value             string
// }

// type Service struct {
// 	ID              string `json:"id"`
// 	Type            string `json:"type"`
// 	ServiceEndpoint string `json:"serviceEndpoint"`
// }

// type EthrDidResolver struct {
// 	Contracts      DIDResolver
// 	Provider       *ethclient.Client
// 	Contract       *ethr.EthereumDIDRegistry
// 	NetworkID      string
// 	BlockTag       string
// 	MaxBlockHeight *big.Int
// }

// func NewEthrDidResolver(options ConfigurationOptions) (*EthrDidResolver, error) {
// 	contracts, err := ConfigureResolverWithNetworks(options)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &EthrDidResolver{
// 		Contracts: contracts,
// 	}, nil
// }

// func (r *EthrDidResolver) getOwner(address common.Address, networkID string, blockTag *big.Int) (common.Address, error) {
// 	ethrDidController := EthrDidController{Address: address, Contract: r.Contracts[networkID]}
// 	return ethrDidController.GetOwner(blockTag)
// }

// func (r *EthrDidResolver) previousChange(address common.Address, networkID string, blockTag *big.Int) (*big.Int, error) {
// 	result, err := r.Contracts[networkID].Changed(&bind.CallOpts{BlockNumber: blockTag}, address)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (r *EthrDidResolver) getBlockMetadata(blockHeight *big.Int, networkID string) (BlockMetadata, error) {
// 	block, err := r.Contracts[networkID].Client.BlockByNumber(context.Background(), blockHeight)
// 	if err != nil {
// 		return BlockMetadata{}, err
// 	}

// 	return BlockMetadata{
// 		BlockNumber: block.Number(),
// 		BlockHash:   block.Hash(),
// 		Timestamp:   block.Time(),
// 	}, nil
// }

// func (r *EthrDidResolver) Resolve(did string, options *ResolutionOptions) (DidDocument, error) {
// 	ethrDid, err := ParseEthrDid(did)
// 	if err != nil {
// 		return DidDocument{}, err
// 	}

// 	networkID := ethrDid.GetNetworkID()
// 	if _, ok := r.Contracts[networkID]; !ok {
// 		return DidDocument{}, fmt.Errorf("network %s not supported", networkID)
// 	}

// 	blockTag := options.BlockTag
// 	if blockTag == nil {
// 		blockTag = r.MaxBlockHeight
// 	}

// 	owner, err := r.getOwner(ethrDid.Address, networkID, blockTag)
// 	if err != nil {
// 		return DidDocument{}, err
// 	}

// 	previousChange, err := r.previousChange(ethrDid.Address, networkID, blockTag)
// 	if err != nil {
// 		return DidDocument{}, err
// 	}

// 	blockMetadata, err := r.getBlockMetadata(blockTag, networkID)
// 	if err != nil {
// 		return DidDocument{}, err
// 	}

// 	return DidDocument{
// 		Context:        "https://www.w3.org/ns/did/v1",
// 		ID:             ethrDid.ToString(),
// 		Controller:     owner.Hex(),
// 		BlockMetadata:  blockMetadata,
// 		PreviousChange: previousChange,
// 	}, nil
// }
