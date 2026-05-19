package output

import (
	"fmt"
	"io"

	"yadro-impulse/internal/domain"
)

func WriteReport(w io.Writer, r domain.Report) error {
	if _, err := fmt.Fprintln(w, "Final report:"); err != nil {
		return fmt.Errorf("write report header: %w", err)
	}
	for _, p := range r.Players {
		_, err := fmt.Fprintf(w, "[%s] %d [%s, %s, %s] HP:%d\n",
			p.Status, p.ID, p.InDungeon, p.AvgFloorClear, p.BossKillTime, p.HP)
		if err != nil {
			return fmt.Errorf("write report line: %w", err)
		}
	}
	return nil
}
