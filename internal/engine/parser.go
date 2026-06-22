package engine

import (
	"fmt"
	"io"
	"os"

	"github.com/openlibing/openlibing-cli/pkg/spc"
	"gopkg.in/yaml.v3"
)

// ParseSPC reads a YAML stream and returns a parsed SPCDefinition.
func ParseSPC(r io.Reader) (*spc.SPCDefinition, error) {
	var def spc.SPCDefinition
	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)
	if err := decoder.Decode(&def); err != nil {
		return nil, fmt.Errorf("parse SPC: %w", err)
	}
	return &def, nil
}

// ParseSPCFile reads a .spc.yaml file and returns a parsed SPCDefinition.
func ParseSPCFile(path string) (*spc.SPCDefinition, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open SPC file %s: %w", path, err)
	}
	defer f.Close()

	def, err := ParseSPC(f)
	if err != nil {
		return nil, fmt.Errorf("file %s: %w", path, err)
	}
	def.FilePath = path
	return def, nil
}
