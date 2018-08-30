LDFLAGS += -X "gitote.com/gitote/gitote/pkg/setting.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
LDFLAGS += -X "gitote.com/gitote/gitote/pkg/setting.BuildGitHash=$(shell git rev-parse HEAD)"

DATA_FILES := $(shell find conf | sed 's/ /\\ /g')
LESS_FILES := $(wildcard public/less/gitote.less public/less/_*.less)
GENERATED  := pkg/bindata/bindata.go public/css/gitote.css

OS := $(shell uname)

TAGS = ""
BUILD_FLAGS = "-v"

RELEASE_ROOT = "release"
RELEASE_GITOTE = "release/gitote"
NOW = $(shell date -u '+%Y%m%d%I%M%S')
GOVET = go tool vet -composites=false -methods=false -structtags=false
GOPATH = $(shell echo $${PWD%%src*})

.PHONY: build bindata clean

.IGNORE: public/css/gitote.css

all: build

check: test

dist: release

web: build
	./gitote web

govet:
	$(GOVET) gitote.go
	$(GOVET) models pkg routes

build: $(GENERATED)
	go install $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)'
	cp '$(GOPATH)/bin/gitote' .

build-dev: $(GENERATED) govet
	go install $(BUILD_FLAGS) -tags '$(TAGS)'
	cp '$(GOPATH)/bin/gitote' .

build-dev-race: $(GENERATED) govet
	go install $(BUILD_FLAGS) -race -tags '$(TAGS)'
	cp '$(GOPATH)/bin/gitote' .

bindata: pkg/bindata/bindata.go

pkg/bindata/bindata.go: $(DATA_FILES)
	./go-bindata -o=$@ -ignore="\\.DS_Store|README.md|TRANSLATORS|auth.d" -pkg=bindata conf/...

less: public/css/gitote.css

public/css/gitote.css: $(LESS_FILES)
	@type lessc >/dev/null 2>&1 && lessc $< >$@ || echo "lessc command not found, skipped."

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm

test:
	go test -cover -race ./...
