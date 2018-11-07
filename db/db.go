package db

import (
	"GOals/models"
	"database/sql"
	"errors"

	// This has to be imported in order to use the sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var tableGoals, tableProgresses string
var dbPath string

func init() {
	dbPath = "./goals.db"
	tableGoals, tableProgresses = "goal", "progress"
}

// Connect connects to a local database using sqlite3 drives
func Connect() (err error) {
	db, err = sql.Open("sqlite3", dbPath)
	return
}

// Disconnect disconnects from the local database
func Disconnect() {
	db.Close()
}

// InsertGoal inserts a new goal into local database
// and updates goal fields like ID, and Progress[].ID
func InsertGoal(goal *models.Goal) error {
	// check if there's a connection to the database
	if db == nil {
		return errors.New("No active connection to database")
	}

	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// insert goal into goals table
	res, err := tx.Exec("INSERT INTO goal(name,date,note) VALUES(?,?,?)", goal.Name, goal.Date, goal.Note)
	if err != nil {
		tx.Rollback()
		return err
	}

	goal.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}

	// check if goal has progresses
	if len(goal.Progress) > 0 {

		for i := range goal.Progress {
			res, err = tx.Exec("INSERT INTO progress (goal_id, value, date, note) VALUES(?,?,?,?)",
				goal.ID, goal.Progress[i].Value, goal.Progress[i].Date, goal.Progress[i].Note)

			if err != nil {
				return err
			}

			goal.Progress[i].ID, err = res.LastInsertId()
			if err != nil {
				return err
			}

		}
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// InsertProgress inserts a new progress into local database
// and updates progress field ID to match the one in the database
func InsertProgress(p *models.Progress, g *models.Goal) error {
	// check if there's a connection to the database
	if db == nil {
		return errors.New("No active connection to database")
	}

	// prepare statement
	stmt, err := db.Prepare("INSERT INTO progress (goal_id, value, date, note) values (?,?,?,?)")
	if err != nil {
		return err
	}

	// exec statement
	res, err := stmt.Exec(g.ID, p.Value, p.Date, p.Note)
	if err != nil {
		return err
	}

	// update progress id
	p.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

// FetchGoals fetches goals from local database
func FetchGoals() (*[]models.Goal, error) {
	// check if there's a connection to the database
	if db == nil {
		return nil, errors.New("No active connection to database")
	}

	// exec query
	rows, err := db.Query("SELECT * FROM " + tableGoals)
	if err != nil {
		return nil, err
	}

	// parse goals
	goals := make([]models.Goal, 0)
	var id int64
	var name, date, note string

	for rows.Next() {
		err = rows.Scan(&id, &name, &date, &note)
		if err != nil {
			return nil, err
		}
		goals = append(goals, models.Goal{
			ID: id, Name: name, Date: date, Note: note, Progress: nil,
		})
	}
	return &goals, nil
}

// FetchGoalsAndProgress fetches all goals and progresses from local database
func FetchGoalsAndProgress() (*[]models.Goal, error) {
	// check if there's a connection to the database
	if db == nil {
		return nil, errors.New("No active connection to database")
	}

	// exec query
	rows, err := db.Query("SELECT g.*, p.id, p.value, p.date, p.note FROM goal as g LEFT JOIN progress as p ON g.id = p.goal_id;")
	if err != nil {
		return nil, err
	}

	// parse goals
	goals := make([]models.Goal, 0)

	// goal fields
	var goalID int64
	var goalName, goalDate, goalNote string

	// progress fields
	var progressID, progressValue sql.NullInt64
	var progressDate, progressNote sql.NullString

	var goal models.Goal
	var lastGoal *models.Goal
	var lastGoalID int64 = -1
	for rows.Next() {
		err = rows.Scan(&goalID, &goalName, &goalDate, &goalNote, &progressID,
			&progressValue, &progressDate, &progressNote)

		if err != nil {
			return nil, err
		}

		if lastGoalID != goalID {
			// new goal

			if lastGoalID != -1 {
				// if this is not the first iteration
				// push previous goal to goals
				goals = append(goals, (*lastGoal))
			}

			// create new goal
			goal = models.Goal{
				ID:       goalID,
				Name:     goalName,
				Date:     goalDate,
				Note:     goalNote,
				Progress: make([]models.Progress, 0),
			}

			// if current goal has a progress
			if progressID.Valid {
				// create progress
				progress := models.Progress{
					ID:    must(progressID.Value()).(int64),
					Value: must(progressValue.Value()).(int64),
					Date:  must(progressDate.Value()).(string),
					Note:  must(progressNote.Value()).(string),
				}

				// add progress to goal
				goal.Progress = append(goal.Progress, progress)
			}
		} else {
			// goal is still the same

			// create progress
			progress := models.Progress{
				ID:    must(progressID.Value()).(int64),
				Value: must(progressValue.Value()).(int64),
				Date:  must(progressDate.Value()).(string),
				Note:  must(progressNote.Value()).(string),
			}

			// add progress to previous goal
			(*lastGoal).Progress = append((*lastGoal).Progress, progress)
		}

		lastGoalID = goalID
		lastGoal = &goal
	}
	// push the last goal to goals
	goals = append(goals, (*lastGoal))

	return &goals, nil
}

// must is used to return only the first value
func must(value interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	// v := value.(int64)
	switch v := value.(type) {
	case int64:
		return int64(v)
	case int:
		return int(v)
	case string:
		return string(v)
	default:
		return v.(string)
	}

}