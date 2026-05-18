package domain

type EventID int

const (
	Register       EventID = iota + 1
	Enter
	KillMonster
	NextFloor
	PrevFloor
	EnterBoss
	KillBoss
	Leave
	CannotContinue
	Heal
	Damage
)

const (
	Disqualified   EventID = iota + 31
	Dead
	ImpossibleMove
)

type Event struct {
	Time     Time
	ID       EventID
	PlayerID int
	Extra    string
}
