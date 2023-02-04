package cmd

import (
	"context"

	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// commands for transaction
const transactionCommandName = "transaction"
const transactionCommandRecord = "record"

// returnTransactionCmd returns the transaction command
func returnTransactionCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   transactionCommandName,
		Short: "manage your transactions in BUX",
		Long: color.GreenString(`
_____________________    _____    _______    _________   _____  ____________________.___________    _______   
\__    ___/\______   \  /  _  \   \      \  /   _____/  /  _  \ \_   ___ \__    ___/|   \_____  \   \      \  
  |    |    |       _/ /  /_\  \  /   |   \ \_____  \  /  /_\  \/    \  \/ |    |   |   |/   |   \  /   |   \ 
  |    |    |    |   \/    |    \/    |    \/        \/    |    \     \____|    |   |   /    |    \/    |    \
  |____|    |____|_  /\____|__  /\____|__  /_______  /\____|__  /\______  /|____|   |___\_______  /\____|__  /
                   \/         \/         \/        \/         \/        \/                      \/         \/`) + `
` + color.YellowString(`
This command is for transaction related commands.

record: records a new transaction in BUX (`+transactionCommandName+` record <xpub> <tx_hex>)
`),
		// Aliases: []string{"tx"},
		Example: applicationName + " " + transactionCommandRecord + " <xpub> <tx_hex>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(transactionCommandName + " requires a subcommand, IE: record, etc.")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Initialize the BUX client
			deferFunc := app.InitializeBUX()
			defer deferFunc()

			// Switch on the subcommand
			if args[0] == transactionCommandRecord { // record a new transaction

				// Check if xpub is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpub is required")
					return
				}

				// Get the xpub
				xpub, err := app.bux.GetXpub(context.Background(), args[1])
				if err != nil {
					chalker.Log(chalker.ERROR, "Error finding xpub: "+err.Error())
					return
				}
				if xpub == nil {
					chalker.Log(chalker.ERROR, "Error: xpub not found")
					return
				}

				// todo: need to get flags for options

				/*// Get the metadata if provided
				modelOps := app.bux.DefaultModelOptions()
				if len(args) == 3 {
					modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(args[2])))
				}*/

				// Display the transaction
				displayModel("")
			} else {
				chalker.Log(chalker.ERROR, "Unknown subcommand")
			}
		},
	}
}
