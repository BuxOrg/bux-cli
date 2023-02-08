package cmd

import (
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-cli/database"
	"github.com/BuxOrg/bux/taskmanager"
	"github.com/mrz1836/go-cachestore"
	"github.com/mrz1836/go-datastore"
	"github.com/mrz1836/go-whatsonchain"
)

// Version of the application
var Version = "v0.1.1"

// Default flag values for various commands
var (
	applicationDirectory string // Folder path for the application resources
	configFile           string // cmd: root
	disableCache         bool   // cmd: root
	draftID              string // cmd: tx
	flushCache           bool   // cmd: root
	generateDocs         bool   // cmd: root
	metadata             string // cmd: tx, xpub, destination
	txConfig             string // cmd: tx
	txHex                string // cmd: tx
	txID                 string // cmd: tx
	verbose              bool   // cmd: root
	wocEnabled           bool   // cmd: tx
	xpubID               string // cmd: destination
)

// Flags for the application
const (
	flagMetadata       = "metadata"
	flagMetadataShort  = "m"
	flagTxConfig       = "txconfig"
	flagTxConfigShort  = "c"
	flagTxDraftID      = "draft"
	flagTxDraftIDShort = "d"
	flagTxHex          = "hex"
	flagTxHexShort     = "x"
	flagTxID           = "txid"
	flagTxIDShort      = "i"
	flagWoc            = "woc"
	flagWocShort       = "w"
	flagXpubID         = "xpubid"
	flagXpubIDShort    = "x"
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
		Chainstate  *ChainstateConfig        `json:"chainstate" mapstructure:"chainstate"`     // Chainstate config
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

	// ChainstateConfig is a configuration for the chainstate
	ChainstateConfig struct {
		Broadcasting       bool   `json:"broadcast" mapstructure:"broadcast"`                     // true for broadcasting
		BroadcastInstantly bool   `json:"broadcast_instantly" mapstructure:"broadcast_instantly"` // true for broadcasting instantly
		P2P                bool   `json:"p2p" mapstructure:"p2p"`                                 // true for p2p
		SyncOnChain        bool   `json:"sync_on_chain" mapstructure:"sync_on_chain"`             // true for syncing on chain
		TaalAPIKey         string `json:"taal_api_key" mapstructure:"taal_api_key"`               // Taal API key
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

	// XpubExtended is an extended xpub struct with the full key
	XpubExtended struct {
		*bux.Xpub
		FullKey string `json:"full_key" mapstructure:"full_key"`
	}

	// Keys is a struct for the private keys, wif, xpriv and xpub
	Keys struct {
		PrivateKey string `json:"private_key" mapstructure:"private_key"`
		WIF        string `json:"wif" mapstructure:"wif"`
		Xpriv      string `json:"xpriv" mapstructure:"xpriv"`
		Xpub       string `json:"xpub" mapstructure:"xpub"`
	}

	// Transaction is a struct for the bux model and whatsonchain transaction
	Transaction struct {
		Bux *bux.Transaction     `json:"bux" mapstructure:"bux"`
		WOC *whatsonchain.TxInfo `json:"woc,omitempty" mapstructure:"woc"`
	}

	// Destination is a struct for the bux model and whatsonchain destination
	Destination struct {
		Bux        *bux.Destination             `json:"bux" mapstructure:"bux"`
		WOCBalance *whatsonchain.AddressBalance `json:"woc_balance,omitempty" mapstructure:"woc_balance"`
		WOCInfo    *whatsonchain.AddressInfo    `json:"woc_info,omitempty" mapstructure:"woc_info"`
	}
)
