package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/openlibing/openlibing-cli/pkg/spc"
	"gopkg.in/yaml.v3"
)

// FormatResult formats a Result according to its Format field.
func FormatResult(result *spc.Result) (string, error) {
	switch result.Format {
	case "json":
		return FormatJSON(result.Rows)
	case "yaml":
		return FormatYAML(result.Rows)
	case "table":
		return FormatTable(result.Rows, nil), nil
	case "raw":
		return FormatRaw(result.Raw), nil
	default:
		return FormatTable(result.Rows, nil), nil
	}
}

// FormatTable renders rows as an aligned terminal table.
func FormatTable(rows []map[string]interface{}, fields []spc.Field) string {
	if len(rows) == 0 {
		return "(no results)"
	}

	// If no fields specified, infer from first row keys
	if len(fields) == 0 {
		first := rows[0]
		for k := range first {
			fields = append(fields, spc.Field{Name: k, Header: k, Path: k})
		}
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Headers
	headers := make([]string, len(fields))
	for i, f := range fields {
		headers[i] = f.Header
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Separator
	seps := make([]string, len(fields))
	for i := range fields {
		seps[i] = strings.Repeat("─", len(headers[i])+4)
	}
	fmt.Fprintln(w, strings.Join(seps, "\t"))

	// Rows
	for _, row := range rows {
		cols := make([]string, len(fields))
		for i, f := range fields {
			val := row[f.Name]
			col := formatValue(val, f.Transform)
			cols[i] = col
		}
		fmt.Fprintln(w, strings.Join(cols, "\t"))
	}

	w.Flush()
	return buf.String()
}

// FormatJSON renders rows as indented JSON.
func FormatJSON(rows []map[string]interface{}) (string, error) {
	if rows == nil {
		rows = []map[string]interface{}{}
	}
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(data), nil
}

// FormatYAML renders rows as YAML.
func FormatYAML(rows []map[string]interface{}) (string, error) {
	if rows == nil {
		rows = []map[string]interface{}{}
	}
	data, err := yaml.Marshal(rows)
	if err != nil {
		return "", fmt.Errorf("yaml marshal: %w", err)
	}
	return string(data), nil
}

// FormatRaw returns raw bytes as a string.
func FormatRaw(raw []byte) string {
	return string(raw)
}

func formatValue(val interface{}, transform string) string {
	if val == nil {
		return "-"
	}

	str := fmt.Sprintf("%v", val)

	switch transform {
	case "upper":
		return strings.ToUpper(str)
	case "lower":
		return strings.ToLower(str)
	case "duration":
		return TransformDuration(val)
	case "":
		return str
	default:
		if strings.HasPrefix(transform, "truncate:") {
			var n int
			fmt.Sscanf(transform, "truncate:%d", &n)
			if len(str) > n {
				return str[:n] + "…"
			}
		}
		return str
	}
}

// TransformDuration converts milliseconds to a human-readable duration string.
func TransformDuration(val interface{}) string {
	var ms int64
	switch v := val.(type) {
	case float64:
		ms = int64(v)
	case int64:
		ms = v
	case int:
		ms = int64(v)
	default:
		return fmt.Sprintf("%v", val)
	}

	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}

	seconds := ms / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, remainingMinutes)
}
