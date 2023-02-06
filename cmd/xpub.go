package cmd

import (
	"context"
	"encoding/json"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/BuxOrg/bux/utils"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/fatih/color"
	"github.com/libsv/go-bk/bip32"
	"github.com/spf13/cobra"
)

// commands for xpub
const xpubCommandGet = "get"
const xpubCommandName = "xpub"
const xpubCommandNew = "new"

// returnXpubCmd returns the xpub command
func returnXpubCmd(app *App) (newCmd *cobra.Command) {
	newCmd = &cobra.Command{
		Use:   xpubCommandName,
		Short: "manage your xpubs in BUX",
		Long: color.GreenString(`
____  _____________ ____ _____________ 
\   \/  /\______   \    |   \______   \
 \     /  |     ___/    |   /|    |  _/
 /     \  |    |   |    |  / |    |   \
/___/\  \ |____|   |______/  |______  /
      \_/                           \/`) + `
` + color.YellowString(`
This command is for xpub (HD-Key) related commands.

new: creates a new xpub in BUX (`+xpubCommandName+` new <xpriv>)
get: get a xpub from BUX (`+xpubCommandName+` get <xpub> | <xpub_id> -m=<metadata_json>)
`),
		// Aliases: []string{"hdkey"},
		Example: applicationName + " " + xpubCommandName + " " + xpubCommandNew + " <xpriv>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(xpubCommandName + " requires a subcommand, IE: " + xpubCommandNew + ", etc.")
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
				chalker.Log(chalker.ERROR, "Error getting metadata: "+err.Error())
				return
			}

			// Switch on the subcommand
			if args[0] == xpubCommandNew { // Create a new xpub

				// Check if xpriv is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpriv is required")
					return
				}

				// Generate the HDKey from the xpriv
				var hdKey *bip32.ExtendedKey
				hdKey, err = bitcoin.GenerateHDKeyFromString(args[1])
				if err != nil {
					chalker.Log(chalker.ERROR, "Error generating: "+err.Error()+", using: "+args[1])
					return
				}

				// Get the xpub
				var xPubString string
				xPubString, err = bitcoin.GetExtendedPublicKey(hdKey)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error deriving xpub: "+err.Error())
					return
				}

				// Get the metadata if provided
				modelOps := app.bux.DefaultModelOptions()
				if len(metadata) > 0 {
					modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(metadata)))
				}

				// Create the xpub in BUX
				var xpub *bux.Xpub
				xpub, err = app.bux.NewXpub(context.Background(), xPubString, modelOps...)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error in BUX: "+err.Error())
					return
				}

				chalker.Log(chalker.INFO, "xpub created: "+xPubString)
				displayModel(xpub)
			} else if args[0] == xpubCommandGet { // Get a xpub from BUX

				// Check if xpub or xpub id is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpub or xpub_id is required")
					return
				}

				// Get the xpub from BUX
				var xpub *bux.Xpub
				if _, err = utils.ValidateXPub(args[1]); err == nil {

					// Get the xpub by xpub
					if xpub, err = app.bux.GetXpub(context.Background(), args[1]); err != nil {
						chalker.Log(chalker.ERROR, "Error getting xpub: "+err.Error())
						return
					}

					// Display the xpub
					displayModel(xpub)
				} else if len(metadata) > 0 {

					// Unmarshal the metadata
					metaData := new(bux.Metadata)
					if err = json.Unmarshal([]byte(metadata), &metaData); err != nil {
						chalker.Log(chalker.ERROR, "Error unmarshalling metadata: "+err.Error())
						return
					}

					// Get the xpubs from BUX
					var xpubs []*bux.Xpub
					if xpubs, err = app.bux.GetXPubs(context.Background(), metaData, nil, nil); err != nil {
						chalker.Log(chalker.ERROR, "Error getting xpubs: "+err.Error())
						return
					}

					// Check if any xpubs were found
					if len(xpubs) == 0 {
						chalker.Log(chalker.ERROR, "No xpubs found")
						return
					}

					// Display the xpubs
					for _, xpub = range xpubs {
						displayModel(xpub)
					}
				} else {

					// Get the xpub from BUX by id
					if xpub, err = app.bux.GetXpubByID(context.Background(), args[1]); err != nil {
						chalker.Log(chalker.ERROR, "Error getting xpub by id: "+err.Error())
						return
					}

					// Display the xpub
					displayModel(xpub)
				}
			} else {
				chalker.Log(chalker.ERROR, "Unknown subcommand")
			}
		},
	}

	// Set the metadata flag
	newCmd.Flags().StringVarP(&metadata, "metadata", "m", "", "Model Metadata")

	return
}
