package cmd

import (
	"errors"

	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/BuxOrg/bux-cli/database"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any sub-commands
var rootCmd = &cobra.Command{
	DisableAutoGenTag: true,
	Use:               applicationName,
	Short:             "Command line app for interacting with a BUX database or server",
	Example:           applicationName + " -h",
	Long: color.GreenString(`
__________ ____ _______  ___         _________ .____    .___ 
\______   \    |   \   \/  /         \_   ___ \|    |   |   |
 |    |  _/    |   /\     /   ______ /    \  \/|    |   |   |
 |    |   \    |  / /     \  /_____/ \     \___|    |___|   |
 |______  /______/ /___/\  \          \______  /_______ \___|
        \/               \_/                 \/        \/  `+Version) + `
` + color.YellowString("Author: MrZ © 2023 github.com/BuxOrg/"+applicationFullName) + `

This CLI app is used for interacting with BUX databases or servers.

Learn more about BUX: https://GetBux.io
`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// Preprocess the command line arguments and flags before executing the root command
	app := commandPreprocessor()
	var err error

	// Create a database connection (Don't require DB for now)
	if app.database, err = database.Connect(applicationName, "db_"+applicationName); err != nil {
		displayError(errors.New("Error connecting to database: " + err.Error()))
	} else {
		// Defer the database disconnection
		defer func(app *App) {
			dbErr := app.database.GarbageCollection()
			if dbErr != nil {
				displayError(errors.New("error in database GarbageCollection: " + dbErr.Error()))
			}

			if dbErr = app.database.Disconnect(); dbErr != nil {
				displayError(errors.New("error in database disconnect: " + dbErr.Error()))
			}
		}(app)
	}

	// Run root command
	er(rootCmd.Execute())

	// Generate documentation from all commands
	if generateDocs {
		generateDocumentation()
	}

	// Flush cache if requested and database is connected
	if flushCache && app.database.Connected {
		if dbErr := app.database.Flush(); dbErr != nil {
			displayError(errors.New("error in database Flush: " + dbErr.Error()))
		} else {
			chalker.Log(chalker.SUCCESS, "Successfully flushed the local database cache")
		}
	}
}
