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
