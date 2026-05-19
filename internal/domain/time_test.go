package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yadro-impulse/internal/domain"
)

func TestParseTime(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   string
		want    domain.Time
		wantErr bool
	}{
		"valid":           {"14:05:00", domain.Time{H: 14, M: 5, S: 0}, false},
		"midnight":        {"00:00:00", domain.Time{H: 0, M: 0, S: 0}, false},
		"max":             {"23:59:59", domain.Time{H: 23, M: 59, S: 59}, false},
		"hour overflow":   {"24:00:00", domain.Time{}, true},
		"minute overflow": {"00:60:00", domain.Time{}, true},
		"missing zero":    {"14:5:0", domain.Time{}, true},
		"not a time":      {"hello", domain.Time{}, true},
		"empty":           {"", domain.Time{}, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.ParseTime(tc.input)
			if tc.wantErr {
				require.ErrorIs(t, err, domain.ErrIncorrectTime)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTime_String(t *testing.T) {
	t.Parallel()

	in := domain.Time{H: 4, M: 7, S: 9}
	assert.Equal(t, "04:07:09", in.String())
}

func TestTime_AddHours(t *testing.T) {
	t.Parallel()

	base := domain.Time{H: 14, M: 5, S: 0}
	plus2 := base.AddHours(2)

	assert.Equal(t, domain.Time{H: 16, M: 5, S: 0}, plus2)
	assert.True(t, base.Before(plus2))
}

func TestSubAndDurationString(t *testing.T) {
	t.Parallel()

	a := domain.Time{H: 14, M: 40, S: 0}
	b := domain.Time{H: 15, M: 4, S: 0}

	assert.Equal(t, "00:24:00", domain.Sub(b, a).String())
	assert.Equal(t, "00:00:00", domain.Sub(a, b).String())
}
