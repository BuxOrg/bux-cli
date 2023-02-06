package cmd

import (
	"context"
	"errors"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/BuxOrg/bux/utils"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// commands for destination
const destinationCommandName = "destination"
const destinationCommandNew = "new"
const destinationCommandGet = "get"

// returnDestinationCmd returns the destination command
func returnDestinationCmd(app *App) (newCmd *cobra.Command) {
	newCmd = &cobra.Command{
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
get: gets an existing destination in BUX (`+destinationCommandName+` get <destination_id | address | locking_script> <xpub_id>)
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

			// Parse Metadata
			var err error
			metadata, err = cmd.Flags().GetString("metadata")
			if err != nil {
				chalker.Log(chalker.ERROR, "Error getting metadata: "+err.Error())
				return
			}

			// Switch on the subcommand
			if args[0] == destinationCommandNew { // Create a new destination

				// Check if xpub is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpub is required")
					return
				}

				// Get the xpub
				var xpub *bux.Xpub
				xpub, err = app.bux.GetXpub(context.Background(), args[1])
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
				if len(metadata) > 0 {
					modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(metadata)))
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
			} else if args[0] == destinationCommandGet { // Get a destination

				// Check if destination ID is provided
				if len(args) < 3 {
					chalker.Log(chalker.ERROR, "Error: (destination id, address or locking script) and xpub_id is required")
					return
				}

				// Get the destination by ID, address or locking script
				var destination *bux.Destination
				destination, err = app.bux.GetDestinationByID(context.Background(), args[2], args[1])
				if err != nil && !errors.Is(err, bux.ErrMissingDestination) {
					chalker.Log(chalker.ERROR, "Error finding destination: "+err.Error())
					return
				}

				// If destination is nil, try to get it by address or locking script
				if destination == nil {
					destination, err = app.bux.GetDestinationByAddress(context.Background(), args[2], args[1])
					if err != nil && errors.Is(err, bux.ErrMissingDestination) {
						destination, err = app.bux.GetDestinationByLockingScript(context.Background(), args[2], args[1])
						if err != nil {
							chalker.Log(chalker.ERROR, "Error finding destination: "+err.Error())
							return
						}
					}
				}

				// Display the destination
				displayModel(destination)
			} else {
				chalker.Log(chalker.ERROR, "Unknown subcommand")
			}
		},
	}

	// Set the metadata flag
	newCmd.Flags().StringVarP(&metadata, "metadata", "m", "", "Model Metadata")

	return
}
