package engine

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// ResolvedRequest holds all template-resolved fields from a Source definition.
// Templates like {{.project_id}} have been rendered against user-supplied params.
type ResolvedRequest struct {
	Method      string
	Endpoint    string            // e.g. "gateway/openlibing-cicd/..."
	QueryParams map[string]string // resolved key=value pairs
	Headers     map[string]string
	Body        string
}

// Resolve renders all Go templates in a Source definition against the provided params.
func Resolve(source *spc.Source, params map[string]interface{}) (*ResolvedRequest, error) {
	ctx := make(map[string]interface{})
	for k, v := range params {
		ctx[k] = fmt.Sprintf("%v", v)
	}

	rr := &ResolvedRequest{
		Method:      source.Method,
		QueryParams: make(map[string]string),
		Headers:     make(map[string]string),
	}

	// Resolve endpoint
	var err error
	rr.Endpoint, err = render("endpoint", source.Endpoint, ctx)
	if err != nil {
		return nil, err
	}

	// Resolve query params
	for k, v := range source.QueryParams {
		rr.QueryParams[k], err = render("query:"+k, v, ctx)
		if err != nil {
			return nil, err
		}
	}

	// Resolve headers
	for k, v := range source.Headers {
		rr.Headers[k], err = render("header:"+k, v, ctx)
		if err != nil {
			return nil, err
		}
	}

	// Resolve body
	if source.Body != "" {
		rr.Body, err = render("body", source.Body, ctx)
		if err != nil {
			return nil, err
		}
	}

	return rr, nil
}

func render(name, tmpl string, ctx map[string]interface{}) (string, error) {
	t, err := template.New(name).Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse [%s]: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("template execute [%s]: %w", name, err)
	}
	return buf.String(), nil
}
