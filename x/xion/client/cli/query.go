package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/burnt-labs/xion/x/xion/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
)

const (
	FlagDenom = "denom"
)

// GetQueryCmd returns the parent command for all x/xion CLi query commands. The
// provided clientCtx should have, at a minimum, a verifier, Tendermint RPC client,
// and marshaler set.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the xion module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetPlatformPercentageCmd(),
	)

	return cmd
}

func GetPlatformPercentageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "platform-percentage",
		Short: "Query for the  platform percentage",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the governance defined platform percentage for value sends.

Example:
  $ %s query %s platform-percentage
`,
				version.AppName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			ctx := cmd.Context()

			res, err := queryClient.PlatformPercentage(ctx, &types.QueryPlatformPercentageRequest{})
			if err != nil {
				return err
			}

			err = clientCtx.PrintString(fmt.Sprintf("%d", res.PlatformPercentage))
			if err != nil {
				return err
			}

			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
