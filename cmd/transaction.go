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
const transactionCommandInfo = "info"

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
				return chalker.Error(transactionCommandName + " requires a subcommand, IE: " + transactionCommandRecord + ", etc.")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Initialize the BUX client
			deferFunc := app.InitializeBUX()
			defer deferFunc()

			// Parse Metadata
			var err error
			if metadata, err = cmd.Flags().GetString(flagMetadata); err != nil {
				displayError(errors.New("error parsing metadata: " + err.Error()))
				return
			}

			// Get the transaction ID
			if txID, err = cmd.Flags().GetString(flagTxID); err != nil {
				displayError(errors.New("error getting txid: " + err.Error()))
				return
			}

			// Get the transaction hex from the flags
			if txHex, err = cmd.Flags().GetString(flagTxHex); err != nil {
				displayError(errors.New("error getting hex: " + err.Error()))
				return
			}

			// Get the optional draft id from the flags
			if draftID, err = cmd.Flags().GetString(flagTxDraftID); err != nil {
				displayError(errors.New("error getting draft id: " + err.Error()))
				return
			}

			// Get the optional woc flag from the flags
			if wocEnabled, err = cmd.Flags().GetBool(flagWoc); err != nil {
				displayError(errors.New("error getting woc flag: " + err.Error()))
				return
			}

			// Switch on the subcommand
			if args[0] == transactionCommandRecord { // record a new transaction

				// Check if xpub is provided
				if len(args) < 2 {
					displayError(ErrXpubIsRequired)
					return
				}

				// Record the transaction
				var tx *bux.Transaction
				tx, err = recordTransaction(context.Background(), app, args[1], draftID, metadata, txID, txHex)
				if err != nil {
					displayError(err)
					return
				}

				// Display the transaction
				displayModel(tx)
			} else if args[0] == transactionCommandInfo { // get transaction info

				// Check if xpub id is provided
				if len(args) < 2 {
					displayError(ErrXpubIDIsRequired)
					return
				}

				// Get the transaction info
				var tx *Transaction
				tx, err = getTransaction(context.Background(), app, args[1], txID, wocEnabled)
				if err != nil {
					displayError(err)
					return
				}

				/*// Run the sync task
				tm := app.bux.Taskmanager()
				err = tm.RunTask(context.Background(), &taskmanager.TaskOptions{
					Arguments: []interface{}{app.bux},
					TaskName:  "sync_transaction_sync",
				})
				if err != nil {
					displayError(err)
					return
				}

				time.Sleep(3 * time.Second)*/

				// Display the transaction
				displayModel(tx)
			} else {
				displayError(ErrUnknownSubcommand)
			}
		},
	}

	// Set the transaction ID flag
	newCmd.Flags().StringVarP(&txID, flagTxID, flagTxIDShort, "", "Transaction ID")

	// Set the transaction hex flag
	newCmd.Flags().StringVarP(&txHex, flagTxHex, flagTxHexShort, "", "Transaction Hex")

	// Set the transaction draft flag
	newCmd.Flags().StringVarP(&txHex, flagTxDraftID, flagTxDraftIDShort, "", "Draft ID (optional)")

	// Set the metadata flag
	newCmd.Flags().StringVarP(&metadata, flagMetadata, flagMetadataShort, "", "Model Metadata")

	// Set the woc flag
	newCmd.Flags().BoolP(flagWoc, flagWocShort, wocEnabled, "Optional flag to use WhatsOnChain for additional transaction data")

	return
}

// recordTransaction records a new transaction
func recordTransaction(ctx context.Context, app *App, xpubKey,
	draftID, metadata, txID, txHex string) (tx *bux.Transaction, err error) {

	// Get the xpub
	var xpub *bux.Xpub
	xpub, err = app.bux.GetXpub(ctx, xpubKey)
	if err != nil {
		return
	} else if xpub == nil {
		displayError(ErrXpubNotFound)
		return
	}

	// Check if txID or txHex is provided
	if len(txHex) == 0 && len(txID) == 0 {
		displayError(ErrTxIDOrHexIsRequired)
		return
	}

	// Check if txID is provided
	if len(txID) > 0 {

		verboseLog(func() {
			chalker.Log(chalker.INFO, "...fetching tx from WOC")
		})

		// Get the transaction hex from the txID using the WhatsOnChain API
		txHex, err = app.bux.Chainstate().WhatsOnChain().GetRawTransactionData(ctx, txID)
		if err != nil {
			return
		}
	}

	// Get the metadata if provided
	modelOps := app.bux.DefaultModelOptions()
	if len(metadata) > 0 {
		modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(metadata)))
	}

	// Record the transaction
	tx, err = app.bux.RecordTransaction(ctx, xpubKey, txHex, draftID, modelOps...)

	return
}

// getTransaction gets a transaction
func getTransaction(ctx context.Context, app *App, xpubID, txID string, wocEnabled bool) (tx *Transaction, err error) {

	// Get the transaction info
	tx = new(Transaction)
	tx.Bux, err = app.bux.GetTransaction(ctx, xpubID, txID)
	if err != nil {
		return
	}

	// Check if WhatsOnChain is enabled
	if wocEnabled {

		verboseLog(func() {
			chalker.Log(chalker.INFO, "...fetching tx from WOC")
		})

		// Get the transaction info from the txHex using the WhatsOnChain API
		tx.WOC, err = app.bux.Chainstate().WhatsOnChain().GetTxByHash(ctx, txID)
		if err != nil {
			return
		}
	}

	return
}
