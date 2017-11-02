.PHONY: clean

GOBIN = build/bin
GO ?= lastest
LINT_FOLDERS := extkeys ./ stats/...
TESTPKG := $(shell go list ./stats/...)

# This is a code for automatic help generator.
# It supports ANSI colors and categories.
# To add new item into help output, simply add comments
# starting with '##'. To add category, use @category.
GREEN  := $(shell tput -Txterm setaf 2)
WHITE  := $(shell tput -Txterm setaf 7)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)
HELP_FUN = \
		   %help; \
		   while(<>) { push @{$$help{$$2 // 'options'}}, [$$1, $$3] if /^([a-zA-Z\-]+)\s*:.*\#\#(?:@([a-zA-Z\-]+))?\s(.*)$$/ }; \
		   print "Usage: make [target]\n\n"; \
		   for (sort keys %help) { \
			   print "${WHITE}$$_:${RESET}\n"; \
			   for (@{$$help{$$_}}) { \
				   $$sep = " " x (32 - length $$_->[0]); \
				   print "  ${YELLOW}$$_->[0]${RESET}$$sep${GREEN}$$_->[1]${RESET}\n"; \
			   }; \
			   print "\n"; \
		   }

help: ##@tasks Show this help
	@perl -e '$(HELP_FUN)' $(MAKEFILE_LIST)

build: ##@tasks Build msgstat binary into build/bin
	go build -o $(GOBON)/msgstat -v ./

ci: lint test-coverage test-units ##@tasks Runs code linting, code coverage and tests for project

clean: ##@tasks Cleanup
	rm -rf build/bin/*

test-units:
	go test $(TESTPKG)

test-coverage:
	go test -coverpkg= $(TESTPKG)


lint: lint-vet lint-gofmt lint-deadcode lint-misspell lint-unparam lint-unused lint-gocyclo lint-errcheck lint-ineffassign lint-interfacer lint-unconvert lint-staticcheck lint-goconst lint-gas lint-varcheck lint-structcheck lint-gosimple

lint-deps:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

lint-vet:
	@echo "lint-vet"
	@gometalinter  --disable-all --enable=vet --deadline=45s  $(LINT_FOLDERS)
lint-golint:
	@echo "lint-golint"
	@gometalinter  --disable-all --enable=golint --deadline=45s  $(LINT_FOLDERS)
lint-gofmt:
	@echo "lint-gofmt"
	@gometalinter  --disable-all --enable=gofmt --deadline=45s  $(LINT_FOLDERS)
lint-deadcode:
	@echo "lint-deadcode"
	@gometalinter  --disable-all --enable=deadcode --deadline=45s  $(LINT_FOLDERS)
lint-misspell:
	@echo "lint-misspell"
	@gometalinter  --disable-all --enable=misspell --deadline=45s  $(LINT_FOLDERS)
lint-unparam:
	@echo "lint-unparam"
	@gometalinter  --disable-all --enable=unparam --deadline=45s  $(LINT_FOLDERS)
lint-unused:
	@echo "lint-unused"
	@gometalinter  --disable-all --enable=unused --deadline=45s  $(LINT_FOLDERS)
lint-gocyclo:
	@echo "lint-gocyclo"
	@gometalinter  --disable-all --enable=gocyclo --cyclo-over=16 --deadline=45s  $(LINT_FOLDERS)
lint-errcheck:
	@echo "lint-errcheck"
	@gometalinter  --disable-all --enable=errcheck --deadline=1m  $(LINT_FOLDERS)
lint-ineffassign:
	@echo "lint-ineffassign"
	@gometalinter  --disable-all --enable=ineffassign --deadline=45s  $(LINT_FOLDERS)
lint-interfacer:
	@echo "lint-interfacer"
	@gometalinter  --disable-all --enable=interfacer --deadline=45s  $(LINT_FOLDERS)
lint-unconvert:
	@echo "lint-unconvert"
	@gometalinter  --disable-all --enable=unconvert --deadline=45s  $(LINT_FOLDERS)
lint-staticcheck:
	@echo "lint-staticcheck"
	@gometalinter  --disable-all --enable=staticcheck --deadline=45s  $(LINT_FOLDERS)
lint-goconst:
	@echo "lint-goconst"
	@gometalinter  --disable-all --enable=goconst --deadline=45s  $(LINT_FOLDERS)
lint-gas:
	@echo "lint-gas"
	@gometalinter  --disable-all --enable=gas --deadline=45s  $(LINT_FOLDERS)
lint-varcheck:
	@echo "lint-varcheck"
	@gometalinter  --disable-all --enable=varcheck --deadline=45s  $(LINT_FOLDERS)
lint-structcheck:
	@echo "lint-structcheck"
	@gometalinter  --disable-all --enable=structcheck --deadline=45s  $(LINT_FOLDERS)
lint-gosimple:
	@echo "lint-gosimple"
	@gometalinter  --disable-all --enable=gosimple --deadline=45s  $(LINT_FOLDERS)
