package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

var schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(512) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX  scheduler_date_idx ON scheduler(date DESC);
`

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Init(dbFile string) error {
	var install bool

	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if install {

		row := db.QueryRow(schema)

		err = row.Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func Close() error {
	return db.Close()
}

func AddTask(task Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

	var id int64
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func Tasks(limit int, search string) ([]Task, error) {
	qArgs := []interface{}{sql.Named("limit", limit)}
	query := `SELECT id, date, title, comment, repeat FROM scheduler`

	if search != "" {
		s, err := time.Parse("02.01.2006", search)
		if err != nil {
			qArgs = append(qArgs, sql.Named("search", "%"+search+"%"))
			query = fmt.Sprintln(query, `WHERE title LIKE :search OR comment LIKE :search`)
		} else {
			qArgs = append(qArgs, sql.Named("search", s.Format("20060102")))
			query = fmt.Sprintln(query, `WHERE date = :search`)
		}
	}

	query = fmt.Sprintln(query, `LIMIT :limit`)

	rows, err := db.Query(query, qArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
