package spc

// SPCDefinition is the complete parsed representation of a .spc.yaml file.
type SPCDefinition struct {
	Name        string      `yaml:"name"`
	Version     string      `yaml:"version"`
	Description string      `yaml:"description"`
	Type        string      `yaml:"type"` // query | action | workflow
	Category    string      `yaml:"category"`
	Tags        []string    `yaml:"tags"`
	Parameters  []Parameter `yaml:"parameters"`
	Source      Source      `yaml:"source"`
	Output      Output      `yaml:"output"`
	AI          *AIConfig   `yaml:"ai,omitempty"`
	Examples    []Example   `yaml:"examples,omitempty"`

	// Runtime metadata (not from YAML)
	Origin   string `yaml:"-"` // "builtin" | "user" | "project"
	FilePath string `yaml:"-"` // absolute path to the SPC file
}

// Parameter defines an input parameter for a Super Power.
type Parameter struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // string | integer | boolean
	Required    bool        `yaml:"required"`
	Default     interface{} `yaml:"default,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
	Validation  *Validation `yaml:"validation,omitempty"`
}

// Validation defines numeric range constraints for a parameter.
type Validation struct {
	Min *int `yaml:"min,omitempty"`
	Max *int `yaml:"max,omitempty"`
}

// Source defines the HTTP data source for a Super Power.
type Source struct {
	Method      string            `yaml:"method"` // GET | POST
	Endpoint    string            `yaml:"endpoint"`
	QueryParams map[string]string `yaml:"query_params,omitempty"`
	Headers     map[string]string `yaml:"headers,omitempty"`
	Body        string            `yaml:"body,omitempty"`
}

// Output defines how to format the response data.
type Output struct {
	Format string  `yaml:"format"` // table | json | yaml | raw
	Fields []Field `yaml:"fields,omitempty"`
}

// Field defines a single column/field in the output.
type Field struct {
	Name      string `yaml:"name"`
	Header    string `yaml:"header"`
	Path      string `yaml:"path"`                // gjson path
	Width     int    `yaml:"width,omitempty"`
	Transform string `yaml:"transform,omitempty"` // upper | lower | duration | truncate:N
}

// AIConfig provides hints for the AI chat mode.
type AIConfig struct {
	PromptHint      string   `yaml:"prompt_hint,omitempty"`
	NaturalLanguage []string `yaml:"natural_language,omitempty"`
	ResultHint      string   `yaml:"result_hint,omitempty"`
}

// Example provides usage examples for --help output.
type Example struct {
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
}

// Result holds the processed output of an SPC execution.
type Result struct {
	Format string                   // output format
	Rows   []map[string]interface{} // extracted rows
	Raw    []byte                   // raw HTTP response body
}
