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
}

func run() {
	wf.Args()
	flag.Parse()
	validRegions := []string{
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-north-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}

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

	args := flag.Args()
	if len(args) > 0 {
		// check if the first argument is a valid region
		possibleRegion := args[0]
		isRegion := false
		for _, r := range validRegions {
			if possibleRegion == r {
				region = possibleRegion
				isRegion = true
				break
			}
		}

		if isRegion {
			// if it was a valid region, get the query
			if len(args) > 1 {
				query = strings.Join(args[1:], " ")
			}
		} else {
			// if it was not a valid region, show region suggestions
			if possibleRegion != "" {
				for _, r := range validRegions {
					if strings.HasPrefix(r, possibleRegion) {
						wf.NewItem(r).
							Subtitle("Use '" + r + "' region").
							Arg(r).
							Autocomplete(r + " ").
							Valid(false)
					}
				}
			}
			return
		}

		// load models
		bedrock.LoadModels(wf, "./assets/fm-"+region+".yaml", false)
		bedrock.LoadModels(wf, "./assets/cri-"+region+".yaml", true)

		// filter results
		if query != "" {
			wf.Filter(strings.ToLower(query))
		}
	}

}

func finalize() {
	if r := recover(); r != nil {
		panic(r)
	}
	if wf.IsEmpty() {
		wf.NewItem("No matching Bedrock model found.").
			Subtitle("Try another query (e.g. `bm us-west-2 nova`, `bm us-east-1 cri sonnet`)").
			Icon(aw.IconNote)
	}
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
