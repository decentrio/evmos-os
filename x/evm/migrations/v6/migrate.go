// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package v6

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/os/x/evm/types"

	v5types "github.com/evmos/os/x/evm/migrations/v6/types"
)

// MigrateStore migrates the x/evm module state from the consensus version 5 to
// version 6. Specifically, it adds the new EVMChannels param.
func MigrateStore(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {
	var (
		paramsV5 v5types.V5Params
		params   types.Params
	)

	store := ctx.KVStore(storeKey)

	paramsV5Bz := store.Get(types.KeyPrefixParams)
	cdc.MustUnmarshal(paramsV5Bz, &paramsV5)

	params.EvmDenom = paramsV5.EvmDenom

	params.ChainConfig = types.ChainConfig{
		HomesteadBlock:      paramsV5.ChainConfig.HomesteadBlock,
		DAOForkBlock:        paramsV5.ChainConfig.DAOForkBlock,
		DAOForkSupport:      paramsV5.ChainConfig.DAOForkSupport,
		EIP150Block:         paramsV5.ChainConfig.EIP150Block,
		EIP150Hash:          paramsV5.ChainConfig.EIP150Hash,
		EIP155Block:         paramsV5.ChainConfig.EIP155Block,
		EIP158Block:         paramsV5.ChainConfig.EIP158Block,
		ByzantiumBlock:      paramsV5.ChainConfig.ByzantiumBlock,
		ConstantinopleBlock: paramsV5.ChainConfig.ConstantinopleBlock,
		PetersburgBlock:     paramsV5.ChainConfig.PetersburgBlock,
		IstanbulBlock:       paramsV5.ChainConfig.IstanbulBlock,
		MuirGlacierBlock:    paramsV5.ChainConfig.MuirGlacierBlock,
		BerlinBlock:         paramsV5.ChainConfig.BerlinBlock,
		LondonBlock:         paramsV5.ChainConfig.LondonBlock,
		ArrowGlacierBlock:   paramsV5.ChainConfig.ArrowGlacierBlock,
		GrayGlacierBlock:    paramsV5.ChainConfig.GrayGlacierBlock,
		MergeNetsplitBlock:  paramsV5.ChainConfig.MergeNetsplitBlock,
		ShanghaiBlock:       paramsV5.ChainConfig.ShanghaiBlock,
		CancunBlock:         paramsV5.ChainConfig.CancunBlock,
	}
	params.AllowUnprotectedTxs = paramsV5.AllowUnprotectedTxs
	params.EVMChannels = types.DefaultEVMChannels
	params.AccessControl = types.DefaultAccessControl
	params.ActiveStaticPrecompiles = types.DefaultStaticPrecompiles

	// Migrate old ExtraEIPs from int64 to string. Since no Evmos EIPs have been
	// created before and activators contains only `ethereum_XXXX` activations,
	// all values will be prefixed with `ethereum_`.
	params.ExtraEIPs = make([]string, 0, len(paramsV5.ExtraEIPs))
	for _, eip := range paramsV5.ExtraEIPs {
		eipName := fmt.Sprintf("ethereum_%d", eip)
		params.ExtraEIPs = append(params.ExtraEIPs, eipName)
	}

	if err := params.Validate(); err != nil {
		return err
	}

	bz := cdc.MustMarshal(&params)

	store.Set(types.KeyPrefixParams, bz)
	return nil
}