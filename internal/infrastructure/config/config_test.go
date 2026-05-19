package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yadro-impulse/internal/infrastructure/config"
)

func writeConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(body), 0o600))
	return path
}

func TestLoad_HappyPath(t *testing.T) {
	t.Parallel()

	path := writeConfig(t, `{"Floors":2,"Monsters":2,"OpenAt":"14:05:00","Duration":2}`)

	cfg, err := config.Load(path)
	require.NoError(t, err)
	assert.Equal(t, 2, cfg.Floors)
	assert.Equal(t, 2, cfg.Monsters)
	assert.Equal(t, 2, cfg.Duration)
}

func TestLoad_BadOpenAt(t *testing.T) {
	t.Parallel()

	path := writeConfig(t, `{"Floors":2,"Monsters":2,"OpenAt":"bad","Duration":2}`)
	_, err := config.Load(path)
	require.Error(t, err)
}

func TestLoad_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := config.Load("/nonexistent/path/config.json")
	require.Error(t, err)
}
