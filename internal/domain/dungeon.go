package domain

type DungeonConfig struct {
	Floors    int
	Monsters  int
	OpenAt    Time
	Duration  int
}

type Dungeon struct {
	Floors   int
	Monsters int
	OpenAt   Time
	CloseAt  Time
}

func NewDungeon(cfg DungeonConfig) (Dungeon, error) {
	switch {
	case cfg.Floors < 1:
		return Dungeon{}, ErrBadFloors
	case cfg.Monsters < 1:
		return Dungeon{}, ErrBadMonsters
	case cfg.Duration < 1:
		return Dungeon{}, ErrBadDuration
	}
	d := Dungeon{
		Floors:   cfg.Floors,
		Monsters: cfg.Monsters,
		OpenAt:   cfg.OpenAt,
		CloseAt:  cfg.OpenAt.AddHours(cfg.Duration),
	}
	if d.closesAfterMidnight() {
		return Dungeon{}, ErrClosesAfterMidnight
	}
	return d, nil
}

func (d Dungeon) IsOpen(t Time) bool {
	return t.Seconds() >= d.OpenAt.Seconds() && t.Seconds() < d.CloseAt.Seconds()
}

func (d Dungeon) closesAfterMidnight() bool {
	return d.CloseAt.Seconds() >= secondsPerDay
}
