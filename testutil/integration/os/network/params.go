// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package network

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	evmtypes "github.com/evmos/os/x/evm/types"
	feemarketypes "github.com/evmos/os/x/feemarket/types"
)

func (n *IntegrationNetwork) UpdateEvmParams(params evmtypes.Params) error {
	return n.app.EVMKeeper.SetParams(n.ctx, params)
}

func (n *IntegrationNetwork) UpdateFeeMarketParams(params feemarketypes.Params) error {
	return n.app.FeeMarketKeeper.SetParams(n.ctx, params)
}

func (n *IntegrationNetwork) UpdateGovParams(params govtypes.Params) error {
	return n.app.GovKeeper.SetParams(n.ctx, params)
}