package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/chalker"
	"github.com/BuxOrg/bux/chainstate"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/go-homedir"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"github.com/tonicpow/go-minercraft"
)

// Added a mutex lock for a race-condition
var viperLock sync.Mutex

// commandPreprocessor is the core application loader (runs before every cmd)
func commandPreprocessor() (app *App) {

	// Create a new application
	app = new(App)

	// Set up the application resources
	setupAppResources(app)

	// Load the configuration
	initConfig(app)

	// Mode is required
	if viper.GetString("mode") == "" {
		er(ErrModeIsRequired)
	}

	// Add config option
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "custom config file (default is $HOME/"+applicationName+"/"+configFileDefault+".json)")

	// Add document generation for all commands
	rootCmd.PersistentFlags().BoolVar(&generateDocs, "docs", false, "generate docs from all commands (./"+docsLocation+")")

	// Add a toggle for disabling request caching
	rootCmd.PersistentFlags().BoolVar(&disableCache, "no-cache", false, "turn off caching for this specific command")

	// Add a toggle for flushing all the local database cache
	rootCmd.PersistentFlags().BoolVar(&flushCache, "flush-cache", false, "flushes ALL cache, empties local temporary database")

	// Add a toggle for verbose logging
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose logging")

	// Add xpriv command
	rootCmd.AddCommand(returnXprivCmd())

	// Add xpub command
	rootCmd.AddCommand(returnXpubCmd(app))

	// Add destination command
	rootCmd.AddCommand(returnDestinationCmd(app))

	// Add transaction command
	rootCmd.AddCommand(returnTransactionCmd(app))

	return
}

// displayModel will display a model in a pretty format
func displayModel(v any) {
	if v == nil {
		displayError(ErrModelIsNil)
		return
	}
	if b, err := json.MarshalIndent(v, "", "  "); err != nil {
		displayError(fmt.Errorf("error marshaling model: %w", err))
	} else {
		chalker.Log(chalker.INFO, string(b))
	}
}

// displayError will display an error in a pretty format
func displayError(err error) {
	if err != nil {
		chalker.Log(chalker.ERROR, fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
}

// er is a basic helper method to catch errors loading the application
func er(err error) {
	if err != nil {
		displayError(err)
		os.Exit(1)
	}
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
		displayError(fmt.Errorf("error generating docs: %w", err))
		return
	}

	// Success
	chalker.Log(chalker.SUCCESS, fmt.Sprintf("Successfully generated documentation for %d commands", len(rootCmd.Commands())))
}

// initConfig reads in config file and ENV variables if set
func initConfig(app *App) {

	// Create a lock for modifying the config
	viperLock.Lock()

	// Unlock the lock when we're done
	defer func() {
		viperLock.Unlock()
	}()

	// Custom configuration file and location
	if configFile != "" {

		verboseLog(func() {
			chalker.Log(chalker.INFO, fmt.Sprintf("Loading custom configuration file: %s...", configFile))
		})

		// Use config file from the flag
		viper.SetConfigFile(configFile)
	} else {

		// Check if the default config file exists
		if _, err := os.Stat(filepath.Join(applicationDirectory, configFileDefault+".json")); err != nil {

			// Read the example config
			var content []byte
			content, err = os.ReadFile("config-example.json")
			er(err)

			// Make a dummy file if it doesn't exist
			var file *os.File
			file, err = os.OpenFile(filepath.Join(applicationDirectory, configFileDefault+".json"), os.O_RDWR|os.O_CREATE, 0755) //nolint:gosec // We don't care about the permissions
			er(err)

			defer func() {
				_ = file.Close()
			}()

			// Write the example config into the file
			_, err = file.WriteString(string(content))
			er(err)
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
	er(viper.ReadInConfig())

	// Unmarshal into values struct
	er(viper.Unmarshal(&app.config))

	// Fix for relative paths in database configuration (SQLite)
	usr, _ := user.Current()
	dir := usr.HomeDir
	if app.config.SQLite != nil && len(app.config.SQLite.DatabasePath) > 0 {
		if app.config.SQLite.DatabasePath == "~" {
			// In case of "~", which won't be caught by the "else if"
			app.config.SQLite.DatabasePath = dir
		} else if strings.HasPrefix(app.config.SQLite.DatabasePath, "~/") {
			// Use strings.HasPrefix so we don't match paths like
			// "/something/~/something/"
			app.config.SQLite.DatabasePath = filepath.Join(dir, app.config.SQLite.DatabasePath[2:])
		}
	}

	// Set the verbose logging from the config file
	if app.config.Verbose {
		app.config.Debug = true
		verbose = true
	}

	verboseLog(func() {
		chalker.Log(chalker.INFO, fmt.Sprintf("...loaded config file: %s", viper.ConfigFileUsed()))
	})
}

// loadBux will load BUX into the app
func loadBux(app *App) (loaded bool) {

	// Start building BUX client options
	var options []bux.ClientOps

	// Flag for debugging
	if app.config.Debug || verbose {
		options = append(options, bux.WithDebugging())
	}

	// Customize the outgoing user agent
	options = append(options, bux.WithUserAgent(app.GetUserAgent()))

	// Switch on the mode
	if app.config.Mode == modeDatabase {

		// Load cache
		if app.config.Cachestore.Engine == cachestore.Redis {
			options = append(options, bux.WithRedis(&cachestore.RedisConfig{
				DependencyMode:        app.config.Redis.DependencyMode,
				MaxActiveConnections:  app.config.Redis.MaxActiveConnections,
				MaxConnectionLifetime: app.config.Redis.MaxConnectionLifetime,
				MaxIdleConnections:    app.config.Redis.MaxIdleConnections,
				MaxIdleTimeout:        app.config.Redis.MaxIdleTimeout,
				URL:                   app.config.Redis.URL,
				UseTLS:                app.config.Redis.UseTLS,
			}))
		} else if app.config.Cachestore.Engine == cachestore.FreeCache {
			options = append(options, bux.WithFreeCache())
		}

		// Set the datastore
		var err error
		if options, err = loadDatastore(options, app); err != nil {
			displayError(fmt.Errorf("error loading datastore: %w", err))
			return
		}

		// Load task manager (redis or taskq)
		if app.config.TaskManager.Engine == taskmanager.TaskQ {
			config := taskmanager.DefaultTaskQConfig(app.config.TaskManager.QueueName)
			if app.config.TaskManager.Factory == taskmanager.FactoryRedis {
				options = append(
					options,
					bux.WithTaskQUsingRedis(
						config,
						&redis.Options{
							Addr: strings.Replace(app.config.Redis.URL, "redis://", "", -1),
						},
					))
			} else {
				options = append(options, bux.WithTaskQ(config, app.config.TaskManager.Factory))
			}
		}

		// Add chainstate options
		options = append(options, bux.WithChainstateOptions(
			app.config.Chainstate.Broadcasting,
			app.config.Chainstate.BroadcastInstantly,
			app.config.Chainstate.P2P,
			app.config.Chainstate.SyncOnChain,
		))

		// Exclude providers (NowNodes needs API key) // todo: make this configurable
		options = append(options, bux.WithExcludedProviders([]string{chainstate.ProviderNowNodes}))

		// todo: allow custom miners for minercraft

		// Custom rates and custom miners
		if len(app.config.Chainstate.TaalAPIKey) > 0 {
			verboseLog(func() {
				chalker.Log(chalker.INFO, "taal api key detected and loaded")
			})
			miners, _ := minercraft.DefaultMiners()
			taal := minercraft.MinerByName(miners, minercraft.MinerTaal)
			taal.Token = app.config.Chainstate.TaalAPIKey
			options = append(options, bux.WithBroadcastMiners([]*chainstate.Miner{{Miner: taal}}))
			options = append(options, bux.WithQueryMiners([]*chainstate.Miner{{Miner: taal}}))
		}

		// Load BUX
		app.bux, err = bux.NewClient(context.Background(), options...)
		if err != nil {
			displayError(errors.New("Error loading BUX: " + err.Error()))
			return
		}

	} else if app.config.Mode == modeServer {
		displayError(ErrServerModeIsNotImplemented)
		return
	} else {
		displayError(ErrUnknownMode)
		return
	}

	// Success on loading?
	if app.bux != nil {
		loaded = true

		verboseLog(func() {
			chalker.Log(chalker.SUCCESS, fmt.Sprintf("Successfully loaded BUX version: %s", app.bux.UserAgent()))
		})

		// Print some basic stats
		if app.config.Verbose {
			printBuxStats(app)
		}
	}

	return
}

// loadDatastore will load the correct datastore based on the engine
func loadDatastore(options []bux.ClientOps, app *App) ([]bux.ClientOps, error) {

	// Select the datastore
	if app.config.Datastore.Engine == datastore.SQLite {
		debug := app.config.Datastore.Debug
		tablePrefix := app.config.Datastore.TablePrefix
		if len(app.config.SQLite.TablePrefix) > 0 {
			tablePrefix = app.config.SQLite.TablePrefix
		}
		options = append(options, bux.WithSQLite(&datastore.SQLiteConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:       debug,
				TablePrefix: tablePrefix,
			},
			DatabasePath: app.config.SQLite.DatabasePath, // "" for in memory
			Shared:       app.config.SQLite.Shared,
		}))
	} else if app.config.Datastore.Engine == datastore.MySQL || app.config.Datastore.Engine == datastore.PostgreSQL {
		debug := app.config.Datastore.Debug
		tablePrefix := app.config.Datastore.TablePrefix
		if len(app.config.SQL.TablePrefix) > 0 {
			tablePrefix = app.config.SQL.TablePrefix
		}

		options = append(options, bux.WithSQL(app.config.Datastore.Engine, &datastore.SQLConfig{
			CommonConfig: datastore.CommonConfig{
				Debug:                 debug,
				MaxConnectionIdleTime: app.config.SQL.MaxConnectionIdleTime,
				MaxConnectionTime:     app.config.SQL.MaxConnectionTime,
				MaxIdleConnections:    app.config.SQL.MaxIdleConnections,
				MaxOpenConnections:    app.config.SQL.MaxOpenConnections,
				TablePrefix:           tablePrefix,
			},
			Driver:    app.config.Datastore.Engine.String(),
			Host:      app.config.SQL.Host,
			Name:      app.config.SQL.Name,
			Password:  app.config.SQL.Password,
			Port:      app.config.SQL.Port,
			TimeZone:  app.config.SQL.TimeZone,
			TxTimeout: app.config.SQL.TxTimeout,
			User:      app.config.SQL.User,
		}))

	} else if app.config.Datastore.Engine == datastore.MongoDB {

		debug := app.config.Datastore.Debug
		tablePrefix := app.config.Datastore.TablePrefix
		if len(app.config.Mongo.TablePrefix) > 0 {
			tablePrefix = app.config.Mongo.TablePrefix
		}
		app.config.Mongo.Debug = debug
		app.config.Mongo.TablePrefix = tablePrefix
		options = append(options, bux.WithMongoDB(app.config.Mongo))
	} else {
		return nil, errors.New("unsupported datastore engine: " + app.config.Datastore.Engine.String())
	}

	// Add the auto migrate
	if app.config.Datastore.AutoMigrate {
		options = append(options, bux.WithAutoMigrate(bux.BaseModels...))
	}

	return options, nil
}

// printBuxStats will print some basic BUX statistics
func printBuxStats(app *App) {
	count, err := app.bux.GetXPubsCount(context.Background(), nil, nil)
	if err != nil {
		displayError(err)
	} else {
		chalker.Log(chalker.SUCCESS, fmt.Sprintf("xpubs found: %d", count))
	}

	count, err = app.bux.GetDestinationsCount(context.Background(), nil, nil)
	if err != nil {
		displayError(err)
	} else {
		chalker.Log(chalker.SUCCESS, fmt.Sprintf("destinations found: %d", count))
	}
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

// verboseLog is a helper method to log only if verbose is enabled
func verboseLog(logLine func()) {
	if verbose {
		logLine()
	}
}

// GetUserAgent will return the outgoing user agent
func (a *App) GetUserAgent() string {
	return "BUX-CLI: " + Version
}

// InitializeBUX will initialize BUX if it is not already initialized
func (a *App) InitializeBUX() (deferFunc func()) {

	// Load BUX if not already loaded
	if a.bux == nil {
		loaded := loadBux(a)

		// Fail if BUX is not loaded
		if !loaded {
			er(ErrFailedToLoadBux)
		}
	}

	// Return a function to close BUX
	deferFunc = func() {
		if a.bux != nil {
			_ = a.bux.Close(context.Background())
		}
	}
	return
}
