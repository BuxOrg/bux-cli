package database

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var applicationName = "buxcli"

func TestConnect(t *testing.T) {
	t.Run("TestConnect Success", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")
		err = db.Disconnect()
		require.NoError(t, err)
	})
}

func TestDB_Disconnect(t *testing.T) {
	t.Run("TestDB_Disconnect Success", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")
		err = db.Disconnect()
		require.NoError(t, err)
		require.Equalf(t, false, db.Connected, "Database is still connected")
	})
}

func TestDB_Set(t *testing.T) {
	t.Run("TestDB_Set Success", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")
		err = db.Set("key", "value", time.Minute)
		require.NoError(t, err)

		var val string
		val, err = db.Get("key")
		require.NoError(t, err)
		require.Equalf(t, "value", val, "Value is not equal to 'value'")

		err = db.Disconnect()
		require.NoError(t, err)
	})

	t.Run("TestDB_Set Error", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")

		err = db.Disconnect()
		require.NoError(t, err)

		err = db.Set("key", "value", time.Minute)
		require.Error(t, err)
		require.Equalf(t, ErrDatabaseNotConnected, err, "Error is not equal to ErrDatabaseNotConnected")
	})
}

func TestDB_Get(t *testing.T) {
	t.Run("TestDB_Get Success", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")
		err = db.Set("key", "value", time.Minute)
		require.NoError(t, err)

		var val string
		val, err = db.Get("key")
		require.NoError(t, err)
		require.Equalf(t, "value", val, "Value is not equal to 'value'")

		err = db.Disconnect()
		require.NoError(t, err)
	})

	t.Run("TestDB_Get Error", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")

		err = db.Disconnect()
		require.NoError(t, err)

		var val string
		val, err = db.Get("key")
		require.Error(t, err)
		require.Equalf(t, "", val, "Value is not empty")
		require.Equalf(t, ErrDatabaseNotConnected, err, "Error is not equal to ErrDatabaseNotConnected")
	})
}

func TestDB_Flush(t *testing.T) {
	t.Run("TestDB_Flush Success", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")
		err = db.Set("key", "value", time.Minute)
		require.NoError(t, err)

		var val string
		val, err = db.Get("key")
		require.NoError(t, err)
		require.Equalf(t, "value", val, "Value is not equal to 'value'")

		err = db.Flush()
		require.NoError(t, err)

		val, err = db.Get("key")
		require.NoError(t, err)
		require.Equalf(t, "", val, "Value is not empty")

		err = db.Disconnect()
		require.NoError(t, err)
	})

	t.Run("TestDB_Flush Error", func(t *testing.T) {
		db, err := Connect(applicationName, "db_"+applicationName)
		require.NoError(t, err)
		require.NotNil(t, db)
		require.Equalf(t, true, db.Connected, "Database is not connected")

		err = db.Disconnect()
		require.NoError(t, err)

		err = db.Flush()
		require.Error(t, err)
		require.Equalf(t, ErrDatabaseNotConnected, err, "Error is not equal to ErrDatabaseNotConnected")
	})
}
