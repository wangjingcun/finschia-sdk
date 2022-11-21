package server

import (
	"fmt"

	"github.com/spf13/cobra"

	ostcmd "github.com/line/ostracon/cmd/ostracon/commands"

	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/server/types"
)

// NewRollbackCmd creates a command to rollback tendermint and multistore state by one height.
func NewRollbackCmd(appCreator types.AppCreator, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "rollback lbm-sdk and tendermint state by one height",
		Long: `
A state rollback is performed to recover from an incorrect application state transition,
when Tendermint has persisted an incorrect app hash and is thus unable to make
progress. Rollback overwrites a state at height n with the state at height n - 1.
The application also roll back to height n - 1. No blocks are removed, so upon
restarting Tendermint the transactions in block n will be re-executed against the
application.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := GetServerContextFromCmd(cmd)
			cfg := ctx.Config
			home := cfg.RootDir
			db, err := openDB(home)
			if err != nil {
				return err
			}
			app := appCreator(ctx.Logger, db, nil, ctx.Viper)
			// rollback tendermint state
			height, hash, err := ostcmd.RollbackState(ctx.Config)
			if err != nil {
				return fmt.Errorf("failed to rollback tendermint state: %w", err)
			}
			// rollback the multistore

			if err := app.CommitMultiStore().RollbackToVersion(height); err != nil {
				return fmt.Errorf("failed to rollback to version: %w", err)
			}

			fmt.Printf("Rolled back state to height %d and hash %X", height, hash)
			return nil
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	return cmd
}
