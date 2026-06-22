package registry

import (
	"fmt"
	"io/fs"
	"strings"
	"sync"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// RegistryImpl implements the SPC registry with 3-layer loading.
type RegistryImpl struct {
	mu    sync.RWMutex
	index map[string]*spc.SPCDefinition // name to definition
	all   []*spc.SPCDefinition          // all definitions in load order
}

// NewRegistry creates an empty registry.
func NewRegistry() *RegistryImpl {
	return &RegistryImpl{
		index: make(map[string]*spc.SPCDefinition),
	}
}

// LoadAll performs the full 3-layer discovery: builtin to user to project.
func (r *RegistryImpl) LoadAll(embedFS fs.FS, embedDir string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Layer 1: Built-in (embedded)
	builtins, _ := LoadFromEmbed(embedFS, embedDir)
	r.merge(builtins)

	// Layer 2: User custom (~/.openlibing/spc/)
	userSPCs, _ := LoadUserSPCs()
	r.merge(userSPCs)

	// Layer 3: Project local (./.openlibing/spc/)
	projectSPCs, _ := LoadProjectSPCs()
	r.merge(projectSPCs)

	return nil
}

// merge inserts or overrides definitions. Later layers take precedence.
func (r *RegistryImpl) merge(defs []*spc.SPCDefinition) {
	for _, def := range defs {
		// Check for override
		if existing, ok := r.index[def.Name]; ok {
			// Remove from all list
			for i, d := range r.all {
				if d == existing {
					r.all = append(r.all[:i], r.all[i+1:]...)
					break
				}
			}
		}
		r.index[def.Name] = def
		r.all = append(r.all, def)
	}
}

// Get retrieves an SPC by exact name.
func (r *RegistryImpl) Get(name string) (*spc.SPCDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.index[name]
	if !ok {
		return nil, fmt.Errorf("SPC '%s' not found — run 'openlibing list' to see available super powers", name)
	}
	return def, nil
}

// ListAll returns all registered SPC definitions.
func (r *RegistryImpl) ListAll() []*spc.SPCDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*spc.SPCDefinition, len(r.all))
	copy(result, r.all)
	return result
}

// ListByCategory returns SPCs filtered by category.
func (r *RegistryImpl) ListByCategory(cat string) []*spc.SPCDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*spc.SPCDefinition
	for _, def := range r.all {
		if strings.EqualFold(def.Category, cat) {
			result = append(result, def)
		}
	}
	return result
}

// Search performs keyword-based SPC search. MVP: keyword matching only.
func (r *RegistryImpl) Search(query string) []*spc.SPCDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return ResolveSearch(query, r.all)
}
