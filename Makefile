SHELL := /bin/bash

PLIST=info.plist
BIN=alfred-bedrock-model-id
ICON=./icon.png
ASSETS=./assets
DIST_FILE=bedrock-model-id.alfredworkflow
VERSION?=0.0.0

all: update-version $(DIST_FILE)

$(BIN):
	go build -o $(BIN) ./main.go

update-version:
	sed -i '' 's/__VERSION__/$(VERSION)/g' $(PLIST)

$(DIST_FILE): $(BIN) $(PLIST) $(ICON) $(YAML) $(ASSETS)
	zip -r $(DIST_FILE) $(BIN) $(PLIST) $(ICON) $(YAML) $(ASSETS)
