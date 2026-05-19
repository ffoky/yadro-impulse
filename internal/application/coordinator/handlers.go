package coordinator

import (
	"strconv"

	"yadro-impulse/internal/domain"
)

func (c *Coordinator) dispatch(player *domain.Player, event domain.Event) error {
	switch event.ID {
	case domain.Enter:
		return c.echo(player.Enter(event.Time, c.dungeon), event)
	case domain.KillMonster:
		return c.echo(player.KillMonster(event.Time, c.dungeon), event)
	case domain.NextFloor:
		return c.echo(player.NextFloor(event.Time, c.dungeon), event)
	case domain.PrevFloor:
		return c.echo(player.PrevFloor(event.Time, c.dungeon), event)
	case domain.EnterBoss:
		return c.echo(player.EnterBossFloor(event.Time, c.dungeon), event)
	case domain.KillBoss:
		return c.echo(player.KillBoss(event.Time), event)
	case domain.Leave:
		return c.echo(player.LeaveDungeon(event.Time), event)
	case domain.CannotContinue:
		return c.handleCannotContinue(player, event)
	case domain.Heal:
		return c.handleHeal(player, event)
	case domain.Damage:
		return c.handleDamage(player, event)
	}
	return c.impossible(event)
}

func (c *Coordinator) echo(err error, event domain.Event) error {
	if err != nil {
		return c.impossible(event)
	}
	return c.write(event)
}

func (c *Coordinator) handleCannotContinue(player *domain.Player, event domain.Event) error {
	player.Disqualify(event.Time)
	if err := c.write(event); err != nil {
		return err
	}
	return c.write(domain.Event{
		Time: event.Time, ID: domain.Disqualified, PlayerID: event.PlayerID,
	})
}

func (c *Coordinator) handleHeal(player *domain.Player, event domain.Event) error {
	amount, err := strconv.Atoi(event.Extra)
	if err != nil {
		return c.impossible(event)
	}
	if err = player.Heal(amount); err != nil {
		return c.impossible(event)
	}
	return c.write(event)
}

func (c *Coordinator) handleDamage(player *domain.Player, event domain.Event) error {
	amount, err := strconv.Atoi(event.Extra)
	if err != nil {
		return c.impossible(event)
	}
	dead, err := player.TakeDamage(event.Time, amount)
	if err != nil {
		return c.impossible(event)
	}
	if writeErr := c.write(event); writeErr != nil {
		return writeErr
	}
	if dead {
		eventDead := domain.Event{
			Time: event.Time, ID: domain.Dead, PlayerID: event.PlayerID,
		}
		return c.write(eventDead)
	}
	return nil
}
