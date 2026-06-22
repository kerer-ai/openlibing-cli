package engine

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/openlibing/openlibing-cli/internal/api"
	"github.com/openlibing/openlibing-cli/pkg/spc"
	"github.com/tidwall/gjson"
)

// SPCLookup is the minimal interface Engine needs from Registry.
type SPCLookup interface {
	Get(name string) (*spc.SPCDefinition, error)
}

// Engine orchestrates the 7-step SPC execution pipeline.
type Engine struct {
	registry SPCLookup
	client   *api.Client
}

// NewEngine creates a new Engine.
func NewEngine(registry SPCLookup, client *api.Client) *Engine {
	return &Engine{
		registry: registry,
		client:   client,
	}
}

// Execute runs the full 7-step pipeline: Parse → Validate → Resolve → Call → Extract → Transform → Output.
func (e *Engine) Execute(name string, params map[string]interface{}) (*spc.Result, error) {
	// Step 1: Parse — lookup SPC from registry
	def, err := e.registry.Get(name)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	// Step 2: Validate — check parameters
	if err := Validate(def, params); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	// Apply defaults for unset parameters
	params = applyDefaults(def, params)

	// Step 3: Resolve — render templates
	resolved, err := Resolve(&def.Source, params)
	if err != nil {
		return nil, fmt.Errorf("resolve: %w", err)
	}

	// Step 4: Call — execute HTTP request
	respBody, err := e.call(resolved)
	if err != nil {
		return nil, fmt.Errorf("call: %w", err)
	}

	// Step 5: Extract — parse JSON and apply gjson paths
	rows, err := extract(def, respBody)
	if err != nil {
		// Degrade: return raw JSON
		return &spc.Result{
			Format: "raw",
			Raw:    respBody,
		}, nil
	}

	// Steps 6 & 7: Transform (applied during format stage in CLI) and Output
	return &spc.Result{
		Format: def.Output.Format,
		Rows:   rows,
		Raw:    respBody,
	}, nil
}

func (e *Engine) call(resolved *ResolvedRequest) ([]byte, error) {
	resp, err := e.client.Do(
		resolved.Method,
		resolved.Endpoint,
		resolved.QueryParams,
		resolved.Headers,
		nil, // body handled via query params for now
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func extract(def *spc.SPCDefinition, rawJSON []byte) ([]map[string]interface{}, error) {
	jsonStr := string(rawJSON)

	// Determine if response wraps data in a "data" array
	dataResult := gjson.Get(jsonStr, "data")
	var arr []gjson.Result
	if dataResult.IsArray() {
		arr = dataResult.Array()
	} else {
		// Try root-level array
		rootResult := gjson.Parse(jsonStr)
		if rootResult.IsArray() {
			arr = rootResult.Array()
		} else {
			// Single object — wrap in array
			arr = []gjson.Result{rootResult}
		}
	}

	var rows []map[string]interface{}
	for _, item := range arr {
		row := make(map[string]interface{})
		for _, field := range def.Output.Fields {
			val := item.Get(field.Path)
			if val.Exists() {
				row[field.Name] = val.Value()
			}
		}
		// If no fields defined, include all keys
		if len(def.Output.Fields) == 0 {
			item.ForEach(func(key, value gjson.Result) bool {
				row[key.String()] = value.Value()
				return true
			})
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func applyDefaults(def *spc.SPCDefinition, params map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range params {
		result[k] = v
	}
	for _, p := range def.Parameters {
		if _, exists := result[p.Name]; !exists && p.Default != nil {
			result[p.Name] = p.Default
		}
	}
	return result
}

// Ensure json import is used
var _ = json.Valid
