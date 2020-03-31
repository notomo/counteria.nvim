test:
	$(MAKE) build
	THEMIS_VIM=nvim THEMIS_ARGS="-e -s --headless" themis

build:
	GO111MODULE=on go build -o ./bin/counteriad ./cmd/counteriad/main.go

clear:
	rm $(HOME)/.local/share/counteria/default.db

.PHONY: test
