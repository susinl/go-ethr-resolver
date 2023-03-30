package controller

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/ethclient"
)

// EthrDidController struct
type EthrDidController struct {
	contract    *ethclient.Client
	signer      *ethclient.Client
	address     common.Address
	did         string
	legacyNonce bool
}

// NewEthrDidController creates a new instance of EthrDidController
func NewEthrDidController(identifier string, contract *ethclient.Client, signer *ethclient.Client, chainNameOrID string, rpcURL string, registry string, legacyNonce bool) *EthrDidController {
	var address common.Address
	var publicKey string
	var network string

	// TODO: Add implementation of interpretIdentifier function
	if identifier != "" {
		address, publicKey, network = configuration.interpretIdentifier(identifier)
	}

	if contract != nil {
		return &EthrDidController{
			contract:    contract,
			signer:      signer,
			address:     address,
			did:         fmt.Sprintf("did:ethr:%s%s", network, publicKey),
			legacyNonce: legacyNonce,
		}
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil
	}

	// TODO: Add implementation of getContractForNetwork function

	return &EthrDidController{
		contract:    client,
		signer:      signer,
		address:     address,
		did:         fmt.Sprintf("did:ethr:%s%s", network, publicKey),
		legacyNonce: legacyNonce,
	}
}

// GetOwner returns the owner of the provided address
func (c *EthrDidController) GetOwner(address common.Address, blockTag *big.Int) (common.Address, error) {
	var result common.Address

	// TODO: Add implementation of contract.functions.identityOwner

	return result, nil
}

// AttachContract returns the contract after attaching it to the signer
func (c *EthrDidController) AttachContract(controller *common.Address) (*ethclient.Client, error) {
	var signer *ethclient.Client

	currentOwner, err := c.GetOwner(c.address, nil)
	if err != nil {
		return nil, err
	}

	if controller != nil {
		currentOwner = *controller
	}

	if c.signer != nil {
		signer = c.signer
	} else {
		signer, err = c.contract.GetSigner(currentOwner)
		if err != nil {
			return nil, err
		}
	}

	return c.contract.NewContract(c.contract.Address, signer), nil
}

// ChangeOwner changes the owner of the provided address
func (c *EthrDidController) ChangeOwner(newOwner string, options ...Option) (*types.Transaction, error) {
	var overrides CallOpts

	for _, option := range options {
		overrides = option(overrides)
	}

	contract, err := c.attachContract(overrides.From)
	if err != nil {
		return nil, err
	}

	return contract.Transact("changeOwner", c.Address, newOwner)
}

func (c *EthrDidController) CreateChangeOwnerHash(newOwner string) (string, error) {
	paddedNonce, err := c.getPaddedNonceCompatibility()
	if err != nil {
		return "", err
	}

	dataToHash := append([]byte(MESSAGE_PREFIX), []byte(c.Contract.Address().String())...)
	dataToHash = append(dataToHash, paddedNonce...)
	dataToHash = append(dataToHash, []byte(c.Address)...)
	dataToHash = append(dataToHash, []byte("changeOwner")...)
	dataToHash = append(dataToHash, []byte(newOwner)...)

	return fmt.Sprintf("%x", sha3.Sum256(dataToHash)), nil
}

func (c *EthrDidController) ChangeOwnerSigned(newOwner string, metaSignature MetaSignature, options ...Option) (*types.Transaction, error) {
	var overrides CallOpts

	for _, option := range options {
		overrides = option(overrides)
	}

	contract, err := c.attachContract(overrides.From)
	if err != nil {
		return nil, err
	}

	return contract.Transact("changeOwnerSigned", c.Address, metaSignature.SigV, metaSignature.SigR, metaSignature.SigS, newOwner)
}

type Identity struct {
	contractAddress common.Address
	legacyNonce     bool
	client          Client
}

// MetaSignature represents the Ethereum signature
type MetaSignature struct {
	SigV uint8
	SigR [32]byte
	SigS [32]byte
}

// AddDelegate adds a delegate to the Identity Contract
func (i *Identity) AddDelegate(ctx context.Context, delegateType string, delegateAddress common.Address, exp uint64, opts ...CallOption) (*types.Receipt, error) {
	contract, err := i.attachContract()
	if err != nil {
		return nil, err
	}

	delegateTypeBytes := stringToBytes32(delegateType)
	addDelegateTx, err := contract.AddDelegate(opts, i.contractAddress, delegateTypeBytes, delegateAddress, exp)
	if err != nil {
		return nil, err
	}

	receipt, err := bind.WaitMined(ctx, i.client, addDelegateTx)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// CreateAddDelegateHash creates a hash for the addDelegate function call
func (i *Identity) CreateAddDelegateHash(delegateType string, delegateAddress common.Address, exp uint64) ([]byte, error) {
	paddedNonce, err := i.getPaddedNonceCompatibility()
	if err != nil {
		return nil, err
	}

	dataToHash := append([]byte(MESSAGE_PREFIX), i.contractAddress.Bytes()...)
	dataToHash = append(dataToHash, paddedNonce...)
	dataToHash = append(dataToHash, i.contractAddress.Bytes()...)

	delegateTypeBytes := []byte(delegateType)
	delegateTypeFormatted := formatBytes32String(delegateTypeBytes)
	expBytes := make([]byte, 32)
	binary.LittleEndian.PutUint64(expBytes, exp)

	dataToHash = append(dataToHash, []byte("addDelegate")...)
	dataToHash = append(dataToHash, delegateTypeFormatted...)
	dataToHash = append(dataToHash, delegateAddress.Bytes()...)
	dataToHash = append(dataToHash, expBytes...)

	return crypto.Keccak256(dataToHash), nil
}

func (identity *Identity) AddDelegateSigned(delegateType string, delegateAddress common.Address, exp uint64, metaSignature MetaSignature, options ...contract.Option) (*types.Receipt, error) {
	delegateTypeBytes := stringToBytes32(delegateType)

	transactOpts, err := bind.NewTransactor(strings.NewReader(key), identity.client.GetPassword())
	if err != nil {
		return nil, err
	}

	transactOpts.GasLimit = 123456
	transactOpts = append(transactOpts, options...)

	contract, err := identity.attachContract(transactOpts.From)
	if err != nil {
		return nil, err
	}

	transactOpts.From = nil

	addDelegateTx, err := contract.AddDelegateSigned(transactOpts, identity.Address, metaSignature.SigV, metaSignature.SigR, metaSignature.SigS, delegateTypeBytes, delegateAddress, exp)
	if err != nil {
		return nil, err
	}

	return identity.client.WaitForReceipt(context.Background(), addDelegateTx.Hash())
}
