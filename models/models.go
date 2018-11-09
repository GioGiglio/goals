// Package models implements the basic structures used among the program.
package models

import (
	"fmt"
	"strings"
	"time"
)

/* STRUCTS */

// Progress to a specific Goal
type Progress struct {
	ID, Value  int64
	Date, Note string
}

// Goal represents a personal objective
type Goal struct {
	ID               int64
	Name, Date, Note string
	Progress         []Progress
}

/* FUNCTIONS */

// ParseDate parses a date
func ParseDate(date string) (string, error) {
	switch date {
	case "today":
		date = time.Now().Format("02/01/2006")
	case "yesterday":
		date = time.Now().Add(-8.64e+13).Format("02/01/2006")
	default:
		// Check if date is valid
		_, err := time.Parse("2/1/2006", date)
		return date, err
	}
	return date, nil
}

/* METHODS */

// AddProgress adds a progress to a Goal
func (g *Goal) AddProgress(progress *Progress) {
	g.Progress = append(g.Progress, *progress)
}

// SetDate sets the date of a Goal
func (g *Goal) SetDate(date string) error {
	date, err := ParseDate(date)
	if err != nil {
		return err
	}
	g.Date = date
	return nil
}

// SetDate sets the date of a Progress
func (p *Progress) SetDate(date string) error {
	date, err := ParseDate(date)
	if err != nil {
		return err
	}
	p.Date = date
	return nil
}

func (g Goal) String() string {
	if len(g.Progress) > 0 {
		// If goals has at least 1 progress, print goal and its progresses
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("%-20s\t%-10s\t%s \n", g.Name, g.Date, g.Note))
		// write progresses
		for _, v := range g.Progress {
			builder.WriteString("\t" + v.String() + "\n")
		}
		return builder.String()
	}
	return fmt.Sprintf("%-20s\t%-10s\t%s \n", g.Name, g.Date, g.Note)
}

func (p Progress) String() string {
	return fmt.Sprintf("[%d%%]\t%-10s\t%s", p.Value, p.Date, p.Note)
}
