package helper

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	identifierMatcher       = `^(.*)?(0x[0-9a-fA-F]{40}|0x[0-9a-fA-F]{66})$`
	nullAddress             = "0x0000000000000000000000000000000000000000"
	defaultRegistryAddress  = "0xdca7ef03e98e0dc2b855be647c39abe984fcf21b"
	defaultJSONRPC          = "http://127.0.0.1:8545/"
	messagePrefix           = "0x1900"
	verificationMethodTypes = "VerificationMethodTypes"
	eventNames              = "EventNames"
	legacyAttrTypes         = "LegacyAttrTypes"
	legacyAlgoMap           = "LegacyAlgoMap"
	errors                  = "Errors"
)

type Address = string
type Uint256 = *big.Int
type Bytes32 = string
type Bytes = string

type ERC1056Event struct {
	Identity       Address
	PreviousChange Uint256
	ValidTo        *Uint256
	EventName      string
	BlockNumber    uint64
}

type DIDOwnerChanged struct {
	ERC1056Event
	Owner Address
}

type DIDAttributeChanged struct {
	ERC1056Event
	Name    Bytes32
	Value   Bytes
	ValidTo Uint256
}

type DIDDelegateChanged struct {
	ERC1056Event
	DelegateType Bytes32
	Delegate     Address
	ValidTo      Uint256
}

type VerificationMethodTypes struct {
	EcdsaSecp256k1VerificationKey2019 string
	EcdsaSecp256k1RecoveryMethod2020  string
	Ed25519VerificationKey2018        string
	RSAVerificationKey2018            string
	X25519KeyAgreementKey2019         string
}

type EventNames struct {
	DIDOwnerChanged     string
	DIDAttributeChanged string
	DIDDelegateChanged  string
}

type LegacyAttrTypes struct {
	sigAuth string
	veriKey string
	enc     string
}

type LegacyAlgoMap struct {
	Secp256k1VerificationKey2018         string
	Ed25519SignatureAuthentication2018   string
	Secp256k1SignatureAuthentication2018 string
	RSAVerificationKey2018               string
	Ed25519VerificationKey2018           string
	X25519KeyAgreementKey2019            string
}

type MetaSignature struct {
	SigV uint8
	SigR Bytes32
	SigS Bytes32
}

func strip0x(input string) string {
	if strings.HasPrefix(input, "0x") {
		return input[2:]
	}
	return input
}

func bytes32toString(input Bytes32) string {
	inputBytes, _ := hex.DecodeString(strip0x(input))
	return strings.TrimRight(string(inputBytes), "\x00")
}

func stringToBytes32(str string) string {
	strBytes := []byte(str)[:32]
	return fmt.Sprintf("0x%s%s", hex.EncodeToString(strBytes), strings.Repeat("0", 66-len(strBytes)*2))
}

func interpretIdentifier(identifier string) (Address, string, string) {
	id := identifier
	network := ""

	if strings.HasPrefix(id, "did:ethr") {
		components := strings.Split(id, ":")
		id = components[len(components)-1]
		if len(components) >= 4 {
			network = strings.Join(components[2:len(components)-1], ":")
		}
	}

	if len(id) > 42 {
		return common.BytesToAddress(crypto.Keccak256(common.FromHex(id)[1:])).Hex(), id, network
	} else {
		return common.HexToAddress(id).Hex(), "", network // checksum address
	}
}

// This function is a placeholder and you need to replace it with the appropriate implementation
func signMetaTxData(identity string, signerAddress string, privateKeyBytes []byte, dataBytes []byte, didReg interface{}) (*big.Int, *big.Int, error) {
	// Replace the following line with your implementation for fetching the nonce
	nonce := big.NewInt(0)

	paddedNonce := common.LeftPadBytes(nonce.Bytes(), 32)
	dataToSign := append([]byte(messagePrefix), paddedNonce...)
	dataToSign = append(dataToSign, []byte(identity)...)
	dataToSign = append(dataToSign, dataBytes...)

	hash := crypto.Keccak256(dataToSign)

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	if err != nil {
		return nil, nil, err
	}

	return r, s, nil
}

type Errors string

const (
	// The resolver has failed to construct the DID document.
	// This can be caused by a network issue, a wrong registry address or malformed logs while parsing the registry history.
	// Please inspect the `DIDResolutionMetadata.message` to debug further.
	NotFound Errors = "notFound"

	// The resolver does not know how to resolve the given DID. Most likely it is not a `did:ethr`.
	InvalidDid Errors = "invalidDid"

	// The resolver is misconfigured or is being asked to resolve a DID anchored on an unknown network
	UnknownNetwork Errors = "unknownNetwork"
)
