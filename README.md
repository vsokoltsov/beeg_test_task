# Beeg Test Task


## Before run

* Create file `.env` in root directory with environment keys. List of them provided in [.env.sample](.env.sample)

## Run

* `make up`

## Stop

* `make down`

## Test

* `make test` (make sure,  that the app is running via `make up`)

## Migrations

* `make migrate_up`
* `make migrate_down`
* `make add_migration ARGS=<name of the migration>`

## Benchmark

* [Here](benchmark.txt)