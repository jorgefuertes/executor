package commands

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhich(t *testing.T) {
	const (
		cmd         = "ls"
		fullPathCmd = "/bin/" + cmd
		sym         = "/tmp/" + cmd
	)

	t.Run("known executable", func(t *testing.T) {
		ok := isExecutable(cmd)
		assert.True(t, ok, "executable '%s' not found", cmd)
		t.Run("full path", func(t *testing.T) {
			ok := isExecutable(fullPathCmd)
			assert.True(t, ok, "executable '%s' not found", fullPathCmd)
		})
	})

	t.Run("symbolic link to executable", func(t *testing.T) {
		err := os.Symlink(fullPathCmd, sym)
		require.NoError(t, err, "cannot create symbolic link '%s-->%s'", sym, fullPathCmd)
		defer os.Remove(sym)

		ok := isExecutable(sym)
		require.True(t, ok, "executable '%s' not found", sym)
	})
}
