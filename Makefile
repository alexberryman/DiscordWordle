.PHONY: help status up down
.DEFAULT_GOAL := help
status: ## sql-migrate status
	set -o allexport; source local.env; set +o allexport; cd internal/turnips; sql-migrate status
up: ## sql-migrate up
	set -o allexport; source local.env; set +o allexport; cd internal/turnips; sql-migrate up; sql-migrate status
down: ## sql-migrate down
	set -o allexport; source local.env; set +o allexport; cd internal/turnips; sql-migrate down; sql-migrate
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
