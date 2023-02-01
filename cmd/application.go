package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// Core application loader (runs before every cmd)
func commandPreprocessor() (app *App, deferFunc func()) {

	// Create a new application
	app = new(App)

	// Set up the application resources
	setupAppResources(app)

	// Load the configuration
	initConfig(app)

	// Mode is required
	if viper.GetString("mode") == "" {
		chalker.Log(chalker.ERROR, "Mode is required")
		os.Exit(1)
	}

	// Load BUX
	loaded := loadBux(app)
	deferFunc = func() {
		if app.bux != nil {
			_ = app.bux.Close(context.Background())
		}
	}

	// Fail if BUX is not loaded
	if !loaded {
		chalker.Log(chalker.ERROR, "Error loading BUX")
		os.Exit(1)
	}

	// Add config option
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Custom config file (default is $HOME/"+applicationName+"/"+configFileDefault+".json)")

	// Add document generation for all commands
	rootCmd.PersistentFlags().BoolVar(&generateDocs, "docs", false, "Generate docs from all commands (./"+docsLocation+")")

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
func initConfig(app *App) {

	// Create a lock for modifying the config
	viperLock.Lock()

	// Custom configuration file and location
	if configFile != "" {

		chalker.Log(chalker.INFO, fmt.Sprintf("Loading custom configuration file: %s...", configFile))

		// Use config file from the flag
		viper.SetConfigFile(configFile)
	} else {

		// Check if the default config file exists
		if _, err := os.Stat(filepath.Join(applicationDirectory, configFileDefault+".json")); err != nil {

			// Read the example config
			var content []byte
			content, err = os.ReadFile("config-example.json")
			if err != nil {
				chalker.Log(chalker.ERROR, fmt.Sprintf("Error reading example config: %s", err.Error()))
			}

			// Make a dummy file if it doesn't exist
			var file *os.File
			file, err = os.OpenFile(filepath.Join(applicationDirectory, configFileDefault+".json"), os.O_RDWR|os.O_CREATE, 0755) //nolint:gosec // We don't care about the permissions
			er(err)

			defer func() {
				_ = file.Close()
			}()

			// Write the example config into the file
			if _, err = file.WriteString(string(content)); err != nil {
				chalker.Log(chalker.ERROR, fmt.Sprintf("Error writing config file: %s", err.Error()))
			}
		}

		// Search config in home directory with name "." (without extension)
		viper.AddConfigPath(applicationDirectory)
		viper.SetConfigName(configFileDefault)
	}

	// Set a replacer for replacing double underscore with nested period
	replacer := strings.NewReplacer(".", "__")
	viper.SetEnvKeyReplacer(replacer)

	// Set the prefix
	viper.SetEnvPrefix(applicationName)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Error reading config file: %s", err.Error()))
	}

	// Unmarshal into values struct
	if err := viper.Unmarshal(&app.config); err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Error unmarshalling config file: %s", err.Error()))
	}

	// Fix for relative paths in database configuration
	usr, _ := user.Current()
	dir := usr.HomeDir
	if app.config.Database != nil && len(app.config.Database.DatabasePath) > 0 {
		if app.config.Database.DatabasePath == "~" {
			// In case of "~", which won't be caught by the "else if"
			app.config.Database.DatabasePath = dir
		} else if strings.HasPrefix(app.config.Database.DatabasePath, "~/") {
			// Use strings.HasPrefix so we don't match paths like
			// "/something/~/something/"
			app.config.Database.DatabasePath = filepath.Join(dir, app.config.Database.DatabasePath[2:])
		}
	}

	// Unlock now that the configuration is complete
	viperLock.Unlock()

	chalker.Log(chalker.INFO, fmt.Sprintf("...loaded config file: %s", viper.ConfigFileUsed()))
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

// loadBux will load BUX into the app
func loadBux(app *App) (loaded bool) {

	// Check the mode
	if app.config.Mode == modeDatabase {

		// Load BUX
		var err error
		app.bux, err = bux.NewClient(
			context.Background(),                   // Set context
			bux.WithAutoMigrate(bux.BaseModels...), // Auto migrate the database
			bux.WithFreeCache(),                    // Use in-memory cache
			bux.WithSQLite(app.config.Database),    // SQL Lite connection
			bux.WithTaskQ(taskmanager.DefaultTaskQConfig("local_queue"), taskmanager.FactoryMemory), // Tasks
		)
		if err != nil {
			chalker.Log(chalker.ERROR, fmt.Sprintf("Error loading BUX: %s", err.Error()))
		}

	} else if app.config.Mode == modeServer {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Mode is not implemented: %s", app.config.Mode))
	} else {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Unknown mode: %s", app.config.Mode))
	}

	chalker.Log(chalker.SUCCESS, fmt.Sprintf("Successfully loaded BUX version: %s", app.bux.UserAgent()))

	// Print some basic stats
	printBuxStats(app)

	if app.bux != nil {
		loaded = true
	}
	return
}

// printBuxStats will print some basic BUX statistics
func printBuxStats(app *App) {
	count, err := app.bux.GetXPubsCount(context.Background(), nil, nil)
	if err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf("Error getting xpub count: %s", err.Error()))
	} else {
		chalker.Log(chalker.SUCCESS, fmt.Sprintf("Xpubs Found: %d", count))
	}
}
