package cmd

import (
	"context"
	"errors"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// commands for transaction
const transactionCommandName = "transaction"
const transactionCommandRecord = "record"

// returnTransactionCmd returns the transaction command
func returnTransactionCmd(app *App) (newCmd *cobra.Command) {

	newCmd = &cobra.Command{
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

record: records a new transaction in BUX (`+transactionCommandName+` record <xpub> -i=<tx_id>)
`),
		// Aliases: []string{"tx"},
		Example: applicationName + " " + transactionCommandRecord + " <xpub> -i=<tx_id>",
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

			// Parse Metadata
			var err error
			metadata, err = cmd.Flags().GetString("metadata")
			if err != nil {
				displayError(errors.New("error parsing metadata: " + err.Error()))
				return
			}

			// Switch on the subcommand
			if args[0] == transactionCommandRecord { // record a new transaction

				// Check if xpub is provided
				if len(args) < 2 {
					displayError(ErrXpubIsRequired)
					return
				}

				// Get the xpub
				var xpub *bux.Xpub
				xpub, err = app.bux.GetXpub(context.Background(), args[1])
				if err != nil {
					displayError(errors.New("error finding xpub: " + err.Error()))
					return
				} else if xpub == nil {
					displayError(ErrXpubNotFound)
					return
				}

				// Get the transaction ID
				txID, err = cmd.Flags().GetString("txid")
				if err != nil {
					displayError(errors.New("error getting txid: " + err.Error()))
					return
				}

				// Check if txID is provided
				if len(txID) > 0 {

					// Get the transaction hex from the txID using the WhatsOnChain API
					txHex, err = app.bux.Chainstate().WhatsOnChain().GetRawTransactionData(context.Background(), txID)
					if err != nil {
						displayError(errors.New("error finding transaction: " + err.Error()))
						return
					}
				}

				// Check if txHex is provided
				if len(txHex) <= 0 {

					// Get the transaction hex from the flags
					txHex, err = cmd.Flags().GetString("hex")
					if err != nil {
						displayError(errors.New("error getting hex: " + err.Error()))
						return
					}
					if len(txHex) <= 0 {
						displayError(ErrTxIDOrHexIsRequired)
						return
					}
				}

				// Get the metadata if provided
				modelOps := app.bux.DefaultModelOptions()
				if len(metadata) > 0 {
					modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(metadata)))
				}

				// Record the transaction
				var tx *bux.Transaction
				tx, err = app.bux.RecordTransaction(context.Background(), args[1], txHex, "", modelOps...)
				if err != nil {
					displayError(errors.New("error recording transaction: " + err.Error()))
					return
				}

				// Display the transaction
				displayModel(tx)
			} else {
				displayError(ErrUnknownSubcommand)
			}
		},
	}

	// Set the transaction ID flag
	newCmd.Flags().StringVarP(&txID, "txid", "i", "", "Transaction ID")

	// Set the transaction hex flag
	newCmd.Flags().StringVarP(&txHex, "hex", "x", "", "Transaction Hex")

	// Set the metadata flag
	newCmd.Flags().StringVarP(&metadata, "metadata", "m", "", "Model Metadata")

	return
}
