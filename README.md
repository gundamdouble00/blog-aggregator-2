# [Blog Aggregator](https://www.boot.dev/courses/build-blog-aggregator-golang)

## Requisite

- [PostgreSQL](https://www.postgresql.org/)
- [Go](https://go.dev/) version 1.27+
- [sqlc](https://sqlc.dev/) version 1.31+
- [goose](https://pressly.github.io/goose/) version 3.28+

## Set up

1. Clone this repo
2. Create a database named `gator2`
3. Create a config file named _.gatorconfig.json_ at the home directory, this file contains a db url and username
4. Migrate the database with `goose`. The migration must be happend in `sql/schema` directory

## Commands

| Name     | Usage                 |
| -------- | --------------------- |
| login    | login \<username\>    |
| register | register \<username\> |
