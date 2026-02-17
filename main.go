package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/toshikwa/alfred-bedrock-model-id/bedrock"
)

var (
	doCheck           bool
	doFetch           bool
	updateIcon        = &aw.Icon{Value: "assets/update-available.png"}
	repo              = "toshikwa/alfred-bedrock-model-id"
	checkForUpdateJob = "checkForUpdate"
	fetchModelsJob    = "fetchModels"
	maxCacheAge       = 24 * time.Hour
	wf                *aw.Workflow
)

func init() {
	flag.BoolVar(&doCheck, "check", false, "check for a new version")
	flag.BoolVar(&doFetch, "fetch", false, "fetch latest models from GitHub")
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

	// fetch models from GitHub
	if doFetch {
		wf.Configure(aw.TextErrors(true))
		log.Println("Fetching latest models from GitHub...")
		if err := bedrock.FetchModels(wf); err != nil {
			wf.FatalError(err)
		}
		return
	}

	// run update check in background
	if wf.UpdateCheckDue() && !wf.IsRunning(checkForUpdateJob) {
		log.Println("Running update check in background...")
		cmd := exec.Command(os.Args[0], "-check")
		if err := wf.RunInBackground(checkForUpdateJob, cmd); err != nil {
			log.Printf("Error starting update check: %s", err)
		}
	}

	// run model fetch in background if cache is stale
	if !wf.IsRunning(fetchModelsJob) && wf.Cache.Expired(bedrock.CacheFileName, maxCacheAge) {
		log.Println("Running model fetch in background...")
		cmd := exec.Command(os.Args[0], "-fetch")
		if err := wf.RunInBackground(fetchModelsJob, cmd); err != nil {
			log.Printf("Error starting model fetch: %s", err)
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
		bedrock.LoadModels(wf)
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
