package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	envFileName     = ".env"
	recursionLevels = 10
)

func TestGetEnv(t *testing.T) {
	tmp := t.TempDir()
	require.NoError(t, os.WriteFile(tmp+"/"+envFileName, []byte("TEST=true"), 0o644))

	maxTmpPath := tmp
	for i := 1; i <= recursionLevels; i++ {
		maxTmpPath += fmt.Sprintf("/%d", i)
	}
	require.NoError(t, os.MkdirAll(maxTmpPath, 0o755))

	t.Run("not found", func(t *testing.T) {
		env, err := getEnv(envFileName, tmp+"/1", 0)
		require.Error(t, err)
		assert.Empty(t, env)
		assert.ErrorIs(t, ErrEnvNotFound, err)
	})

	t.Run("found without recursion", func(t *testing.T) {
		env, err := getEnv(envFileName, tmp, 0)
		require.NoError(t, err)
		assert.NotEmpty(t, env)
		assert.Equal(t, "true", env["TEST"])
	})

	t.Run("found with recursion", func(t *testing.T) {
		path := tmp

		currDir, err := os.Getwd()
		require.NoError(t, err)

		defer func() {
			endDir, err := os.Getwd()
			require.NoError(t, err)
			require.Equal(t, currDir, endDir)
		}()

		for level := 1; level <= recursionLevels; level++ {
			path += fmt.Sprintf("/%d", level)
			env, err := getEnv(envFileName, path, level)
			require.NoError(t, err)
			assert.NotEmpty(t, env)
			assert.Equal(t, "true", env["TEST"])
		}
	})
}
