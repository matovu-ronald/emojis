# Quickstart: Like-emoji search

1) Add the module (if not already):
```bash
go get github.com/matovu-ronald/emojis
```

2) Import and call the like search:
```go
package main

import (
    "fmt"

    "github.com/matovu-ronald/emojis/search"
)

func main() {
    results := search.Like(search.LikeParams{
        SeedEmoji: "👍",
        Include:   []string{"hand"},
        Exclude:   []string{"skin-tone"},
        Limit:     5,
    })

    for _, e := range results {
        fmt.Printf("%s %v\n", e.Label, e.Tags)
    }
}
```

3) Run tests to validate behavior:
```bash
go test ./...
```

4) Expected behavior:
- Deterministic ordering for identical inputs.
- Exclude tags remove any matching candidates; include tags require at least one match when provided.
- Seed emoji is not echoed in the result set.
