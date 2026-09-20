# GO FINAL PROJECT

This is a task scheduler written in go. It supports repetition rules with days interval, yearly, on specified weekdays and on specified days of the month in optionally specified months.

### Solved tasks with asterisks:

- Port setting with `TODO_PORT` environment variable

- Database setting with `TODO_DBFILE` environment variable

- Week and month repetition rules

- Task search with title, comment or date

- JWT Authentication System

## Using

### System dependencies

- docker >=29.8.0
- docker-buildx >= 0.37.1

### Quickstart

```bash
# This script builds and runs docker container
./build_run.sh
```

### Building from source

```bash
# Building docker image
docker build -t venin/go-final-project

# Launching docker container
docker run -d -p 7540:7540 venin/go-final-project
```

### Environment variables

- `TODO_PORT` - overrides server port
- `TODO_DBFILE` - overrides path to database file
- `TODO_PASSWORD` - sets password for api

## Testing

### Tests

**Launch server before starting any tests!**

Start all tests
```bash
go test ./tests
```

Test frontend
```bash
go test -run '^TestApp$' -count=1 ./tests
```

Test sql database
```bash
go test -run '^TestDB$' -count=1 ./tests
```

Test repetition rules solver
```bash
go test -run '^TestNextDate$' -count=1 ./tests
```

Test `POST` on `api/task`
```bash
go test -run '^TestAddTask$' -count=1 ./tests
```

Test `POST` on `api/task` and `GET` on `api/tasks`
```bash
go test -run '^TestTasks$' -count=1 ./tests
```

Test `PUT` on `api/task`
```bash
go test -run '^TestEditTask$' -count=1 ./tests
```

Test `POST` on `api/task/done`
```bash
go test -run '^TestDone$' -count=1 ./tests
```

Test `DELETE` on `api/task`
```bash
go test -run '^TestDelTask$' -count=1 ./tests
```

### Testing settings

Testing settings is configured through [`tests/settings.go`](tests/settings.go) file.

Parameters:

- `Port` - port where server listens for incoming connections. Overriden by `TODO_PORT` environment variable.
- `DBFile` - path to SQLite database. Overriden by `TODO_DBFILE` environment variable.
- `FullNextDate` - controls whether "week" and "month" repetition rules are tested or not.
- `Search` - controls whether searching api is tested or not.
- `Token` - JWT used to authenticate on server. Leave empty if password is not set.

## Known issues

- If `TODO_DBFILE` is set, tests search for database file inside `tests` folder, use `./...` instead of `./tests`