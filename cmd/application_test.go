package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppConfig_GetUserAgent(t *testing.T) {
	t.Parallel()

	t.Run("get valid user agent", func(t *testing.T) {
		app := new(App)
		require.NotNil(t, app)
		agent := app.GetUserAgent()
		assert.Equal(t, "BUX-CLI: "+Version, agent)
	})
}
