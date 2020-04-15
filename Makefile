test:
	$(MAKE) build
	THEMIS_VIM=nvim THEMIS_ARGS="-e -s --headless" themis

build:
	GO111MODULE=on go build -o ./bin/counteriad ./cmd/counteriad/main.go

DB := $(HOME)/.local/share/counteria/default.db

clear:
	rm $(DB)

show:
	sqlite3 $(DB) .tables
	@echo
	sqlite3 $(DB) .schema
	@echo
	sqlite3 $(DB) 'select * from tasks;'
	@echo
	sqlite3 $(DB) 'select * from done_tasks;'
	@echo
	sqlite3 $(DB) 'select * from task_rule_lines;'

.PHONY: test
.PHONY: build
.PHONY: clear
.PHONY: show
