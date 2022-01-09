.PHONY: help status up down delete-db reset-db
.DEFAULT_GOAL := help
status: ## sql-migrate status
	set -o allexport; source local.env; set +o allexport; cd internal/turnips; sql-migrate status
up: ## sql-migrate up
	set -o allexport; source local.env; set +o allexport; cd internal/turnips; sql-migrate up; sql-migrate status
down: ## sql-migrate down
	set -o allexport; source local.env; set +o allexport; cd internal/turnips; sql-migrate down; sql-migrate status
delete-db: ## deletes docker volume for database and restarts container
	docker-compose stop database
	docker-compose rm -f database
	docker volume rm -f discordwordle_database-data
	docker-compose up -d database
reset-db: ## delete database contents and apply migrations
	make delete-db
	for i in 1 2 3 4 5; do docker-compose exec database pg_isready && break || sleep 1; done
	make up
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
