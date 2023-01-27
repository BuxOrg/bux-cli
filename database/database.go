/*
Package database is for a local storage layer for the application
*/
package database

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v3"
)

// DB is the database struct
type DB struct {
	database  *badger.DB
	Connected bool
}

// ErrDatabaseNotConnected is returned when the database is not connected
var ErrDatabaseNotConnected = errors.New("database is not connected")

// Connect will make a new database connection and new folder/file(s) if needed
func Connect(folder, database string) (db *DB, err error) {

	// Get the home dir
	var home string
	if home, err = os.UserHomeDir(); err != nil {
		return nil, err
	}

	// Set the database file and connect (disable logging for now)
	opts := badger.DefaultOptions(filepath.Join(home, folder, database)).WithLogger(nil)
	// opts.EventLogging = false

	// Create the new database connection
	db = new(DB)
	db.database, err = badger.Open(opts)
	if err != nil {
		return nil, err
	}
	db.Connected = true
	return
}

// Disconnect will close the db connection
func (db *DB) Disconnect() error {
	err := db.database.Close()
	if err != nil {
		return err
	}
	db.Connected = false
	return nil
}

// isConnected will check if the database is connected
func isConnected(db *DB) bool {
	if db == nil || !db.Connected {
		return false
	}
	return true
}

// Set will store a new key/value pair (expiration optional)
func (db *DB) Set(key, value string, ttl time.Duration) error {
	if !isConnected(db) {
		return ErrDatabaseNotConnected
	}
	return db.database.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value))
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}

		return txn.SetEntry(entry)
	})
}

// Get will retrieve a value from a key (if found)
func (db *DB) Get(key string) (string, error) {
	if !isConnected(db) {
		return "", ErrDatabaseNotConnected
	}
	var valCopy []byte
	err := db.database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		valCopy, err = item.ValueCopy(nil)
		return err
	})

	// Not found (don't return an error, as we want to use this as cache)
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = nil
	}

	return string(valCopy), err
}

// Flush will empty the entire database
func (db *DB) Flush() error {
	if !isConnected(db) {
		return ErrDatabaseNotConnected
	}
	return db.database.DropAll()
}

// GarbageCollection will clean up some garbage in the database (reduces space, etc.)
func (db *DB) GarbageCollection() (err error) {
	if !isConnected(db) {
		return ErrDatabaseNotConnected
	}
	err = db.database.RunValueLogGC(0.5)
	if errors.Is(err, badger.ErrNoRewrite) {
		return nil
	}
	return
}
