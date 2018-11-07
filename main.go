package main

// TODO: create progresses table into db.

import (
	"GOals/db"
	"GOals/models"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey"
)

// global variables
var goals *[]models.Goal

func init() {
	// parse command line arguments
	/*
		parser := argparse.NewParser("GOals", "Personal goals register and tracker")

		s := parser.List("a", "add", nil)

		err := parser.Parse(os.Args)
		if err != nil {
			fmt.Fprintln(os.Stderr, parser.Usage(err))
		}

		fmt.Println(*s)
		os.Exit(0)
	*/

}

func main() {
	// connect to local database
	err := db.Connect()
	checkErr(err)

	// disconnect from database on exit
	defer db.Disconnect()

	// fetch goals
	goals, err = db.FetchGoalsAndProgress()
	checkErr(err)
	//fmt.Println((*goals))

	for _, v := range *goals {
		fmt.Println(v)
	}

	editGoal(promptSelectGoal(nil))
	return

	/*
		goal := createGoal(promptGoal())
		goal.AddProgress(createProgress(promptProgress()))
		goal.AddProgress(createProgress(promptProgress()))
		db.InsertGoal(goal)
	*/

	/*
		goal := (*goals)[len((*goals))-1]
		goal.AddProgress(createProgress(promptProgress()))

		db.InsertProgress(&goal.Progress[len(goal.Progress)-1], &goal)

		fmt.Println(goal)
	*/
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func promptSelectGoal(matches *[]models.Goal) *models.Goal {
	goal := ""
	options := make([]string, 0, len(goal))

	for _, v := range *goals {
		options = append(options, v.Name)
	}

	prompt := &survey.Select{
		Message: "Choose a goal:",
		Options: options,
	}
	survey.AskOne(prompt, &goal, nil)

	return getGoal(goal)
}

func editGoal(goal *models.Goal) {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Event name:",
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
				Message: "Event date:",
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
				Message: "Event note:",
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

	goal.Name, goal.Date, goal.Note = ans.Name, ans.Date, ans.Note

	fmt.Println(goal)
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
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Name: ")
		scanner.Scan()
		name = scanner.Text()

		if len(name) > 0 && len(name) <= 20 {
			break
		}
		fmt.Fprintln(os.Stderr, "-- length constraint not respected")
	}

	for {
		fmt.Print("Date (dd/mm/yyyy): ")
		scanner.Scan()
		date = scanner.Text()

		var err error
		date, err = models.ParseDate(date)

		if err != nil {
			fmt.Fprintln(os.Stderr, "-- invalid date ")
		} else {
			break
		}
	}

	for {
		fmt.Print("Note: ")
		scanner.Scan()
		note = scanner.Text()

		if len(note) <= 50 {
			break
		}
		fmt.Fprintln(os.Stderr, "-- length constraint not respected")
	}
	return
}

func promptProgress() (value int64, date, note string) {
	scanner := bufio.NewScanner(os.Stdin)
	var err error
	for {
		fmt.Print("Value (0..100): ")
		scanner.Scan()
		value, err = strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "-- invalid value")
			continue
		}
		if value >= 0 && value <= 100 {
			break
		} else {
			fmt.Fprintln(os.Stderr, "-- invalid value")
		}
	}

	for {
		fmt.Print("Date (dd/mm/yyyy): ")
		scanner.Scan()
		date = scanner.Text()

		var err error
		date, err = models.ParseDate(date)

		if err != nil {
			fmt.Fprintln(os.Stderr, "-- invalid date ")
		} else {
			break
		}
	}

	for {
		fmt.Print("Note: ")
		scanner.Scan()
		note = scanner.Text()

		if len(note) <= 50 {
			break
		}
		fmt.Fprintln(os.Stderr, "-- length constraint not respected")
	}
	return
}
