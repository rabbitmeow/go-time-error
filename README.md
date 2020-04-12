## Overview

This is simple Go REST API to reproduce error when accessing [mysql time data type](https://dev.mysql.com/doc/refman/8.0/en/time.html) using [mysql driver](https://github.com/go-sql-driver/mysql).
Read my explanation here :

- en version : 
- id version :

## How to run

### With docker

- Run your [docker](https://www.docker.com/)
- Clone this repo
- Run `cd go-time-error/after` or `cd go-time-error/before`
- Linux or macOS : run `make start` for starting application and `make stop` for stopping application
- Windows :  `docker-compose up -d && docker image prune` for starting application and `docker-compose down` for taking it down.

### Without docker

- Make sure latest version of Golang is installed and setup properly
- Mysql installed
- Create your database and you can import `prepare_shift.sql`
- Run `cd go-time-error/after` or `cd go-time-error/before`
- Modify the `config.yaml` file
- Run `go run main.go`
