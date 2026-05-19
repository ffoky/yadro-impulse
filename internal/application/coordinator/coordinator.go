package coordinator

import (
	"fmt"
	"strconv"

	"yadro-impulse/internal/domain"
)

type EventWriter interface {
	Write(domain.Event) error
}

type Coordinator struct {
	dungeon  domain.Dungeon
	players  map[int]*domain.Player
	out      EventWriter
	closed bool
}

func New(d domain.Dungeon, out EventWriter) *Coordinator {
	return &Coordinator{
		dungeon: d,
		players: make(map[int]*domain.Player),
		out:     out,
	}
}

func (c *Coordinator) Handle(event domain.Event) error {
	if !c.closed && !event.Time.Before(c.dungeon.CloseAt) {
		c.close()
	}

	if event.ID == domain.Register {
		return c.handleRegister(event)
	}

	player, ok := c.players[event.PlayerID]
	if !ok {
		player = domain.NewPlayer(event.PlayerID)
		player.Disqualify(event.Time)
		c.players[event.PlayerID] = player
		return c.write(domain.Event{
			Time: event.Time, ID: domain.Disqualified, PlayerID: event.PlayerID,
		})
	}

	if player.State == domain.StateDisqualified {
		return nil
	}

	return c.dispatch(player, event)
}

func (c *Coordinator) CloseDungeon() {
	c.close()
}

func (c *Coordinator) Report() domain.Report {
	return domain.MakeReport(c.players)
}

func (c *Coordinator) close() {
	if c.closed {
		return
	}
	c.closed = true
	for _, player := range c.players {
		if player.IsInGame() {
			player.ForceFail(c.dungeon.CloseAt)
		}
	}
}

func (c *Coordinator) handleRegister(event domain.Event) error {
	if _, ok := c.players[event.PlayerID]; ok {
		return c.impossible(event)
	}
	player := domain.NewPlayer(event.PlayerID)
	if err := player.Register(); err != nil {
		return c.impossible(event)
	}
	c.players[event.PlayerID] = player
	return c.write(event)
}

func (c *Coordinator) impossible(event domain.Event) error {
	return c.write(domain.Event{
		Time:     event.Time,
		ID:       domain.ImpossibleMove,
		PlayerID: event.PlayerID,
		Extra:    strconv.Itoa(int(event.ID)),
	})
}

func (c *Coordinator) write(event domain.Event) error {
	if err := c.out.Write(event); err != nil {
		return fmt.Errorf("write event: %w", err)
	}
	return nil
}
