package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		repeat VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT ""
	);
	CREATE INDEX scheduler_date ON scheduler(date);
`

var ErrNotInitialized = errors.New("database is not initialized")

type Database struct {
	inner *sql.DB
}

func (self *Database) Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	self.inner = db

	if install {
		_, err := self.inner.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}
