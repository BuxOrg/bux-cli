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
		Short: "manage and interact with destinations in BUX",
		Long: color.GreenString(`
________  ___________ ____________________.___ _______      ________________.___________    _______   
\______ \ \_   _____//   _____/\__    ___/|   |\      \    /  _  \__    ___/|   \_____  \   \      \  
 |    |  \ |    __)_ \_____  \   |    |   |   |/   |   \  /  /_\  \|    |   |   |/   |   \  /   |   \ 
 |    '   \|        \/        \  |    |   |   /    |    \/    |    \    |   |   /    |    \/    |    \
 /_______  /_______  /_______  /  |____|   |___\____|__  /\____|__  /____|   |___\_______  /\____|__  /
		 \/        \/        \/                        \/         \/                     \/         \/`) + `
` + color.YellowString(`
This command is for destination (address, locking script) related commands.

new: creates a new destination in BUX (`+destinationCommandName+` new <xpub>)
get: gets an existing destination in BUX (`+destinationCommandName+` get <destination_id | address | locking_script> -x=<xpub_id>)
`),
		Aliases: []string{"address"},
		Example: applicationName + " " + destinationCommandName + " " + destinationCommandNew + " <xpub>",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(destinationCommandName + " requires a subcommand, IE: " + destinationCommandNew + ", etc.")
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

			// Get the optional woc flag from the flags
			if wocEnabled, err = cmd.Flags().GetBool(flagWoc); err != nil {
				displayError(errors.New("error getting woc flag: " + err.Error()))
				return
			}

			// Not a valid subcommand
			if args[0] != destinationCommandNew && args[0] != destinationCommandGet {
				displayError(ErrUnknownSubcommand)
				return
			}

			// Switch on the subcommand
			if args[0] == destinationCommandNew { // Create a new destination

				// Check if xpub is provided
				if len(args) < 2 {
					displayError(ErrXpubIsRequired)
					return
				}

				// Create the destination
				var destination *Destination
				destination, err = newDestination(context.Background(), app, args[1], metadata)
				if err != nil {
					displayError(errors.New("error creating destination: " + err.Error()))
					return
				}

				// Display the destination
				displayModel(destination)

			} else if args[0] == destinationCommandGet { // Get a destination

				// Check if xpub id is provided
				if len(xpubID) <= 0 {
					displayError(ErrXpubIDIsRequired)
					return
				}

				// Get the destination
				var destination *Destination
				destination, err = getDestination(context.Background(), app, args[1], xpubID, wocEnabled)
				if err != nil {
					displayError(errors.New("error getting destination: " + err.Error()))
					return
				}

				// Display the destination
				displayModel(destination)
			}
		},
	}

	// Set the metadata flag
	newCmd.Flags().StringVarP(&metadata, flagMetadata, flagMetadataShort, "", "Model Metadata")

	// Set the xpub id flag
	newCmd.Flags().StringVarP(&xpubID, flagXpubID, flagXpubIDShort, "", "Xpub ID")

	// Set the woc flag
	newCmd.Flags().BoolP(
		flagWoc, flagWocShort, wocEnabled,
		"Optional flag to use WhatsOnChain for additional address data",
	)

	return
}

// newDestination creates a new destination
// app: the app
// xpubKey: the xpub key
func newDestination(ctx context.Context, app *App,
	xpubKey, metadata string) (destination *Destination, err error) {

	destination = new(Destination)

	var xpub *bux.Xpub
	xpub, err = app.bux.GetXpub(ctx, xpubKey)
	if err != nil {
		return
	}
	if xpub == nil {
		err = errors.New("xpub not found")
		return
	}

	// Get the metadata if provided
	modelOps := app.bux.DefaultModelOptions()
	if len(metadata) > 0 {
		modelOps = append(modelOps, bux.WithMetadataFromJSON([]byte(metadata)))
	}

	// Create the destination
	destination.Bux, err = app.bux.NewDestination(
		ctx, xpubKey, utils.ChainExternal, utils.ScriptTypePubKeyHash, false, modelOps...,
	)

	return
}

// getDestination gets a destination by ID, address or locking script
// app: the app
// idOrAddressOrScript: the destination ID, address or locking script
// xpubID: the xpub ID
func getDestination(ctx context.Context, app *App, idOrAddressOrScript, xpubID string,
	wocEnabled bool) (destination *Destination, err error) {

	destination = new(Destination)

	// Get the destination by ID, address or locking script
	destination.Bux, err = app.bux.GetDestinationByID(ctx, xpubID, idOrAddressOrScript)
	if err != nil && !errors.Is(err, bux.ErrMissingDestination) {
		return
	}

	// If destination is nil, try to get it by address or locking script
	if destination.Bux == nil {
		destination.Bux, err = app.bux.GetDestinationByAddress(ctx, xpubID, idOrAddressOrScript)
		if err != nil && errors.Is(err, bux.ErrMissingDestination) {
			destination.Bux, err = app.bux.GetDestinationByLockingScript(ctx, xpubID, idOrAddressOrScript)
			if err != nil {
				err = errors.New("error finding destination: " + err.Error())
			}
		}
	}

	// If destination is not nil and WhatsOnChain is enabled, get the address data
	if destination.Bux != nil && len(destination.Bux.Address) > 0 && wocEnabled {

		// Get the address info from WhatsOnChain
		destination.WOCInfo, err = app.bux.Chainstate().WhatsOnChain().AddressInfo(ctx, destination.Bux.Address)
		if err != nil {
			return
		}

		// Get the balance from WhatsOnChain
		destination.WOCBalance, err = app.bux.Chainstate().WhatsOnChain().AddressBalance(ctx, destination.Bux.Address)
		if err != nil {
			return
		}
	}

	return
}
