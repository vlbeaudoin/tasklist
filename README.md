# vlbeaudoin/tasklist

Maintains a list of tasks.

## usage

1. In this cloned repository, build the application with `go build .` or `go install .`

2. Fill in the config file (see [config file](#config-file) for file layout.)

3. Run the bare binary (built or installed) to access the help.

## config file

Default config file location is `~/.tasklist.yaml`, use flag `--config` to set a different location.

Supports any file format that `github.com/spf13/viper` supports.

### db.type

Currently, only the `sqlite` type is available, but any type supported by `gorm` could be implemented (postgres is planned).

### db.path

Any path accessible in READ/WRITE by the user running the binary is valid. Points to a location where the database file will reside (sqlite only).

### general.list_after_add

If true, will call the equivalent of `tasklist list` after each call to `tasklist add <New Task Label>`.

## external resources

https://github.com/divrhino/studybuddy

https://github.com/gophercises/task

https://github.com/vlbeaudoin/studybuddy
