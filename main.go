package main

// TODO: create progresses table into db.

import (
	"GOals/db"
	"GOals/models"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// global variables
var goals *[]models.Goal

func init() {
	// parse command line arguments
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

	return

	//createGoal(promptGoal())

	//goal := (*goals)[len(*goals)-1]
	goal := (*goals)[0]
	val, date, note := promptProgress()
	createProgress(&goal, val, date, note)

	/*
		val, date, note = promptProgress()
		createProgress(&goal, val, date, note)
	*/
	fmt.Println(goal)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func createGoal(name, date, note string) {
	// check length constraints
	if len(name) > 20 || len(note) > 50 {
		panic("Lenght contraints are not respected")
	}
	goal := *(models.CreateGoal(name, note))
	err := goal.SetDate(date)
	checkErr(err)

	// insert goal into database
	id, err := db.InsertGoal(&goal)
	checkErr(err)
	// assign id to goal
	goal.ID = id
}

func createProgress(goal *models.Goal, value int64, date, note string) {
	// check constraints
	if value > 100 || len(note) > 50 {
		panic("Lenght contraints are not respected")
	}

	progress := (models.CreateProgress(value, note))
	err := progress.SetDate(date)
	checkErr(err)

	// add progress to goal
	goal.AddProgress(progress)

	// insert progress into db
	progress.ID, err = db.InsertProgress(progress, goal.ID)
	checkErr(err)
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
