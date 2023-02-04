package cmd

import (
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/database"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
)

// Version of the application
var Version = "v0.1.0"

// Default flag values for various commands
var (
	applicationDirectory string // Folder path for the application resources
	configFile           string // cmd: root
	disableCache         bool   // cmd: root
	flushCache           bool   // cmd: root
	generateDocs         bool   // cmd: root
	verbose              bool   // cmd: root
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

type (

	// App is the main application struct
	// This is used to pass around the application configuration and services
	App struct {
		applicationDirectory string              // Folder path for the application resources
		bux                  bux.ClientInterface // BUX Client
		config               *Config             // Application configuration
		database             *database.DB        // CLI Application database (internal buxcli DB)
	}

	// Config is the configuration for the application and BUX
	Config struct {
		Cachestore  *CachestoreConfig        `json:"cachestore" mapstructure:"cachestore"`     // Cachestore config
		Datastore   *DatastoreConfig         `json:"datastore" mapstructure:"datastore"`       // Datastore config
		Debug       bool                     `json:"debug" mapstructure:"debug"`               // Debug mode
		Mode        string                   `json:"mode" mapstructure:"mode"`                 // Mode is either database or server
		Mongo       *datastore.MongoDBConfig `json:"mongodb" mapstructure:"mongodb"`           // MongoDB config
		Redis       *RedisConfig             `json:"redis" mapstructure:"redis"`               // Redis config
		SQL         *datastore.SQLConfig     `json:"sql" mapstructure:"sql"`                   // SQL config (MySQL, Postgres, etc)
		SQLite      *datastore.SQLiteConfig  `json:"sqlite" mapstructure:"sqlite"`             // SQLite config
		TaskManager *TaskManagerConfig       `json:"task_manager" mapstructure:"task_manager"` // TaskManager config
		Verbose     bool                     `json:"verbose" mapstructure:"verbose"`           // Verbose mode (also enables debug)
	}

	// CachestoreConfig is the configuration for the cachestore
	CachestoreConfig struct {
		Engine cachestore.Engine `json:"engine" mapstructure:"engine"` // Cache engine to use (redis, freecache)
	}

	// DatastoreConfig is a configuration for the datastore
	DatastoreConfig struct {
		AutoMigrate bool             `json:"auto_migrate" mapstructure:"auto_migrate"` // loads a blank database
		Debug       bool             `json:"debug" mapstructure:"debug"`               // true for sql statements
		Engine      datastore.Engine `json:"engine" mapstructure:"engine"`             // mysql, sqlite
		TablePrefix string           `json:"table_prefix" mapstructure:"table_prefix"` // pre_users (pre)
	}

	// RedisConfig is a configuration for Redis cachestore or taskmanager
	RedisConfig struct {
		DependencyMode        bool          `json:"dependency_mode" mapstructure:"dependency_mode"`                 // Only in Redis with script enabled
		MaxActiveConnections  int           `json:"max_active_connections" mapstructure:"max_active_connections"`   // Max active connections
		MaxConnectionLifetime time.Duration `json:"max_connection_lifetime" mapstructure:"max_connection_lifetime"` // Max connection lifetime
		MaxIdleConnections    int           `json:"max_idle_connections" mapstructure:"max_idle_connections"`       // Max idle connections
		MaxIdleTimeout        time.Duration `json:"max_idle_timeout" mapstructure:"max_idle_timeout"`               // Max idle timeout
		URL                   string        `json:"url" mapstructure:"url"`                                         // Redis URL connection string
		UseTLS                bool          `json:"use_tls" mapstructure:"use_tls"`                                 // Flag for using TLS
	}

	// TaskManagerConfig is a configuration for the taskmanager
	TaskManagerConfig struct {
		Engine    taskmanager.Engine  `json:"engine" mapstructure:"engine"`         // taskq, machinery
		Factory   taskmanager.Factory `json:"factory" mapstructure:"factory"`       // Factory (memory, redis)
		QueueName string              `json:"queue_name" mapstructure:"queue_name"` // test_queue
	}
)
