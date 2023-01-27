package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

// Core application loader (runs before every cmd)
func commandPreprocessor() (app *App) {

	// Create a new application
	app = new(App)

	// Set up the application resources
	setupAppResources(app)

	// Load the configuration
	cobra.OnInitialize(initConfig)

	// Add config option
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Custom config file (default is $HOME/"+applicationName+"/"+configFileDefault+".yaml)")

	// Add document generation for all commands
	rootCmd.PersistentFlags().BoolVar(&generateDocs, "docs", false, "Generate docs from all commands (./"+docsLocation+")")

	// Add a toggle for request tracing
	rootCmd.PersistentFlags().BoolVarP(&skipTracing, "skip-tracing", "t", false, "Turn off request tracing information")

	// Add a toggle for disabling request caching
	rootCmd.PersistentFlags().BoolVar(&disableCache, "no-cache", false, "Turn off caching for this specific command")

	// Add a toggle for flushing all the local database cache
	rootCmd.PersistentFlags().BoolVar(&flushCache, "flush-cache", false, "Flushes ALL cache, empties local database")

	return
}

// er is a basic helper method to catch errors loading the application
func er(err error) {
	if err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Error: %s...", err.Error()))
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set
func initConfig() {

	// Custom configuration file and location
	if configFile != "" {

		chalker.Log(chalker.INFO, fmt.Sprintf("Loading custom configuration file: %s...", configFile))

		// Use config file from the flag
		viper.SetConfigFile(configFile)
	} else {

		// Make a dummy file if it doesn't exist
		file, err := os.OpenFile(filepath.Join(applicationDirectory, configFileDefault+".yaml"), os.O_RDONLY|os.O_CREATE, 0600)
		er(err)
		_ = file.Close() // Error is not needed here, just close and continue

		// Search config in home directory with name "." (without extension)
		viper.AddConfigPath(applicationDirectory)
		viper.SetConfigName(configFileDefault)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Error reading config file: %s", err.Error()))
	}

	// chalker.Log(chalker.INFO, fmt.Sprintf("...loaded config file: %s", viper.ConfigFileUsed()))
}

// generateDocumentation will generate all documentation about each command
func generateDocumentation() {

	// Replace the colorful logs in terminal (displays in Cobra docs) (color numbers generated)
	replacer := strings.NewReplacer("[32m", "```", "[33m", "```\n", "[39m", "", "[22m", "", "[36m", "", "[1m", "", "[40m", "", "[49m", "", "\u001B", "", "[0m", "")
	rootCmd.Long = replacer.Replace(rootCmd.Long)

	// Loop all command, adjust the Long description, re-add command
	for _, command := range rootCmd.Commands() {
		rootCmd.RemoveCommand(command)
		command.Long = replacer.Replace(command.Long)
		rootCmd.AddCommand(command)
	}

	// Generate the Markdown docs
	if err := doc.GenMarkdownTree(rootCmd, docsLocation); err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Error generating docs: %s", err.Error()))
		return
	}

	// Success
	chalker.Log(chalker.SUCCESS, fmt.Sprintf("Successfully generated documentation for %d commands", len(rootCmd.Commands())))
}

// setupAppResources will set up the local application directories
func setupAppResources(app *App) {

	// Find home directory
	home, err := homedir.Dir()
	er(err)

	// Set the path
	app.applicationDirectory = filepath.Join(home, applicationName)

	// Detect if we have a program folder (windows)
	_, err = os.Stat(app.applicationDirectory)
	if err != nil {
		// If it does not exist, make one!
		if os.IsNotExist(err) {
			er(os.MkdirAll(app.applicationDirectory, os.ModePerm))
		}
	}

	// Set the global application directory (used for config)
	applicationDirectory = app.applicationDirectory
}
