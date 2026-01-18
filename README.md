[![Go Reference](https://pkg.go.dev/badge/github.com/matovu-ronald/emojis.svg)](https://pkg.go.dev/github.com/matovu-ronald/emojis)

# emojis

A Go module that provides emoji search functionality. Search for emojis by description, tags, and categories.

## Overview

The `emojis` module is a Go package that offers efficient emoji searching capabilities. It allows you to search for emojis using inclusion and exclusion parameters based on emoji descriptions and tags.

## Installation

```bash
go get github.com/matovu-ronald/emojis
```

## Usage

### Basic Search

```go
package main

import (
	"fmt"
	"github.com/matovu-ronald/emojis/search"
)

func main() {
	params := search.Params{
		Include: []string{"fruit"},
	}

	results := search.ByDescription(params)
	for _, emoji := range results {
		fmt.Println(emoji.Label)
	}
}
```

### Like (Similar Emoji) Search

```go
package main

import (
	"fmt"
	"github.com/matovu-ronald/emojis/search"
)

func main() {
	results := search.Like(search.LikeParams{
		SeedEmoji: "🍇",
		Include:   []string{"fruit"},
		Exclude:   []string{"citrus"},
		Limit:     3,
	})

	for _, emoji := range results {
		fmt.Println(emoji.Emoji, emoji.Label)
	}
}
```

### Search Parameters

The `Params` struct supports:

- **Include**: Slice of strings to include in the search (searches emoji labels and tags)
- **Exclude**: Slice of strings to exclude from the search results

## API

### Package: `search`

#### `ByDescription(params Params) []Emoji`

Searches for emojis by description using the provided search parameters.

**Parameters:**

- `params` - A `Params` struct containing Include and Exclude slices

**Returns:**

- Slice of matching `Emoji` objects

#### `Params`

```go
type Params struct {
	Include []string // Labels/tags to include in search
	Exclude []string // Labels/tags to exclude from search
}

#### `Like(params LikeParams) []Emoji`

Searches for emojis related to a seed emoji or label using tag/label overlap with deterministic ordering and optional include/exclude filters.

```go
type LikeParams struct {
	SeedEmoji string
	SeedLabel string
	Include   []string
	Exclude   []string
	Limit     int
}
```

**Notes**

- Seed emoji is excluded from results.
- Include filters require at least one match; exclude filters remove any match.
- Ordering is deterministic: score desc, then label asc; limit defaults to 50 and clamps to the dataset size.
```

#### `Emoji`

Represents an emoji with its metadata (Label, Tags, etc.)

## License

See LICENSE file for details.
