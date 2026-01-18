package search

import (
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestByDesc(t *testing.T) {
	tests := map[string]struct {
		params Params
		want   []string
	}{
		"fruit": {
			Params{Include: []string{"fruit"}},
			[]string{"🍇", "🍈", "🍉", "🍊", "🍋"},
		},
		"cat": {
			Params{Include: []string{"cat"}},
			[]string{"🐱"},
		},
		"animal faces": {
			Params{Include: []string{"face"}, Exclude: []string{"smile", "laugh", "grin", "upside-down"}},
			[]string{"🐵", "🐶", "🐱", "🐯", "🦊"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ByDescription(test.params)
			if len(result) != len(test.want) {
				t.Fatalf("Search ByDescription: want %d emojis, got %d", len(test.want), len(result))
			}

			for _, emoji := range result {
				if !slices.Contains(test.want, emoji.Emoji) {
					t.Fatalf("Search ByDescription: expected %s to be in %#v", emoji.Emoji, test.want)
				}
			}
		})
	}
}

func TestNormalizeTerms(t *testing.T) {
	got := normalizeTerms([]string{" Face ", "face", "HAND", ""})
	want := []string{"face", "hand"}

	if !slices.Equal(got, want) {
		t.Fatalf("normalizeTerms mismatch: got %v, want %v", got, want)
	}
}

func TestSplitTerms(t *testing.T) {
	got := splitTerms("  Grinning Face with  Eyes  ")
	want := []string{"grinning", "face", "with", "eyes"}

	if !slices.Equal(got, want) {
		t.Fatalf("splitTerms mismatch: got %v, want %v", got, want)
	}
}

func TestOverlapScore(t *testing.T) {
	seed := map[string]struct{}{"face": {}, "grin": {}, "open": {}}
	candidate := map[string]struct{}{"face": {}, "smile": {}, "open": {}, "laugh": {}}

	got := overlapScore(seed, candidate)
	if got <= 0 {
		t.Fatalf("overlapScore expected positive, got %f", got)
	}

	if got >= 1 {
		t.Fatalf("overlapScore expected < 1, got %f", got)
	}
}

func TestLikeSimilarityOrdering(t *testing.T) {
	params := LikeParams{SeedEmoji: "😀"}
	result := Like(params)

	if len(result) < 3 {
		t.Fatalf("expected at least 3 related emojis (SC-001), got %d", len(result))
	}

	wantTop := []string{"🦊", "😅", "😆"}
	for i, emoji := range wantTop {
		if result[i].Emoji != emoji {
			t.Fatalf("unexpected order at %d: got %s want %s", i, result[i].Emoji, emoji)
		}
	}

	repeat := Like(params)
	if !equalEmojiLists(result, repeat) {
		t.Fatalf("expected deterministic results across runs (SC-002)")
	}

	for i := 0; i < 20; i++ {
		run := Like(params)
		if !equalEmojiLists(result, run) {
			t.Fatalf("determinism failed on run %d (SC-002)", i+1)
		}
	}
}

func TestLikeUnknownAndEmptySeed(t *testing.T) {
	if got := Like(LikeParams{SeedEmoji: ""}); len(got) != 0 {
		t.Fatalf("expected empty result for empty seed, got %v", got)
	}

	if got := Like(LikeParams{SeedLabel: "nonexistent"}); len(got) != 0 {
		t.Fatalf("expected empty result for unknown seed label, got %v", got)
	}
}

func TestLikeIncludeExcludeFilters(t *testing.T) {
	tests := []struct {
		name   string
		params LikeParams
		forbid []string
	}{
		{
			name:   "face exclude laugh",
			params: LikeParams{SeedEmoji: "😀", Include: []string{"face"}, Exclude: []string{"laugh"}},
			forbid: []string{"laugh"},
		},
		{
			name:   "face exclude grin",
			params: LikeParams{SeedEmoji: "😀", Include: []string{"face"}, Exclude: []string{"grin"}},
			forbid: []string{"grin"},
		},
		{
			name:   "face exclude upside-down",
			params: LikeParams{SeedEmoji: "😀", Include: []string{"face"}, Exclude: []string{"upside-down"}},
			forbid: []string{"upside-down"},
		},
		{
			name:   "fruit exclude citrus",
			params: LikeParams{SeedEmoji: "🍇", Include: []string{"fruit"}, Exclude: []string{"citrus"}},
			forbid: []string{"citrus"},
		},
		{
			name:   "fruit exclude grape",
			params: LikeParams{SeedEmoji: "🍇", Include: []string{"fruit"}, Exclude: []string{"grape"}},
			forbid: []string{"grape"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Like(tt.params)
			if len(result) == 0 {
				t.Fatalf("expected filtered results, got none")
			}

			for _, emo := range result {
				for _, term := range tt.forbid {
					if slices.Contains(emo.Tags, term) || strings.Contains(emo.Label, term) {
						t.Fatalf("exclude tag matched %q in %+v", term, emo)
					}
				}
			}
		})
	}
}

func TestLikeLimitClampingAndStability(t *testing.T) {
	params := LikeParams{SeedEmoji: "😀", Limit: 2}
	first := Like(params)
	if len(first) != 2 {
		t.Fatalf("expected limit of 2 results, got %d", len(first))
	}

	params.Limit = -1
	if got := Like(params); len(got) != 0 {
		t.Fatalf("expected empty result for negative limit, got %d", len(got))
	}

	params.Limit = 100
	params.SeedEmoji = "🍇"
	res := Like(params)
	if len(res) == 0 {
		t.Fatalf("expected results for grape seed, got none")
	}

	second := Like(LikeParams{SeedEmoji: "🍇", Limit: 100})
	if !equalEmojiLists(res, second) {
		t.Fatalf("expected stable ordering across runs for grape seed")
	}
}

func TestLikePerformance(t *testing.T) {
	params := LikeParams{SeedEmoji: "😀", Limit: 50}
	start := time.Now()
	for i := 0; i < 5000; i++ {
		_ = Like(params)
	}

	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("expected like search to complete within 500ms for 5000 runs, got %v", elapsed)
	}
}

func equalEmojiLists(a, b []Emoji) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Emoji != b[i].Emoji {
			return false
		}
	}
	return true
}

func ExampleLike() {
	results := Like(LikeParams{SeedEmoji: "🍇", Include: []string{"fruit"}, Exclude: []string{"citrus"}, Limit: 3})
	for _, e := range results {
		fmt.Println(e.Emoji, e.Label)
	}
	// Output:
	// 🍈 melon
	// 🍉 watermelon
	// 🍊 tangerine
}
