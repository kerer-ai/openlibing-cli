package registry

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/config"
	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// LoadFromEmbed reads all .spc.yaml files from an embedded filesystem.
func LoadFromEmbed(embedFS fs.FS, dir string) ([]*spc.SPCDefinition, error) {
	entries, err := fs.ReadDir(embedFS, dir)
	if err != nil {
		return nil, err
	}

	var defs []*spc.SPCDefinition
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".spc.yaml") {
			continue
		}

		f, err := embedFS.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		defer f.Close()

		def, err := engine.ParseSPC(f)
		if err != nil {
			continue // skip malformed SPCs
		}
		def.Origin = "builtin"
		defs = append(defs, def)
	}

	return defs, nil
}

// LoadFromDir scans a directory for .spc.yaml files and parses them.
func LoadFromDir(dir, origin string) ([]*spc.SPCDefinition, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil // directory doesn't exist, no SPCs to load
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var defs []*spc.SPCDefinition
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".spc.yaml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		def, err := engine.ParseSPCFile(path)
		if err != nil {
			continue // skip malformed SPCs
		}
		def.Origin = origin
		defs = append(defs, def)
	}

	return defs, nil
}

// LoadUserSPCs loads SPCs from ~/.openlibing/spc/
func LoadUserSPCs() ([]*spc.SPCDefinition, error) {
	dir, err := config.SPCUserDir()
	if err != nil {
		return nil, nil
	}
	return LoadFromDir(dir, "user")
}

// LoadProjectSPCs loads SPCs from ./.openlibing/spc/
func LoadProjectSPCs() ([]*spc.SPCDefinition, error) {
	return LoadFromDir(".openlibing/spc", "project")
}
