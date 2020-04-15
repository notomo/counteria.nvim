test:
	$(MAKE) build
	THEMIS_VIM=nvim THEMIS_ARGS="-e -s --headless" themis

build:
	GO111MODULE=on go build -o ./bin/counteriad ./cmd/counteriad/main.go

DB := $(HOME)/.local/share/counteria/default.db

clear:
	rm $(DB)

ARG := .help
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

.PHONY: test
.PHONY: build
.PHONY: clear
.PHONY: db_exec
.PHONY: show
