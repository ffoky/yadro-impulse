package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDungeon(t *testing.T) Dungeon {
	t.Helper()
	openAt, err := ParseTime("14:00:00")
	require.NoError(t, err)
	d, err := NewDungeon(DungeonConfig{Floors: 2, Monsters: 2, OpenAt: openAt, Duration: 2})
	require.NoError(t, err)
	return d
}

func mustEnter(t *testing.T) *Player {
	t.Helper()
	p := NewPlayer(1)
	require.NoError(t, p.Register())
	d := newDungeon(t)
	enter, _ := ParseTime("14:10:00")
	require.NoError(t, p.Enter(enter, d))
	return p
}

func TestPlayer_RegisterTwiceFails(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	require.NoError(t, p.Register())
	require.ErrorIs(t, p.Register(), ErrAlreadyRegistered)
}

func TestPlayer_EnterWithoutRegister(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	d := newDungeon(t)
	enter, _ := ParseTime("14:10:00")
	require.ErrorIs(t, p.Enter(enter, d), ErrNotRegistered)
}

func TestPlayer_EnterAfterClose(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	require.NoError(t, p.Register())
	d := newDungeon(t)
	late, _ := ParseTime("17:00:00")
	require.ErrorIs(t, p.Enter(late, d), ErrDungeonClosed)
}

func TestPlayer_HealCapsAtMaxHP(t *testing.T) {
	t.Parallel()

	p := mustEnter(t)
	hit, _ := ParseTime("14:11:00")
	_, err := p.TakeDamage(hit, 50)
	require.NoError(t, err)
	require.NoError(t, p.Heal(200))
	assert.Equal(t, 100, p.HP)
}

func TestPlayer_DamageKills(t *testing.T) {
	t.Parallel()

	p := mustEnter(t)
	hit, _ := ParseTime("14:11:00")
	dead, err := p.TakeDamage(hit, 200)
	require.NoError(t, err)
	assert.True(t, dead)
	assert.Equal(t, 0, p.HP)
	assert.Equal(t, StateDead, p.State)
}

func TestPlayer_KillMonsterRecordsFloorTime(t *testing.T) {
	t.Parallel()

	p := mustEnter(t)
	d := newDungeon(t)

	k1, _ := ParseTime("14:12:00")
	require.NoError(t, p.KillMonster(k1, d))
	k2, _ := ParseTime("14:15:00")
	require.NoError(t, p.KillMonster(k2, d))

	require.Len(t, p.timeline.floorClearTimes, 1)
	assert.Equal(t, "00:05:00", p.timeline.floorClearTimes[0].String())
}

func TestPlayer_PrevFloorOnFirstFails(t *testing.T) {
	t.Parallel()

	p := mustEnter(t)
	d := newDungeon(t)
	now, _ := ParseTime("14:11:00")
	require.ErrorIs(t, p.PrevFloor(now, d), ErrNoPrevFloor)
}

func TestPlayer_KillBossOnNonBossFloorFails(t *testing.T) {
	t.Parallel()

	p := mustEnter(t)
	now, _ := ParseTime("14:11:00")
	require.ErrorIs(t, p.KillBoss(now), ErrNotOnBossFloor)
}
