package bedrock

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	aw "github.com/deanishe/awgo"
	"gopkg.in/yaml.v3"
)

const CacheFileName = "models.yaml"

var (
	modelsURL       = "https://raw.githubusercontent.com/toshikwa/alfred-bedrock-model-id/main/assets/models.yaml"
	bundledYAMLPath = "./assets/models.yaml"
)

type Model struct {
	Name string `yaml:"modelName"`
	Id   string `yaml:"modelId"`
}

func FetchModels(wf *aw.Workflow) error {
	resp, err := http.Get(modelsURL)
	if err != nil {
		return fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Validate YAML before saving
	var models []Model
	if err := yaml.Unmarshal(data, &models); err != nil {
		return fmt.Errorf("invalid YAML from remote: %w", err)
	}
	if len(models) == 0 {
		return fmt.Errorf("remote YAML contains no models")
	}

	cachePath := filepath.Join(wf.CacheDir(), CacheFileName)
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	log.Printf("Fetched %d models from GitHub and cached to %s", len(models), cachePath)
	return nil
}

func LoadModels(wf *aw.Workflow) {
	var yamlFile []byte
	var err error

	// Try cache first
	cachePath := filepath.Join(wf.CacheDir(), CacheFileName)
	yamlFile, err = os.ReadFile(cachePath)

	// Fall back to bundled YAML
	if err != nil {
		log.Printf("Cache not found, falling back to bundled YAML: %s", err)
		yamlFile, err = os.ReadFile(bundledYAMLPath)
		if err != nil {
			log.Fatal(err)
		}
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
