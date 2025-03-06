package main

import (
	"flag"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/toshikwa/alfred-bedrock-model-id/bedrock"
)

var (
	doCheck           bool
	region            string
	query             string
	updateIcon        = &aw.Icon{Value: "assets/update-available.png"}
	repo              = "toshikwa/alfred-bedrock-model-id"
	checkForUpdateJob = "checkForUpdate"
	wf                *aw.Workflow
)

func init() {
	flag.BoolVar(&doCheck, "check", false, "check for a new version")
	wf = aw.New(update.GitHub(repo))
	bedrock.LoadModels(wf, "./assets/models.yaml")
}

func run() {
	wf.Args()
	flag.Parse()

	// check for update
	if doCheck {
		wf.Configure(aw.TextErrors(true))
		log.Println("Checking for updates...")
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
		return
	}
	defer finalize()

	// filter models by query
	args := flag.Args()
	if len(args) > 0 {
		query = strings.TrimSpace(strings.Join(args[1:], " "))
		if query != "" {
			wf.Filter(query)
		}
	}
}

func finalize() {
	if r := recover(); r != nil {
		panic(r)
	}
	if wf.IsEmpty() {
		wf.NewItem("No matching Bedrock model found.").
			Subtitle("Try another query (e.g. `bm nova`, `bm us sonnet`)").
			Icon(aw.IconNote)
	}
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
