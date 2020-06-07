include .env

.PHONY: up
up: docker-compose.yml
	@echo "$@"
	docker-compose up

.PHONY: down
down: docker-compose.yml
	@echo "$@"
	docker-compose stop

.PHONY: clean
clean:
	@echo "$@"
	docker kill $(shell docker ps -a -q) || true
	docker rm $(shell docker ps -a -q) || true

.PHONY: shell
shell:
	@echo "$@"
	docker exec -it app \
		/bin/sh

.PHONY: test
test:
ifndef $(ARGS)
	@echo 'no ARGS around'
	$(eval ARGS := "./tests/...")
endif
	docker exec -it beeg_db /bin/sh -c \
		"mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e 'drop database if exists ${MYSQL_TEST_DATABASE}; create database ${MYSQL_TEST_DATABASE}' "
	docker exec -it -e APP_ENV=test -e MYSQL_DATABASE="${MYSQL_TEST_DATABASE}" app \
			go test ${ARGS} -v

.PHONY: add_migration
add_migration:
	@echo "$@"
	docker exec -it app /bin/sh -c \
		"goose -dir /app/app/migrations/ mysql \"${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/${MYSQL_DATABASE}\" create ${ARGS} sql"

.PHONY: migrate_up
migrate_up:
	@echo "$@"
	docker exec -it app /bin/sh -c \
		"goose -dir /app/app/migrations/ mysql \"${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/${MYSQL_DATABASE}\" up"

.PHONY: migrate_down
migrate_down:
	@echo "$@"
	docker exec -it app /bin/sh -c \
		"goose -dir /app/app/migrations/ mysql \"${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:3306)/${MYSQL_DATABASE}\" down"