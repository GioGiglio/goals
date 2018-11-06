package models

import (
	"fmt"
	"strings"
	"time"
)

// Structs

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

// Functions

// CreateGoal creates a base Goal element
func CreateGoal(name, note string) *Goal {
	return &Goal{
		Name:     name,
		Note:     note,
		Date:     "",
		Progress: nil,
		ID:       -1,
	}
}

// CreateProgress creates a base Progress element
func CreateProgress(value int64, note string) *Progress {
	return &Progress{
		Value: value,
		Note:  note,
		Date:  "",
		ID:    -1,
	}
}

// ParseDate parses a date
func ParseDate(date string) (string, error) {
	switch date {
	case "today":
		date = time.Now().Format("02/01/2006")
	case "yesterday":
		date = time.Now().Add(-8.64e+13).Format("02/01/2006")
	default:
		// Check if date is valid
		_, err := time.Parse("02/01/2006", date)
		return date, err
	}
	return date, nil
}

/* Methods */

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
		builder.WriteString(fmt.Sprintf("%-20s\t%-15s\t%-10s \n", g.Name, g.Note, g.Date))
		// write progresses
		for _, v := range g.Progress {
			builder.WriteString("\t" + v.String() + "\n")
		}
		return builder.String()
	}
	return fmt.Sprintf("%-20s\t%-15s\t%-10s \n", g.Name, g.Note, g.Date)
}

func (p Progress) String() string {
	return fmt.Sprintf("[%d%%]\t%-20s\t%-10s", p.Value, p.Note, p.Date)
}
