package cmd

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// commands for xpub
const xpubCommandName = "xpub"
const xpubCommandCreate = "create"

// returnXpubCmd returns the xpub command
func returnXpubCmd(app *App) *cobra.Command {
	return &cobra.Command{
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

create: creates a new xpub in BUX (`+xpubCommandName+` create <xpriv>)
`),
		// Aliases: []string{"hdkey"},
		Example: applicationName + " " + xpubCommandName + " create <xpriv>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(xpubCommandName + " requires a subcommand, IE: create, etc.")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Initialize the BUX client
			deferFunc := app.InitializeBUX()
			defer deferFunc()

			// Switch on the subcommand
			if args[0] == xpubCommandCreate { // Create a new xpub

				// Check if xpriv is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpriv is required")
					return
				}

				// Generate the HDKey from the xpriv
				hdKey, err := bitcoin.GenerateHDKeyFromString(args[1])
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

				// Create the xpub in BUX
				var xpub *bux.Xpub
				xpub, err = app.bux.NewXpub(context.Background(), xPubString, app.bux.DefaultModelOptions()...)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error in BUX: "+err.Error())
					return
				}

				chalker.Log(chalker.INFO, "xpub: "+xPubString)
				chalker.Log(chalker.INFO, "xpub BUX id: "+xpub.ID)
			} else {
				chalker.Log(chalker.ERROR, "Unknown subcommand")
			}
		},
	}
}
