package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
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
	defer finalize()

	// check for update
	if doCheck {
		wf.Configure(aw.TextErrors(true))
		log.Println("Checking for updates...")
		if err := wf.CheckForUpdate(); err != nil {
			wf.FatalError(err)
		}
		return
	}

	// execute command to check
	if wf.UpdateCheckDue() && !wf.IsRunning(checkForUpdateJob) {
		log.Println("Running update check in background...")
		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(checkForUpdateJob, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	// filter models by query
	args := flag.Args()
	if len(args) > 0 {
		query = strings.TrimSpace(strings.Join(args[0:], " "))
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
