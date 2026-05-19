package output

import (
	"fmt"
	"io"

	"yadro-impulse/internal/domain"
)

type Writer struct {
	w io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (out *Writer) Write(ev domain.Event) error {
	text := format(ev)
	if text == "" {
		return nil
	}
	if _, err := fmt.Fprintf(out.w, "[%s] %s\n", ev.Time, text); err != nil {
		return fmt.Errorf("write line: %w", err)
	}
	return nil
}

var messages = map[domain.EventID]string{
	domain.Register:       "Player [%d] registered",
	domain.Enter:          "Player [%d] entered the dungeon",
	domain.KillMonster:    "Player [%d] killed the monster",
	domain.NextFloor:      "Player [%d] went to the next floor",
	domain.PrevFloor:      "Player [%d] went to the previous floor",
	domain.EnterBoss:      "Player [%d] entered the boss's floor",
	domain.KillBoss:       "Player [%d] killed the boss",
	domain.Leave:          "Player [%d] left the dungeon",
	domain.Disqualified:   "Player [%d] is disqualified",
	domain.Dead:           "Player [%d] is dead",
	domain.CannotContinue: "Player [%d] cannot continue due to [%s]",
	domain.Heal:           "Player [%d] has restored [%s] of health",
	domain.Damage:         "Player [%d] recieved [%s] of damage",
	domain.ImpossibleMove: "Player [%d] makes imposible move [%s]",
}

func format(ev domain.Event) string {
	tmpl, ok := messages[ev.ID]
	if !ok {
		return ""
	}
	if ev.Extra != "" {
		return fmt.Sprintf(tmpl, ev.PlayerID, ev.Extra)
	}
	return fmt.Sprintf(tmpl, ev.PlayerID)
}
