package main

// TODO: setup function
// TODO: Check for len(goals) before edit or remove operations

import (
	"GOals/db"
	"GOals/models"
	"GOals/prompt"
	"flag"
	"fmt"
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

	// parse command line arguments
	var flagNew string
	var flagEdit string
	var flagRemove string

	flag.StringVar(&flagNew, "new", "", "Add a new [goal | progress]")
	flag.StringVar(&flagEdit, "edit", "", "Edit an existing [goal | progress]")
	flag.StringVar(&flagRemove, "remove", "", "Remove an existing [goal | progress]")
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
		return
	}

	if flagEdit != "" {
		switch flagEdit {
		case "goal":
			onEditGoal()
		case "progress":
			onEditProgress()
		default:
			flag.Usage()
		}
		return
	}

	if flagRemove != "" {
		switch flagRemove {
		case "goal":
			onRemoveGoal()
		case "progress":
			onRemoveProgress()
		default:
			flag.Usage()
		}
		return
	}

	onShowGoals()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func onAddGoal() {
	goal := createGoal(prompt.InsertGoal(goals))

	if prompt.Confirm("Add progress?") {
		goal.AddProgress(createProgress(prompt.InsertProgress(&goal.Progress)))
	}

	err := db.InsertGoal(goal)
	checkErr(err)

	fmt.Println("-- goal added")
}

func onAddProgress() {
	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	goal := prompt.SelectGoal(goals)
	progress := createProgress(prompt.InsertProgress(&goal.Progress))
	goal.AddProgress(progress)
	err := db.InsertProgress(progress, goal)
	checkErr(err)

	fmt.Println("-- progress added")
}

func onShowGoals() {
	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	for i := range *goals {
		fmt.Println((*goals)[i])
	}
}

func onEditGoal() {
	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	goal := prompt.SelectGoal(goals)
	goal = prompt.EditGoal(goal)

	err := db.UpdateGoalNoProgress(goal)
	checkErr(err)

	fmt.Println("-- goal edited")
}

func onEditProgress() {
	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	goal := prompt.SelectGoal(goals)
	progress := prompt.SelectProgress(goal)
	progress = prompt.EditProgress(progress)

	err := db.UpdateProgress(progress)
	checkErr(err)

	fmt.Println("-- progress edited")
}

func onRemoveGoal() {
	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	goal := prompt.SelectGoal(goals)
	fmt.Println(goal)
	if prompt.Confirm("Remove goal? (this action cannot be undone)") {
		db.RemoveGoal(goal.ID)
		fmt.Println("-- goal removed")
	} else {
		fmt.Println("-- canceled")
	}
}

func onRemoveProgress() {
	if len((*goals)) == 0 {
		fmt.Println("-- no goals")
		return
	}

	goal := prompt.SelectGoal(goals)

	// check if goal has progresses
	if len((*goal).Progress) == 0 {
		fmt.Println("-- goal has no progress")
		return
	}
	progress := prompt.SelectProgress(goal)
	if prompt.Confirm("Remove progress? (this action cannot be undone)") {
		db.RemoveProgress(progress.ID)
		fmt.Println("-- progress removed")
	} else {
		fmt.Println("-- canceled")
	}
}

func createGoal(name, date, note string) *models.Goal {

	return &models.Goal{
		Name: name,
		Date: date,
		Note: note,
	}
}

func createProgress(value int64, date, note string) *models.Progress {
	return &models.Progress{
		Value: value,
		Date:  date,
		Note:  note,
	}
}
