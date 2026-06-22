package registry

import (
	"sort"
	"strings"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// searchResult holds a scored SPC match.
type searchResult struct {
	def   *spc.SPCDefinition
	score float64
}

// ResolveSearch performs keyword-based matching of a query against SPC definitions.
// Scoring: name match * 3.0 + description match * 2.0 + tag match * 1.0 + category match * 1.5
func ResolveSearch(query string, defs []*spc.SPCDefinition) []*spc.SPCDefinition {
	query = strings.ToLower(query)
	var results []searchResult

	for _, def := range defs {
		score := 0.0

		// Name keyword match
		if strings.Contains(strings.ToLower(def.Name), query) {
			score += 3.0
		}

		// Description keyword match
		if strings.Contains(strings.ToLower(def.Description), query) {
			score += 2.0
		}

		// Category match
		if strings.Contains(strings.ToLower(def.Category), query) {
			score += 1.5
		}

		// Tag match
		for _, tag := range def.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				score += 1.0
				break
			}
		}

		// Also check individual words in query against name/desc
		words := strings.Fields(query)
		for _, word := range words {
			if len(word) < 3 {
				continue
			}
			for _, tag := range def.Tags {
				if strings.Contains(strings.ToLower(tag), word) {
					score += 0.5
				}
			}
		}

		if score > 0 {
			results = append(results, searchResult{def: def, score: score})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	result := make([]*spc.SPCDefinition, len(results))
	for i, r := range results {
		result[i] = r.def
	}
	return result
}
