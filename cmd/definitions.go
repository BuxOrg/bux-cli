package cmd

import (
	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/database"
	"github.com/mrz1836/go-datastore"
)

// Version is set manually (also make:build overwrites this value from GitHub's latest tag)
var Version = "v0.1.0"

// Default flag values for various commands
var (
	applicationDirectory string // Folder path for the application resources
	configFile           string // cmd: root
	disableCache         bool   // cmd: root
	flushCache           bool   // cmd: root
	generateDocs         bool   // cmd: root
)

// Defaults for the application
const (
	applicationFullName = "bux-cli"       // Full name of the application (long version)
	applicationName     = "buxcli"        // Application name (binary) (short version
	configFileDefault   = "config"        // Config file name
	docsLocation        = "docs/commands" // Default location for command documentation
	modeDatabase        = "database"      // Mode for database
	modeServer          = "server"        // Mode for server
)

// App is the main application struct
type App struct {
	applicationDirectory string              // Folder path for the application resources
	bux                  bux.ClientInterface // BUX Client
	database             *database.DB        // Database connection (internal buxcli DB)
	config               struct {
		Mode     string                  `json:"mode"`     // Mode is either database or server
		Database *datastore.SQLiteConfig `json:"database"` // Database Config - SQLite
	}
}
