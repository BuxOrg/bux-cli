package cmd

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/BuxOrg/bux/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

/*
destination, err := client.NewDestination(
			ctx, rawXPub, utils.ChainExternal, utils.ScriptTypePubKeyHash, false, opts...,
		)
*/

// commands for destination
const destinationCommandName = "destination"
const destinationCommandNew = "new"

// returnDestinationCmd returns the destination command
func returnDestinationCmd(app *App) *cobra.Command {
	return &cobra.Command{
		Use:   destinationCommandName,
		Short: "manage your destinations in BUX",
		Long: color.GreenString(`
________  ___________ ____________________.___ _______      ________________.___________    _______   
\______ \ \_   _____//   _____/\__    ___/|   |\      \    /  _  \__    ___/|   \_____  \   \      \  
 |    |  \ |    __)_ \_____  \   |    |   |   |/   |   \  /  /_\  \|    |   |   |/   |   \  /   |   \ 
 |    '   \|        \/        \  |    |   |   /    |    \/    |    \    |   |   /    |    \/    |    \
 /_______  /_______  /_______  /  |____|   |___\____|__  /\____|__  /____|   |___\_______  /\____|__  /
		 \/        \/        \/                        \/         \/                     \/         \/`) + `
` + color.YellowString(`
This command is for destination (address) related commands.

new: creates a new destination in BUX (`+destinationCommandName+` new <xpub>)
`),
		// Aliases: []string{"address"},
		Example: applicationName + " " + destinationCommandNew + " <xpub>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(destinationCommandName + " requires a subcommand, IE: new, etc.")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Initialize the BUX client
			deferFunc := app.InitializeBUX()
			defer deferFunc()

			// Switch on the subcommand
			if args[0] == destinationCommandNew { // Create a new destination

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

				// Get the metadata if provided
				modelOps := app.bux.DefaultModelOptions()
				if len(args) == 3 {
					modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(args[2])))
				}

				// Create the destination
				var destination *bux.Destination
				destination, err = app.bux.NewDestination(
					context.Background(), args[1], utils.ChainExternal, utils.ScriptTypePubKeyHash, false, modelOps...,
				)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error creating destination: "+err.Error())
					return
				}

				// Display the destination
				displayModel(destination)
			} else {
				chalker.Log(chalker.ERROR, "Unknown subcommand")
			}
		},
	}
}
