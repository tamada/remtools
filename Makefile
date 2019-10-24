GO=go
NAME := remtools
VERSION := 5.0.0
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)'
	-X 'main.revision=$(REVISION)'

all: test build

update_version:
	@for i in README.md docs/content/_index.md; do\
	    sed -e 's!Version-[0-9.]*-yellowgreen!Version-${VERSION}-yellowgreen!g' -e 's!tag/v[0-9.]*!tag/v${VERSION}!g' $$i > a ; mv a $$i; \
	done

	@sed 's/const VERSION = .*/const VERSION = "${VERSION}"/g' common/config.go > a
	@mv a common/config.go
	@echo "Replace version to \"${VERSION}\""

setup: deps update_version
	git submodule update --init

test: setup format lint
	$(GO) test -covermode=count -coverprofile=coverage.out $$(go list ./... | grep -v vendor)

build: setup
	$(GO) build -o $(NAME) -v

lint: setup
	$(GO) vet $$(go list ./... | grep -v vendor)
	for pkg in $$(go list ./... | grep -v vendor); do \
		golint -set_exit_status $$pkg || exit $$?; \
	done

install: test build
	$(GO) install $(LDFLAGS)

clean:
	$(GO) clean
	rm -rf $(NAME)
