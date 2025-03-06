SHELL := /bin/bash

PLIST=info.plist
BIN=alfred-bedrock-model-id
ICON=./icon.png
ASSETS=./assets
DIST_FILE=bedrock-model-id.alfredworkflow

all: $(DIST_FILE)

$(BIN):
	go build -o $(BIN) ./main.go

$(DIST_FILE): $(BIN) $(PLIST) $(ICON) $(YAML) $(ASSETS)
	zip -r $(DIST_FILE) $(BIN) $(PLIST) $(ICON) $(YAML) $(ASSETS)