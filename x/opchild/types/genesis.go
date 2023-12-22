package types

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/address"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// NewGenesisState creates a new GenesisState instance
func NewGenesisState(params Params, validators []Validator) *GenesisState {
	return &GenesisState{
		Params:              params,
		LastValidatorPowers: []LastValidatorPower{},
		Validators:          validators,
		Exported:            false,
	}
}

// DefaultGenesisState gets the raw genesis raw message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:               DefaultParams(),
		LastValidatorPowers:  []LastValidatorPower{},
		Validators:           []Validator{},
		Exported:             false,
		NextL2Sequence:       1,
		FinalizedL1Sequences: []uint64{},
	}
}

// ValidateGenesis performs basic validation of rollup genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data *GenesisState, ac address.Codec) error {
	if err := validateGenesisStateValidators(data.Validators); err != nil {
		return err
	}

	return data.Params.Validate(ac)
}

func validateGenesisStateValidators(validators []Validator) error {
	addrMap := make(map[string]bool, len(validators))

	for i := 0; i < len(validators); i++ {
		val := validators[i]
		consPk, err := val.ConsPubKey()
		if err != nil {
			return err
		}

		strKey := string(consPk.Bytes())

		if _, ok := addrMap[strKey]; ok {
			consAddr, err := val.GetConsAddr()
			if err != nil {
				return err
			}
			return fmt.Errorf("duplicate validator in genesis state: moniker %v, address %v", val.Moniker, consAddr)
		}

		addrMap[strKey] = true
	}

	return nil
}

// GetGenesisStateFromAppState returns x/opchild GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.Validators {
		if err := g.Validators[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
