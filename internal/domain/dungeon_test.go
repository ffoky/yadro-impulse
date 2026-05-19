package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yadro-impulse/internal/domain"
)

func mustParseTime(t *testing.T, s string) domain.Time {
	t.Helper()
	tm, err := domain.ParseTime(s)
	require.NoError(t, err)
	return tm
}

func TestNewDungeon_HappyPath(t *testing.T) {
	t.Parallel()

	d, err := domain.NewDungeon(domain.DungeonConfig{
		Floors: 2, Monsters: 2, OpenAt: mustParseTime(t, "14:00:00"), Duration: 2,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, d.Floors)
	assert.Equal(t, 2, d.Monsters)
}

func TestNewDungeon_RejectsZeroFloors(t *testing.T) {
	t.Parallel()

	_, err := domain.NewDungeon(domain.DungeonConfig{
		Floors: 0, Monsters: 2, OpenAt: mustParseTime(t, "14:00:00"), Duration: 2,
	})
	require.ErrorIs(t, err, domain.ErrBadFloors)
}

func TestNewDungeon_RejectsZeroMonsters(t *testing.T) {
	t.Parallel()

	_, err := domain.NewDungeon(domain.DungeonConfig{
		Floors: 2, Monsters: 0, OpenAt: mustParseTime(t, "14:00:00"), Duration: 2,
	})
	require.ErrorIs(t, err, domain.ErrBadMonsters)
}

func TestNewDungeon_RejectsCrossMidnight(t *testing.T) {
	t.Parallel()

	_, err := domain.NewDungeon(domain.DungeonConfig{
		Floors: 2, Monsters: 2, OpenAt: mustParseTime(t, "23:00:00"), Duration: 2,
	})
	require.ErrorIs(t, err, domain.ErrClosesAfterMidnight)
}

func TestNewDungeon_RejectsBoundaryExactly24h(t *testing.T) {
	t.Parallel()

	_, err := domain.NewDungeon(domain.DungeonConfig{
		Floors: 2, Monsters: 2, OpenAt: mustParseTime(t, "00:00:00"), Duration: 24,
	})
	require.ErrorIs(t, err, domain.ErrClosesAfterMidnight)
}
