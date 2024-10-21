// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v6 "github.com/evmos/os/x/evm/migrations/v6"
	"github.com/evmos/os/x/evm/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace types.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, legacySubspace types.Subspace) Migrator {
	return Migrator{
		keeper:         keeper,
		legacySubspace: legacySubspace,
	}
}

// Migrate5to6 migrates the store from consensus version 5 to 6
func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	// As current params's consensus version is 7, we merge the migrate logic
	// Migrate5to6 and Migrate6to7 into this one and make Migrate6to7 empty
	return v6.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate6to7 migrates the store from consensus version 6 to 7
func (m Migrator) Migrate6to7(ctx sdk.Context) error {
	// As current params's consensus version is 7, we merge the migrate logic
	// Migrate5to6 and Migrate6to7 into this one and make Migrate6to7 empty
	return nil
}
