package domain

import "errors"

var (
	ErrIncorrectTime     = errors.New("bad time format")
	ErrBadFloors         = errors.New("bad floors")
	ErrBadMonsters       = errors.New("bad monsters")
	ErrBadDuration       = errors.New("bad duration")
	ErrClosesAfterMidnight = errors.New("dungeon closes after midnight")
	ErrNotRegistered     = errors.New("player not registered")
	ErrAlreadyRegistered = errors.New("player already registered")
	ErrDungeonClosed     = errors.New("dungeon closed")
	ErrAlreadyIn         = errors.New("already in dungeon")
	ErrNotIn             = errors.New("not in dungeon")
	ErrNoMonsters        = errors.New("no monsters on floor")
	ErrFloorNotCleared   = errors.New("floor not cleared")
	ErrNoNextFloor       = errors.New("no next floor")
	ErrNoPrevFloor       = errors.New("no previous floor")
	ErrNotOnBossFloor    = errors.New("not on boss floor")
	ErrOnBossFloor       = errors.New("already on boss floor")
	ErrBossAlreadyDead   = errors.New("boss already killed")
	ErrBadAmount         = errors.New("bad amount")
)
