// Package prompt provides functions for prompting the user to insert something.
package prompt

import (
	"GOals/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey"
)

// ShouldAddProgress asks the user wheter or not to add a new progress.
// Note: This function gets invoked after user have inserted a new goal.
func ShouldAddProgress() bool {
	value := false
	p := &survey.Confirm{
		Message: "Add progress?",
	}
	survey.AskOne(p, &value, nil)
	return value
}

// InsertGoal prompts the user to insert the fields of a new goal.
func InsertGoal() (name, date, note string) {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Goal name:",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) > 20 || len(str) == 0 {
					return errors.New("Lenght contraint not respected")
				}
				return nil
			},
		},
		{
			Name: "Date",
			Prompt: &survey.Input{
				Message: "Goal date:",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if _, err := models.ParseDate(str); !ok || err != nil {
					return errors.New("Invalid date format")
				}
				return nil
			},
		},
		{
			Name: "Note",
			Prompt: &survey.Input{
				Message: "Goal note:",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) > 50 {
					return errors.New("Lenght contraint not respected")
				}
				return nil
			},
		},
	}

	ans := struct {
		Name, Date, Note string
	}{}

	err := survey.Ask(qs, &ans)
	if err != nil {
		panic(err)
	}

	ans.Date, _ = models.ParseDate(ans.Date)
	return ans.Name, ans.Date, ans.Note
}

// InsertProgress prompts the user to insert the fields of a new progress.
func InsertProgress() (int64, string, string) {
	qs := []*survey.Question{
		{
			Name: "Value",
			Prompt: &survey.Input{
				Message: "Progress value (0..100):",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if val, err := strconv.ParseInt(str, 10, 64); !ok || err != nil || val < 0 || val > 100 {
					return errors.New("invalid value")
				}
				return nil
			},
		},
		{
			Name: "Date",
			Prompt: &survey.Input{
				Message: "Progress date:",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if _, err := models.ParseDate(str); !ok || err != nil {
					return errors.New("Invalid date format")
				}
				return nil
			},
		},
		{
			Name: "Note",
			Prompt: &survey.Input{
				Message: "Progress note:",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) > 50 {
					return errors.New("Lenght contraint not respected")
				}
				return nil
			},
		},
	}

	ans := struct {
		Value      int64
		Date, Note string
	}{}

	err := survey.Ask(qs, &ans)
	if err != nil {
		panic(err)
	}

	ans.Date, _ = models.ParseDate(ans.Date)
	return ans.Value, ans.Date, ans.Note
}

// SelectGoal prompts the user to select one of the goals
// in source parameter
func SelectGoal(source *[]models.Goal) *models.Goal {
	goal := ""
	options := make([]string, 0, len(goal))

	for _, v := range *source {
		options = append(options, v.Name)
	}

	p := &survey.Select{
		Message: "Choose a goal:",
		Options: options,
	}
	survey.AskOne(p, &goal, nil)

	// TODO: remove
	fmt.Println("-- prompt FilterMessage: ")

	return &(*source)[p.SelectedIndex]
}

// SelectProgress prompts the user to select one of the
// progresses of parameter goal
func SelectProgress(goal *models.Goal) *models.Progress {
	options := make([]string, 0, len((*goal).Progress))
	var selected string
	for _, v := range (*goal).Progress {
		options = append(options, v.String())
	}

	p := &survey.Select{
		Message: "Choose a progress",
		Options: options,
	}

	survey.AskOne(p, &selected, nil)

	return &goal.Progress[p.SelectedIndex]
}

// EditGoal prompts the user to re-insert goal fields
func EditGoal(goal *models.Goal) *models.Goal {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Goal name:",
				Default: (*goal).Name,
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) > 20 {
					return errors.New("Lenght contraint not respected")
				}
				return nil
			},
		},
		{
			Name: "Date",
			Prompt: &survey.Input{
				Message: "Goal date:",
				Default: (*goal).Date,
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if _, err := models.ParseDate(str); !ok || err != nil {
					return errors.New("Invalid date format")
				}
				return nil
			},
		},
		{
			Name: "Note",
			Prompt: &survey.Input{
				Message: "Goal note:",
				Default: (*goal).Note,
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) > 50 {
					return errors.New("Lenght contraint not respected")
				}
				return nil
			},
		},
	}

	ans := struct {
		Name, Date, Note string
	}{}

	err := survey.Ask(qs, &ans)
	if err != nil {
		panic(err)
	}

	goal.Name, goal.Note = ans.Name, ans.Note
	goal.Date, _ = models.ParseDate(ans.Date)

	return goal
}

// EditProgress prompts the user to re-insert progress fields.
func EditProgress(p *models.Progress) *models.Progress {
	qs := []*survey.Question{
		{
			Name: "value",
			Prompt: &survey.Input{
				Message: "Progress value (0..100):",
				Default: strconv.FormatInt(p.Value, 10),
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if val, err := strconv.ParseInt(str, 10, 64); !ok || err != nil || val < 0 || val > 100 {
					return errors.New("invalid value")
				}
				return nil
			},
		},
		{
			Name: "Date",
			Prompt: &survey.Input{
				Message: "Progress date:",
				Default: (*p).Date,
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if _, err := models.ParseDate(str); !ok || err != nil {
					return errors.New("Invalid date format")
				}
				return nil
			},
		},
		{
			Name: "Note",
			Prompt: &survey.Input{
				Message: "Progress note:",
				Default: (*p).Note,
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); !ok || len(str) > 50 {
					return errors.New("Lenght contraint not respected")
				}
				return nil
			},
		},
	}

	ans := struct {
		Value      int64
		Date, Note string
	}{}

	err := survey.Ask(qs, &ans)
	if err != nil {
		panic(err)
	}

	ans.Date, _ = models.ParseDate(ans.Date)

	p.Value, p.Date, p.Note = ans.Value, ans.Date, ans.Note

	return p
}
