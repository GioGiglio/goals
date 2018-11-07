package main

// TODO: create progresses table into db.

import (
	"GOals/db"
	"GOals/models"
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey"
)

// global variables
var goals *[]models.Goal

func init() {
	// connect to local database
	err := db.Connect()
	checkErr(err)

	// fetch goals
	goals, err = db.FetchGoalsAndProgress()
	checkErr(err)
}

func main() {
	// disconnect from database on exit
	defer db.Disconnect()

	var flagNew string
	var flagEdit string

	flag.StringVar(&flagNew, "new", "", "Add a new [goal | progress]")
	flag.StringVar(&flagEdit, "edit", "", "Edit an existing [goal | progress]")
	flag.Parse()

	if flagNew != "" {
		switch flagNew {
		case "goal":
			onAddGoal()
		case "progress":
			onAddProgress()
		default:
			flag.Usage()
		}
	}

	if flagEdit != "" {
		switch flagEdit {
		case "goal":
			onEditGoal()
		case "progress":
			fmt.Println("-- not implemented :(")
		default:
			flag.Usage()
		}
	}

	onShowGoals()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func onAddGoal() {
	goal := createGoal(promptGoal())

	shouldAddProgress := false
	prompt := &survey.Confirm{
		Message: "Add progress?",
	}
	survey.AskOne(prompt, &shouldAddProgress, nil)

	if shouldAddProgress {
		goal.AddProgress(createProgress(promptProgress()))
	}
	err := db.InsertGoal(goal)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("-- goal added")
}

func onAddProgress() {
	goal := promptSelectGoal(nil)
	progress := createProgress(promptProgress())
	goal.AddProgress(progress)
	err := db.InsertProgress(progress, goal)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("-- progress added")
}

func onShowGoals() {
	goals, err := db.FetchGoalsAndProgress()
	if err != nil {
		fmt.Println(err)
	}

	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	for i := range *goals {
		fmt.Println((*goals)[i])
	}
}

func onEditGoal() {
	goal := promptSelectGoal(nil)
	goal = promptEditGoal(goal)

	err := db.UpdateGoalNoProgress(goal)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("-- goal edited")
}

func promptSelectGoal(src *[]models.Goal) *models.Goal {
	var source *[]models.Goal
	if src != nil {
		source = src
	} else {
		source = goals
	}
	goal := ""
	options := make([]string, 0, len(goal))

	for _, v := range *source {
		options = append(options, v.Name)
	}

	prompt := &survey.Select{
		Message: "Choose a goal:",
		Options: options,
	}
	survey.AskOne(prompt, &goal, nil)

	return getGoal(goal)
}

func promptEditGoal(goal *models.Goal) *models.Goal {
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

func getGoal(name string) *models.Goal {
	for i := range *goals {
		if (*goals)[i].Name == name {
			return &(*goals)[i]
		}
	}
	return nil
}

func createGoal(name, date, note string) *models.Goal {
	// check length constraints
	if len(name) > 20 || len(note) > 50 {
		panic("Lenght contraints are not respected")
	}
	goal := (models.CreateGoal(name, note))
	err := goal.SetDate(date)
	checkErr(err)
	return goal
}

func createProgress(value int64, date, note string) *models.Progress {
	// check constraints
	if value > 100 || len(note) > 50 {
		panic("Lenght contraints are not respected")
	}

	progress := (models.CreateProgress(value, note))
	err := progress.SetDate(date)
	checkErr(err)
	return progress
}

func promptGoal() (name, date, note string) {
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

func promptProgress() (int64, string, string) {
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
