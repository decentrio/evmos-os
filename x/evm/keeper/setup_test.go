package keeper_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/params"
	testconstants "github.com/evmos/os/testutil/constants"
	"github.com/evmos/os/testutil/integration/os/factory"
	"github.com/evmos/os/testutil/integration/os/grpc"
	"github.com/evmos/os/testutil/integration/os/keyring"
	"github.com/evmos/os/testutil/integration/os/network"
	evmtypes "github.com/evmos/os/x/evm/types"
	feemarkettypes "github.com/evmos/os/x/feemarket/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	network *network.UnitTestNetwork
	handler grpc.Handler
	keyring keyring.Keyring
	factory factory.TxFactory

	enableFeemarket  bool
	enableLondonHF   bool
	mintFeeCollector bool
}

type UnitTestSuite struct {
	suite.Suite
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	s.enableFeemarket = false
	s.enableLondonHF = true
	suite.Run(t, s)

	// Run UnitTestSuite
	unitTestSuite := new(UnitTestSuite)
	suite.Run(t, unitTestSuite)
}

func (suite *KeeperTestSuite) SetupTest() {
	keys := keyring.New(2)
	// Set custom balance based on test params
	customGenesis := network.CustomGenesisState{}
	feemarketGenesis := feemarkettypes.DefaultGenesisState()
	if s.enableFeemarket {
		feemarketGenesis.Params.EnableHeight = 1
		feemarketGenesis.Params.NoBaseFee = false
	} else {
		feemarketGenesis.Params.NoBaseFee = true
	}
	customGenesis[feemarkettypes.ModuleName] = feemarketGenesis

	if !s.enableLondonHF {
		evmGenesis := evmtypes.DefaultGenesisState()
		maxInt := sdkmath.NewInt(math.MaxInt64)
		evmGenesis.Params.EvmDenom = testconstants.ExampleAttoDenom
		evmGenesis.Params.ChainConfig.LondonBlock = &maxInt
		evmGenesis.Params.ChainConfig.ArrowGlacierBlock = &maxInt
		evmGenesis.Params.ChainConfig.GrayGlacierBlock = &maxInt
		evmGenesis.Params.ChainConfig.MergeNetsplitBlock = &maxInt
		evmGenesis.Params.ChainConfig.ShanghaiBlock = &maxInt
		evmGenesis.Params.ChainConfig.CancunBlock = &maxInt
		customGenesis[evmtypes.ModuleName] = evmGenesis
	}

	if s.mintFeeCollector {
		// mint some coin to fee collector
		coins := sdk.NewCoins(sdk.NewCoin(testconstants.ExampleAttoDenom, sdkmath.NewInt(int64(params.TxGas)-1)))
		balances := []banktypes.Balance{
			{
				Address: authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
				Coins:   coins,
			},
		}
		bankGenesis := banktypes.DefaultGenesisState()
		bankGenesis.Balances = balances
		customGenesis[banktypes.ModuleName] = bankGenesis
	}

	nw := network.NewUnitTestNetwork(
		network.WithPreFundedAccounts(keys.GetAllAccAddrs()...),
		network.WithCustomGenesis(customGenesis),
	)
	gh := grpc.NewIntegrationHandler(nw)
	tf := factory.New(nw, gh)

	s.network = nw
	s.factory = tf
	s.handler = gh
	s.keyring = keys
}
