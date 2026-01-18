package search

import (
	"slices"
	"strings"
)

const defaultLikeLimit = 50

// Params used to specify search parameters
type Params struct {
	Include []string // Include slice of strings to include in search
	Exclude []string // Exclude slice of strings to exclude in search
}

// LikeParams defines input for Like search.
type LikeParams struct {
	SeedEmoji string
	SeedLabel string
	Include   []string
	Exclude   []string
	Limit     int
}

// ByDescription searches emojis using the params
func ByDescription(params Params) (result []Emoji) {

	for _, emo := range emojis {
		if shouldExclude(emo, params.Exclude) {
			continue
		}

		for _, include := range params.Include {
			include = strings.ToLower(include)
			if strings.Contains(emo.Label, include) || slices.Contains(emo.Tags, include) {
				result = append(result, emo)
			}
		}
	}

	return
}

// ByTags searches emojis that match all the provided tags
func ByTags(tags ...string) (result []Emoji) {
tagLoop:
	for _, emo := range emojis {
		for _, tag := range tags {
			tag = strings.ToLower(tag)
			if !slices.Contains(emo.Tags, tag) {
				continue tagLoop
			}
		}
		result = append(result, emo)
	}

	return
}

// Like returns emojis related to a seed emoji or label, ordered deterministically.
func Like(params LikeParams) []Emoji {
	normalizedInclude := normalizeTerms(params.Include)
	normalizedExclude := normalizeTerms(params.Exclude)
	seedEmoji := strings.TrimSpace(params.SeedEmoji)
	seedLabel := strings.TrimSpace(params.SeedLabel)

	if seedEmoji == "" && seedLabel == "" {
		return nil
	}

	limit := params.Limit
	if limit == 0 {
		limit = defaultLikeLimit
	}
	if limit < 0 {
		limit = 0
	}

	seed, seedFound := findSeed(seedEmoji)
	seedTerms := seedTermsFromInput(seed, seedFound, seedLabel)
	if len(seedTerms) == 0 || limit == 0 {
		return nil
	}

	type scored struct {
		emoji Emoji
		score float64
	}

	var scoredResults []scored

	for _, candidate := range emojis {
		if seedFound && candidate.Emoji == seed.Emoji {
			continue
		}

		candidateTerms := makeTermSet(candidate.Label, candidate.Tags)
		if shouldFilter(candidate, candidateTerms, normalizedInclude, normalizedExclude) {
			continue
		}

		score := overlapScore(seedTerms, candidateTerms)
		if score == 0 {
			continue
		}

		scoredResults = append(scoredResults, scored{emoji: candidate, score: score})
	}

	if len(scoredResults) == 0 {
		return nil
	}

	slices.SortFunc(scoredResults, func(a, b scored) int {
		if a.score == b.score {
			return strings.Compare(a.emoji.Label, b.emoji.Label)
		}
		if a.score > b.score {
			return -1
		}
		return 1
	})

	if limit > len(scoredResults) {
		limit = len(scoredResults)
	}

	result := make([]Emoji, 0, limit)
	for i := 0; i < limit; i++ {
		result = append(result, scoredResults[i].emoji)
	}

	return result
}

func shouldFilter(emo Emoji, termSet map[string]struct{}, includes, excludes []string) bool {
	if matchesTerms(emo, termSet, excludes) {
		return true
	}

	if len(includes) == 0 {
		return false
	}

	return !matchesTerms(emo, termSet, includes)
}

func matchesTerms(emo Emoji, termSet map[string]struct{}, terms []string) bool {
	for _, term := range terms {
		if _, ok := termSet[term]; ok {
			return true
		}
		if strings.Contains(emo.Label, term) {
			return true
		}
	}

	return false
}

func seedTermsFromInput(seed Emoji, found bool, seedLabel string) map[string]struct{} {
	if found {
		return makeTermSet(seed.Label, seed.Tags)
	}

	labelTerms := splitTerms(seedLabel)
	return makeTermSet(strings.Join(labelTerms, " "), nil)
}

func findSeed(seedEmoji string) (Emoji, bool) {
	for _, emo := range emojis {
		if emo.Emoji == seedEmoji {
			return emo, true
		}
	}
	return Emoji{}, false
}

func makeTermSet(label string, tags []string) map[string]struct{} {
	set := map[string]struct{}{}

	for _, term := range splitTerms(label) {
		set[term] = struct{}{}
	}

	for _, tag := range tags {
		normalized := strings.TrimSpace(strings.ToLower(tag))
		if normalized == "" {
			continue
		}
		set[normalized] = struct{}{}
	}

	return set
}

func splitTerms(value string) []string {
	normalized := strings.TrimSpace(strings.ToLower(value))
	if normalized == "" {
		return nil
	}
	return strings.Fields(normalized)
}

func normalizeTerms(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))

	for _, v := range values {
		normalized := strings.TrimSpace(strings.ToLower(v))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}

	return result
}

func overlapScore(seed, candidate map[string]struct{}) float64 {
	if len(seed) == 0 || len(candidate) == 0 {
		return 0
	}
	shared := 0
	for term := range seed {
		if _, ok := candidate[term]; ok {
			shared++
		}
	}
	union := len(seed) + len(candidate) - shared
	if union == 0 {
		return 0
	}

	return float64(shared) / float64(union)
}

// shouldExclude checks emoji tags and labels for exclusions
func shouldExclude(emo Emoji, excludes []string) bool {
	for _, exclude := range excludes {
		exclude = strings.ToLower(exclude)
		if strings.Contains(emo.Label, exclude) || slices.Contains(emo.Tags, exclude) {
			return true
		}
	}

	return false
}
