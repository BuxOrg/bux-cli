package cmd

import (
	"encoding/hex"
	"errors"

	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/fatih/color"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/wif"
	"github.com/spf13/cobra"
)

// commands for xpriv
const xprivCommandName = "xpriv"
const xprivCommandNew = "new"
const xprivCommandInfo = "info"

// returnXprivCmd returns the xpriv command
func returnXprivCmd() *cobra.Command {
	return &cobra.Command{
		Use:   xprivCommandName,
		Short: "create new xpriv keys and see additional info",
		Long: color.GreenString(`
____  _______________________._______   ____
\   \/  /\______   \______   \   \   \ /   /
 \     /  |     ___/|       _/   |\   Y   / 
 /     \  |    |    |    |   \   | \     /  
/___/\  \ |____|    |____|_  /___|  \___/   
      \_/                  \/`) + `
` + color.YellowString(`
This command is for xpriv key related commands. These commands are read-only and no data is stored on the BUX servers.

new: creates a new xpriv key (`+xprivCommandName+` `+xprivCommandNew+`)
info: gets the xpub, WIF and other info from the xpriv key (`+xprivCommandName+` `+xprivCommandInfo+` <xpriv>)
`),
		// Aliases: []string{"priv"},
		Example: applicationName + " " + xprivCommandName + " " + xprivCommandNew,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(xprivCommandName + " requires a subcommand, IE: " + xprivCommandNew + ", " + xprivCommandInfo + ", etc.")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			keys := new(Keys)

			// Switch on the subcommand
			if args[0] == xprivCommandNew { // Create a new xpriv key

				// Create a new xpriv key
				key, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
				if err != nil {
					displayError(errors.New("error generating new xpriv: " + err.Error()))
					return
				}
				keys.Xpriv = key.String()

				displayModel(keys)
			} else if args[0] == xprivCommandInfo { // Get the xpub, WIF and other info from the xpriv key

				// Check if xpriv is provided
				if len(args) < 2 {
					displayError(ErrXprivIsRequired)
					return
				}
				keys.Xpriv = args[1]

				// Get the hd key from the xpriv
				key, err := bitcoin.GenerateHDKeyFromString(keys.Xpriv)
				if err != nil {
					displayError(errors.New("error generating HD key from xpriv: " + err.Error()))
					return
				}

				// Get the public key from the hd key
				keys.Xpub, err = bitcoin.GetExtendedPublicKey(key)
				if err != nil {
					displayError(errors.New("error generating xpub from xpriv: " + err.Error()))
					return
				}

				// Get the private key from the hd key
				var privateKey *bec.PrivateKey
				privateKey, err = bitcoin.GetPrivateKeyFromHDKey(key)
				if err != nil {
					displayError(errors.New("error generating private key from xpriv: " + err.Error()))
					return
				}
				keys.PrivateKey = hex.EncodeToString(privateKey.Serialise())

				// Get the WIF from the private key
				var wifKey *wif.WIF
				wifKey, err = bitcoin.PrivateKeyToWif(keys.PrivateKey)
				if err != nil {
					displayError(errors.New("error generating WIF from xpriv: " + err.Error()))
					return
				}
				keys.WIF = wifKey.String()

				displayModel(keys)
			} else {
				displayError(ErrUnknownSubcommand)
			}
		},
	}
}
