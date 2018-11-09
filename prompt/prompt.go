// Package prompt provides functions for prompting the user to insert something.
package prompt

import (
	"errors"
	"goals/models"
	"strconv"

	"github.com/AlecAivazis/survey"
)

// Confirm prompts the user to answer Yes or No to a particular question.
func Confirm(question string) bool {
	value := false
	p := &survey.Confirm{
		Message: question,
	}
	survey.AskOne(p, &value, nil)
	return value
}

// InsertGoal prompts the user to insert the fields of a new goal.
func InsertGoal(goals *[]models.Goal) (name, date, note string) {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Goal name:",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if !ok || len(str) > 20 || len(str) == 0 {
					return errors.New("Lenght contraint not respected")
				}
				for i := range *goals {
					if (*goals)[i].Name == str {
						return errors.New("goal " + str + " already exists")
					}
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
func InsertProgress(progresses *[]models.Progress) (int64, string, string) {
	qs := []*survey.Question{
		{
			Name: "Value",
			Prompt: &survey.Input{
				Message: "Progress value (0..100):",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				value, err := strconv.ParseInt(str, 10, 64)
				if !ok || err != nil || value < 0 || value > 100 {
					return errors.New("invalid value")
				}
				for i := range *progresses {
					if (*progresses)[i].Value == value {
						return errors.New("progress with value " + str + "% already exists")
					}
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
// in goals array parameter
func SelectGoal(goals *[]models.Goal) *models.Goal {
	goal := ""
	options := make([]string, 0, len(goal))

	for _, v := range *goals {
		options = append(options, v.Name)
	}

	p := &survey.Select{
		Message: "Choose a goal:",
		Options: options,
	}
	survey.AskOne(p, &goal, nil)
	return getGoal(goals, goal)
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

	progressID := selected[1:3]
	if progressID[1] == '%' {
		progressID = progressID[:1]
	}

	progressIDInt, err := strconv.ParseInt(progressID, 10, 64)
	if err != nil {
		panic(err)
	}

	return getProgress(&goal.Progress, progressIDInt)
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

func getGoal(goals *[]models.Goal, goalName string) *models.Goal {
	for i := range *goals {
		if (*goals)[i].Name == goalName {
			return &(*goals)[i]
		}
	}
	return nil
}

func getProgress(progresses *[]models.Progress, progressValue int64) *models.Progress {
	for i := range *progresses {
		if (*progresses)[i].Value == progressValue {
			return &(*progresses)[i]
		}
	}
	return nil
}
