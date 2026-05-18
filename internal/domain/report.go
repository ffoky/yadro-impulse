package domain

import "sort"

type Report struct {
	Players []PlayerReport
}

type PlayerReport struct {
	ID            int
	Status        string
	InDungeon     Duration
	AvgFloorClear Duration
	BossKillTime  Duration
	HP            int
}

func MakeReport(players map[int]*Player) Report {
	ids := make([]int, 0, len(players))
	for id := range players {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	r := Report{Players: make([]PlayerReport, 0, len(ids))}
	for _, id := range ids {
		r.Players = append(r.Players, playerReport(players[id]))
	}
	return r
}

func playerReport(p *Player) PlayerReport {
	pr := PlayerReport{ID: p.ID, HP: p.HP}

	switch p.State {
	case StateSuccess:
		pr.Status = "SUCCESS"
	case StateDisqualified, StateUnknown, StateRegistered:
		pr.Status = "DISQUAL"
	default:
		pr.Status = "FAIL"
	}

	pr.InDungeon = Sub(p.timeline.exitedDungeonAt, p.timeline.enteredDungeonAt)

	if n := len(p.timeline.floorClearTimes); n > 0 {
		var total Duration
		for _, d := range p.timeline.floorClearTimes {
			total += d
		}
		pr.AvgFloorClear = Duration(int(total) / n)
	}

	pr.BossKillTime = p.timeline.bossKillTime
	return pr
}
