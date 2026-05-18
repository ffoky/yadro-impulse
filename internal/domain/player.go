package domain

const maxHP = 100

type State int

const (
	StateUnknown State = iota
	StateRegistered
	StateInDungeon
	StateOnBossFloor
	StateLeft
	StateDead
	StateDisqualified
	StateSuccess
	StateFail
)

type floorProgress struct {
	enteredAt Time
	killed    int
}

type timeline struct {
	enteredDungeonAt Time
	exitedDungeonAt  Time
	floorClearTimes  []Duration
	bossEnteredAt    Time
	bossKilled       bool
	bossKillTime     Duration
}

type Player struct {
	ID       int
	State    State
	HP       int
	floor    int
	cleared  []bool
	progress floorProgress
	timeline timeline
}

func NewPlayer(id int) *Player {
	return &Player{ID: id, HP: maxHP}
}

func (p *Player) IsInGame() bool {
	return p.State == StateInDungeon || p.State == StateOnBossFloor
}

func (p *Player) Register() error {
	if p.State != StateUnknown {
		return ErrAlreadyRegistered
	}
	p.State = StateRegistered
	return nil
}

func (p *Player) Enter(t Time, d Dungeon) error {
	if p.State == StateUnknown {
		return ErrNotRegistered
	}
	if p.State != StateRegistered {
		return ErrAlreadyIn
	}
	if !d.IsOpen(t) {
		return ErrDungeonClosed
	}
	p.State = StateInDungeon
	p.HP = maxHP
	p.floor = 1
	p.cleared = make([]bool, d.Floors)
	p.progress = floorProgress{enteredAt: t}
	p.timeline = timeline{enteredDungeonAt: t}
	return nil
}

func (p *Player) KillMonster(t Time, d Dungeon) error {
	if !p.IsInGame() {
		return ErrNotIn
	}
	if p.State == StateOnBossFloor || p.floor == d.Floors {
		return ErrNoMonsters
	}
	if p.progress.killed >= d.Monsters {
		return ErrNoMonsters
	}
	p.progress.killed++
	if p.progress.killed == d.Monsters && !p.cleared[p.floor-1] {
		p.cleared[p.floor-1] = true
		p.timeline.floorClearTimes = append(p.timeline.floorClearTimes, Sub(t, p.progress.enteredAt))
	}
	return nil
}

func (p *Player) NextFloor(t Time, d Dungeon) error {
	if !p.IsInGame() {
		return ErrNotIn
	}
	if p.State == StateOnBossFloor || p.floor >= d.Floors {
		return ErrNoNextFloor
	}
	if !p.cleared[p.floor-1] {
		return ErrFloorNotCleared
	}
	p.floor++
	killed := 0
	if p.floor-1 < len(p.cleared) && p.cleared[p.floor-1] {
		killed = d.Monsters
	}
	p.progress = floorProgress{enteredAt: t, killed: killed}
	return nil
}

func (p *Player) PrevFloor(t Time, d Dungeon) error {
	if !p.IsInGame() {
		return ErrNotIn
	}
	if p.State == StateOnBossFloor || p.floor <= 1 {
		return ErrNoPrevFloor
	}
	p.floor--
	p.progress = floorProgress{enteredAt: t, killed: d.Monsters}
	return nil
}

func (p *Player) EnterBossFloor(t Time, d Dungeon) error {
	if !p.IsInGame() {
		return ErrNotIn
	}
	if p.State == StateOnBossFloor {
		return ErrOnBossFloor
	}
	if p.floor != d.Floors {
		return ErrNotOnBossFloor
	}
	p.State = StateOnBossFloor
	p.timeline.bossEnteredAt = t
	return nil
}

func (p *Player) KillBoss(t Time) error {
	if p.State != StateOnBossFloor {
		return ErrNotOnBossFloor
	}
	if p.timeline.bossKilled {
		return ErrBossAlreadyDead
	}
	p.timeline.bossKilled = true
	p.timeline.bossKillTime = Sub(t, p.timeline.bossEnteredAt)
	return nil
}

func (p *Player) LeaveDungeon(t Time) error {
	if !p.IsInGame() {
		return ErrNotIn
	}
	p.timeline.exitedDungeonAt = t
	if p.timeline.bossKilled {
		p.State = StateSuccess
	} else {
		p.State = StateLeft
	}
	return nil
}

func (p *Player) Heal(amount int) error {
	if !p.IsInGame() {
		return ErrNotIn
	}
	if amount <= 0 {
		return ErrBadAmount
	}
	p.HP += amount
	if p.HP > maxHP {
		p.HP = maxHP
	}
	return nil
}

func (p *Player) TakeDamage(t Time, amount int) (dead bool, err error) {
	if !p.IsInGame() {
		return false, ErrNotIn
	}
	if amount <= 0 {
		return false, ErrBadAmount
	}
	p.HP -= amount
	if p.HP <= 0 {
		p.HP = 0
		p.State = StateDead
		p.timeline.exitedDungeonAt = t
		return true, nil
	}
	return false, nil
}

func (p *Player) Disqualify(t Time) {
	if p.IsInGame() {
		p.timeline.exitedDungeonAt = t
	}
	p.State = StateDisqualified
}

func (p *Player) ForceFail(t Time) {
	p.timeline.exitedDungeonAt = t
	p.State = StateFail
}
