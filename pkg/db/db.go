package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"os"
	"time"

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
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Repeat  string `json:"repeat"`
	Comment string `json:"comment"`
}

var sql_db *sql.DB

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

	sql_db = db

	if install {
		_, err := sql_db.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddTask(task *Task) (int64, error) {
	if sql_db == nil {
		return 0, ErrNotInitialized
	}

	if task == nil {
		return 0, ErrTaskIsNil
	}

	var id int64
	query := `INSERT INTO scheduler (date, title, repeat, comment) VALUES (@date, @title, @repeat, @comment);`

	res, err := sql_db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("repeat", task.Repeat),
		sql.Named("comment", task.Comment))
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

func GetTasks(limit int, search string) ([]*Task, error) {
	if sql_db == nil {
		return []*Task{}, ErrNotInitialized
	}

	queryStart := `SELECT * FROM scheduler `
	querySearch := `WHERE title LIKE @search OR comment LIKE @search `
	querySearchDate := `WHERE date = @date `
	queryEnd := `ORDER BY date ASC LIMIT @limit`

	query := queryStart
	var date string
	if len(search) != 0 {
		t, err := time.Parse("02.01.2006", search)
		if err != nil {
			query += querySearch
		} else {
			date = t.Format("20060102")
			query += querySearchDate
		}

	}
	query += queryEnd

	rows, err := sql_db.Query(query,
		sql.Named("search", "%"+search+"%"),
		sql.Named("limit", limit),
		sql.Named("date", date))
	if err != nil {
		return []*Task{}, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0, limit)
	for rows.Next() {
		task := new(Task)
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Repeat, &task.Comment)
		if err != nil {
			return []*Task{}, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return []*Task{}, err
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	if sql_db == nil {
		return nil, ErrNotInitialized
	}

	query := `SELECT * FROM scheduler WHERE id = @id`
	row := sql_db.QueryRow(query, sql.Named("id", id))

	task := new(Task)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Repeat, &task.Comment)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func UpdateTask(task *Task) error {
	if sql_db == nil {
		return ErrNotInitialized
	}

	query := `UPDATE scheduler SET date = @date, title = @title, repeat = @repeat, comment = @comment WHERE id = @id`
	result, err := sql_db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("repeat", task.Repeat),
		sql.Named("comment", task.Comment))
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTask(id string) error {
	if sql_db == nil {
		return ErrNotInitialized
	}

	query := `DELETE FROM scheduler WHERE id = @id`
	result, err := sql_db.Exec(query,
		sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateDate(next string, id string) error {
	if sql_db == nil {
		return ErrNotInitialized
	}

	query := `UPDATE scheduler SET date = @next WHERE id = @id`
	result, err := sql_db.Exec(query,
		sql.Named("id", id),
		sql.Named("next", next))
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
