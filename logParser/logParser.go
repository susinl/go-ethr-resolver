package logParser

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC1056Event struct {
	Identity       string
	PreviousChange *big.Int
	ValidTo        *big.Int
	EventName      string
	BlockNumber    uint64
}

var contractAbi abi.ABI

func populateEventMetaClass(log *types.Log, blockNumber uint64) (ERC1056Event, error) {
	result := make(map[string]interface{})

	var eventName string
	var event abi.Event
	var found bool
	for _, e := range contractAbi.Events {
		if e.ID == log.Topics[0] {
			eventName = e.Name
			event = e
			found = true
			break
		}
	}

	if !found {
		return ERC1056Event{}, errors.New("event not found in contract ABI")
	}

	values, err := contractAbi.Unpack(event.RawName, log.Data)
	if err != nil {
		return ERC1056Event{}, err
	}

	if len(values) != 2 {
		return ERC1056Event{}, errors.New("malformed event input. wrong number of arguments")
	}

	result["identity"] = values[0].(string)
	result["previousChange"] = values[1].(*big.Int)

	erc1056Event := ERC1056Event{
		Identity:       result["identity"].(string),
		PreviousChange: result["previousChange"].(*big.Int),
		EventName:      eventName,
		BlockNumber:    blockNumber,
	}

	return erc1056Event, nil
}

// the ABI is not exported in the package "github.com/ethereum/go-ethereum/accounts/abi/bind", so it has to be passed as ana dditional argument
func logDecoder(contract *bind.BoundContract, contractAbi *abi.ABI, logs []*types.Log) ([]ERC1056Event, error) {
	results := make([]ERC1056Event, len(logs))

	for i, log := range logs {
		_, err := contractAbi.EventByID(log.Topics[0])
		if err != nil {
			return nil, err
		}

		event, err := populateEventMetaClass(log, log.BlockNumber)
		if err != nil {
			return nil, err
		}

		results[i] = event
	}

	return results, nil
}
