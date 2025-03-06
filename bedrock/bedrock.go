package bedrock

import (
	"log"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
	"gopkg.in/yaml.v3"
)

type Model struct {
	Name string `yaml:"modelName"`
	Id   string `yaml:"modelId"`
}

func LoadModels(wf *aw.Workflow, yamlPath string) {
	// load models from yaml
	yamlFile, err := os.ReadFile(yamlPath)
	if err != nil {
		log.Fatal(err)
	}
	models := []Model{}
	if err = yaml.Unmarshal(yamlFile, &models); err != nil {
		log.Fatal(err)
	}

	// add models to workflow
	for _, model := range models {
		name := model.Name
		id := model.Id
		match := name + " " + strings.ReplaceAll(strings.Join(strings.Split(id, "."), " "), "-", ".")
		wf.
			NewItem(name).
			Valid(true).
			Subtitle(id).
			Var("action", "run-script").
			Match(match).
			UID(id).
			Arg(id)
	}
}
