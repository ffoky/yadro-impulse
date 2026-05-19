package eventsreader_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/infrastructure/eventsreader"
)

func TestReader_Lines(t *testing.T) {
	t.Parallel()

	input := `[14:00:00] 1 1
[14:27:00] 2 11 60
[14:30:00] 1 9 too tired to continue

[14:31:00] 1 8
`

	want := []domain.Event{
		{Time: domain.Time{H: 14}, ID: domain.Register, PlayerID: 1},
		{Time: domain.Time{H: 14, M: 27}, ID: domain.Damage, PlayerID: 2, Extra: "60"},
		{Time: domain.Time{H: 14, M: 30}, ID: domain.CannotContinue, PlayerID: 1, Extra: "too tired to continue"},
		{Time: domain.Time{H: 14, M: 31}, ID: domain.Leave, PlayerID: 1},
	}

	r := eventsreader.New(strings.NewReader(input))
	got := make([]domain.Event, 0, len(want))
	for {
		ev, err := r.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)
		got = append(got, ev)
	}

	assert.Equal(t, want, got)
}

func TestReader_BadLine(t *testing.T) {
	t.Parallel()

	r := eventsreader.New(strings.NewReader("not a valid line\n"))
	_, err := r.Next()
	require.ErrorIs(t, err, eventsreader.ErrBadLine)
}
