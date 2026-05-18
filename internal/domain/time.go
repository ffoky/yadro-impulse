package domain

import "fmt"

const (
	secondsPerHour   = 3600
	secondsPerMinute = 60
	hoursPerDay      = 24
	secondsPerDay    = secondsPerHour * hoursPerDay
	maxHour          = 23
	maxMinute        = 59
	maxSecond        = 59
	timeTokens       = 3
	timeStrLen       = 8
)

type Time struct {
	H, M, S int
}

func ParseTime(s string) (Time, error) {
	if len(s) != timeStrLen {
		return Time{}, fmt.Errorf("parse %q: %w", s, ErrIncorrectTime)
	}
	var t Time
	n, err := fmt.Sscanf(s, "%02d:%02d:%02d", &t.H, &t.M, &t.S)
	if err != nil || n != timeTokens {
		return Time{}, fmt.Errorf("parse %q: %w", s, ErrIncorrectTime)
	}
	if t.H < 0 || t.H > maxHour || t.M < 0 || t.M > maxMinute || t.S < 0 || t.S > maxSecond {
		return Time{}, fmt.Errorf("parse %q: %w", s, ErrIncorrectTime)
	}
	return t, nil
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d:%02d", t.H, t.M, t.S)
}

func (t Time) Seconds() int {
	return t.H*secondsPerHour + t.M*secondsPerMinute + t.S
}

func (t Time) AddHours(h int) Time {
	total := t.Seconds() + h*secondsPerHour
	return Time{
		H: total / secondsPerHour,
		M: (total / secondsPerMinute) % secondsPerMinute,
		S: total % secondsPerMinute,
	}
}

func (t Time) Before(other Time) bool { return t.Seconds() < other.Seconds() }

type Duration int

func Sub(a, b Time) Duration {
	return Duration(a.Seconds() - b.Seconds())
}

func (d Duration) String() string {
	s := int(d)
	if s < 0 {
		s = 0
	}
	return fmt.Sprintf("%02d:%02d:%02d",
		s/secondsPerHour,
		(s/secondsPerMinute)%secondsPerMinute,
		s%secondsPerMinute,
	)
}
