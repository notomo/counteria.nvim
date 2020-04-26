test:
	$(MAKE) build
	THEMIS_VIM=nvim THEMIS_ARGS="-e -s --headless" themis

build:
	GO111MODULE=on go build -o ./bin/counteriad ./cmd/counteriad/main.go

DB := $(HOME)/.local/share/counteria/default.db

clear:
	rm $(DB)

db_exec: ARG := .help
db_exec:
	sqlite3 $(DB) '.headers on' '.mode column' '$(ARG)'
	@echo

DB_EXEC := @$(MAKE) --no-print-directory db_exec ARG=
show:
	$(DB_EXEC).tables
	$(DB_EXEC).schema
	$(DB_EXEC)'PRAGMA table_info(tasks);'
	$(DB_EXEC)'PRAGMA table_info(done_tasks);'
	$(DB_EXEC)'PRAGMA table_info(task_rule_lines);'
	$(DB_EXEC)'SELECT * FROM tasks;'
	$(DB_EXEC)'SELECT * FROM done_tasks;'
	$(DB_EXEC)'SELECT * FROM task_rule_lines;'

lint:
	staticcheck ./...
	scopelint --set-exit-status ./...
	golint -set_exit_status ./...
	test -z "`goimports -d ./`" || (echo "`goimports -d ./`"; exit 1)
	swityp -target github.com/notomo/counteria.nvim/src/domain/model.TaskRuleType ./...
	swityp -target github.com/notomo/counteria.nvim/src/domain/model.PeriodUnit ./...

setup:
	cat tools.go | awk -F'"' '/_/ {print $$2}' | xargs -tI {} go install {}

.PHONY: test
.PHONY: build
.PHONY: clear
.PHONY: db_exec
.PHONY: show
.PHONY: lint
.PHONY: setup
