package db

import (
	"database/sql"
	"os"

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
	ID      int    `json:"id"`
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
