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

	// query
	args := flag.Args()
	query := strings.TrimSpace(strings.Join(args[0:], " "))

	if query == "" || strings.HasPrefix("update", query) {
		// update
		if wf.UpdateAvailable() {
			wf.Configure(aw.SuppressUIDs(true))
			wf.NewItem("[alfred-bedrock-model-id] An update is available!!").
				Subtitle("Press Enter to install update").
				Valid(false).
				Autocomplete("workflow:update").
				Icon(updateIcon)
		} else {
			wf.NewItem("Search for Bedrock model ID...").
				Subtitle("e.g. `bm us sonnet`, `bm apac nova`")
		}
	} else {
		// filter models
		bedrock.LoadModels(wf, "./assets/models.yaml")
		wf.Filter(query)
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
