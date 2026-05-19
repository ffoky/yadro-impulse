package eventsreader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"yadro-impulse/internal/domain"
)

const (
	minTokens       = 2
	tokensWithExtra = 3
)

var ErrBadLine = errors.New("bad event line")

type Reader struct {
	sc   *bufio.Scanner
	line int
}

func New(r io.Reader) *Reader {
	return &Reader{sc: bufio.NewScanner(r)}
}

func (r *Reader) Next() (domain.Event, error) {
	for r.sc.Scan() {
		r.line++
		text := strings.TrimSpace(r.sc.Text())
		if text == "" {
			continue
		}
		ev, err := parse(text)
		if err != nil {
			return domain.Event{}, fmt.Errorf("line %d: %w", r.line, err)
		}
		return ev, nil
	}
	if err := r.sc.Err(); err != nil {
		return domain.Event{}, fmt.Errorf("scan: %w", err)
	}
	return domain.Event{}, io.EOF
}

func parse(line string) (domain.Event, error) {
	t, rest, err := parseTimePrefix(line)
	if err != nil {
		return domain.Event{}, err
	}
	return parseBody(t, rest)
}

func parseTimePrefix(line string) (domain.Time, string, error) {
	if !strings.HasPrefix(line, "[") {
		return domain.Time{}, "", ErrBadLine
	}
	end := strings.IndexByte(line, ']')
	if end < 0 {
		return domain.Time{}, "", ErrBadLine
	}
	t, err := domain.ParseTime(line[1:end])
	if err != nil {
		return domain.Time{}, "", fmt.Errorf("parse time: %w", err)
	}
	return t, strings.TrimSpace(line[end+1:]), nil
}

func parseBody(t domain.Time, body string) (domain.Event, error) {
	parts := strings.SplitN(body, " ", tokensWithExtra)
	if len(parts) < minTokens {
		return domain.Event{}, ErrBadLine
	}
	pid, err := strconv.Atoi(parts[0])
	if err != nil {
		return domain.Event{}, fmt.Errorf("player id: %w", err)
	}
	eid, err := strconv.Atoi(parts[1])
	if err != nil {
		return domain.Event{}, fmt.Errorf("event id: %w", err)
	}
	extra := ""
	if len(parts) == tokensWithExtra {
		extra = parts[2]
	}
	return domain.Event{
		Time:     t,
		ID:       domain.EventID(eid),
		PlayerID: pid,
		Extra:    extra,
	}, nil
}
