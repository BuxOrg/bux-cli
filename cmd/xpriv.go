package cmd

import (
	"encoding/hex"

	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/bitcoinschema/go-bitcoin/v2"
	"github.com/fatih/color"
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bk/wif"
	"github.com/spf13/cobra"
)

// commands for xpriv
const xprivCommandCreate = "create"
const xprivCommandName = "xpriv"
const xprivCommandWIF = "wif"
const xprivCommandXpub = "xpub"

// returnXprivCmd returns the xpriv command
func returnXprivCmd() *cobra.Command {
	return &cobra.Command{
		Use:   xprivCommandName,
		Short: "create or derive xpriv keys",
		Long: color.GreenString(`
____  _______________________._______   ____
\   \/  /\______   \______   \   \   \ /   /
 \     /  |     ___/|       _/   |\   Y   / 
 /     \  |    |    |    |   \   | \     /  
/___/\  \ |____|    |____|_  /___|  \___/   
      \_/                  \/`) + `
` + color.YellowString(`
This command is for xpriv key related commands.

create: creates a new xpriv key (`+xprivCommandName+` `+xprivCommandCreate+`)

wif: gets the WIF from the xpriv key (`+xprivCommandName+` `+xprivCommandWIF+` <xpriv>)

xpub: gets the xpub from the xpriv key (`+xprivCommandName+` `+xprivCommandXpub+` <xpriv>)
`),
		// Aliases: []string{""},
		Example: applicationName + " " + xprivCommandName + " create",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return chalker.Error(xprivCommandName + " requires a subcommand, IE: create, wif, etc.")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			// Switch on the subcommand
			if args[0] == xprivCommandCreate { // Create a new xpriv key

				// Create a new xpriv key
				key, err := bitcoin.GenerateHDKey(bitcoin.SecureSeedLength)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error: "+err.Error())
					return
				}

				chalker.Log(chalker.INFO, "xpriv: "+key.String())
			} else if args[0] == xprivCommandWIF { // Get the WIF from the private key

				// Check if xpriv is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpriv is required")
					return
				}

				// Get the hd key from the xpriv
				key, err := bitcoin.GenerateHDKeyFromString(args[1])
				if err != nil {
					chalker.Log(chalker.ERROR, "Error generating: "+err.Error())
					return
				}

				// Get the private key from the hd key
				var privateKey *bec.PrivateKey
				privateKey, err = bitcoin.GetPrivateKeyFromHDKey(key)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error generating private key: "+err.Error())
					return
				}

				// Get the WIF from the private key
				var wifKey *wif.WIF
				wifKey, err = bitcoin.PrivateKeyToWif(hex.EncodeToString(privateKey.Serialise()))
				if err != nil {
					chalker.Log(chalker.ERROR, "Error: "+err.Error())
					return
				}

				chalker.Log(chalker.INFO, "xpriv: "+args[1])
				chalker.Log(chalker.INFO, "wif: "+wifKey.String())
			} else if args[0] == xprivCommandXpub {

				// Check if xpriv is provided
				if len(args) < 2 {
					chalker.Log(chalker.ERROR, "Error: xpriv is required")
					return
				}

				// Get the hd key from the xpriv
				key, err := bitcoin.GenerateHDKeyFromString(args[1])
				if err != nil {
					chalker.Log(chalker.ERROR, "Error generating: "+err.Error())
					return
				}

				// Get the public key from the hd key
				var publicKey string
				publicKey, err = bitcoin.GetExtendedPublicKey(key)
				if err != nil {
					chalker.Log(chalker.ERROR, "Error generating public key: "+err.Error())
					return
				}

				chalker.Log(chalker.INFO, "xpriv: "+args[1])
				chalker.Log(chalker.INFO, "xpub: "+publicKey)
			} else {
				chalker.Log(chalker.ERROR, "Unknown subcommand")
			}
		},
	}
}
