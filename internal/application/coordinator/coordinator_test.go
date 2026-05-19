package coordinator_test

import (
	"bytes"
	_ "embed"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yadro-impulse/internal/application/coordinator"
	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/infrastructure/eventsreader"
	"yadro-impulse/internal/infrastructure/output"
)

//go:embed testdata/events.txt
var exampleEvents string

//go:embed testdata/expected.txt
var exampleExpected string

func TestCoordinator_GoldenExample(t *testing.T) {
	t.Parallel()

	at, err := domain.ParseTime("14:05:00")
	require.NoError(t, err)
	dungeon, err := domain.NewDungeon(domain.DungeonConfig{Floors: 2, Monsters: 2, OpenAt: at, Duration: 2})
	require.NoError(t, err)

	var buf bytes.Buffer
	coord := coordinator.New(dungeon, output.NewWriter(&buf))
	reader := eventsreader.New(strings.NewReader(exampleEvents))

	for {
		ev, readErr := reader.Next()
		if errors.Is(readErr, io.EOF) {
			break
		}
		require.NoError(t, readErr)
		require.NoError(t, coord.Handle(ev))
	}
	coord.CloseDungeon()
	require.NoError(t, output.WriteReport(&buf, coord.Report()))

	assert.Equal(t, exampleExpected, buf.String())
}

func TestCoordinator_DungeonCloseFailsActivePlayer(t *testing.T) {
	t.Parallel()

	at, _ := domain.ParseTime("14:00:00")
	dungeon, _ := domain.NewDungeon(domain.DungeonConfig{Floors: 2, Monsters: 2, OpenAt: at, Duration: 1})

	var buf bytes.Buffer
	coord := coordinator.New(dungeon, output.NewWriter(&buf))

	events := []domain.Event{
		{Time: domain.Time{H: 14}, ID: domain.Register, PlayerID: 1},
		{Time: domain.Time{H: 14, M: 5}, ID: domain.Enter, PlayerID: 1},
	}
	for _, ev := range events {
		require.NoError(t, coord.Handle(ev))
	}

	coord.CloseDungeon()
	r := coord.Report()
	require.Len(t, r.Players, 1)
	assert.Equal(t, "FAIL", r.Players[0].Status)
}
