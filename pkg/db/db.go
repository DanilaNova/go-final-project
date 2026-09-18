package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"os"

	_ "modernc.org/sqlite"
)

const DateFormat = "20060102"

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
var ErrTaskIsNil = errors.New("task is nil")

type Database struct {
	inner *sql.DB
}

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Repeat  string `json:"repeat"`
	Comment string `json:"comment"`
}

var DB Database

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	install := false
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	DB.inner = db

	if install {
		_, err := DB.inner.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddTask(task *Task) (int64, error) {
	if task == nil {
		return 0, ErrTaskIsNil
	}

	var id int64
	query := `INSERT INTO scheduler (date, title, repeat, comment) VALUES (@date, @title, @repeat, @comment);`

	res, err := DB.inner.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("repeat", task.Repeat),
		sql.Named("comment", task.Comment))
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}
