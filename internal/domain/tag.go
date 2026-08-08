package domain

// Tag is a topic sources are categorized under - a table row, not a fixed enum. A Source has
// exactly one; Article doesn't reference Tag at all - it's reached transitively via
// Article.SourceID.
type Tag struct {
	ID   string
	Name string
}
