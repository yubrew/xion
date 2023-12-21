package feeabs_tests

import (
	"context"
	"github.com/strangelove-ventures/interchaintest/v4"
	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v4/chain/cosmos/wasm"
	"github.com/strangelove-ventures/interchaintest/v4/ibc"
	"github.com/strangelove-ventures/interchaintest/v4/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestXionFeeAbstraction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	ctx := context.Background()

	client, network := interchaintest.DockerSetup(t)

	numVals := 1
	numFullNodes := 1
	nobleGenesisWrapper := genesisWrapper{}

	// Build Chains configs
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:    "xion",
			Version: "local",
			ChainConfig: ibc.ChainConfig{
				GasPrices: "0.0uxion",
			},
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:    "osmosis",
			Version: "v21.0.0",
			ChainConfig: ibc.ChainConfig{
				GasPrices:      "0.005uosmo",
				EncodingConfig: wasm.WasmEncoding(),
			},
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		nobleChainSpec(ctx,
			&nobleGenesisWrapper,
			"noble",
			numVals,
			numFullNodes,
			true, true, true, true,
		),
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	xion, noble, osmosis := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain), chains[2].(*cosmos.CosmosChain)

	// Build relayer
	r := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
	).Build(t, client, network)

	// Build Interchain
	ic := interchaintest.NewInterchain().
		AddChain(xion).
		AddChain(noble).
		AddChain(osmosis).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  xion,
			Chain2:  noble,
			Relayer: r,
			Path:    "xion-noble",
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  xion,
			Chain2:  osmosis,
			Relayer: r,
			Path:    "xion-osmosis",
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  osmosis,
			Chain2:  noble,
			Relayer: r,
			Path:    "osmosis-noble",
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)
	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})
}
