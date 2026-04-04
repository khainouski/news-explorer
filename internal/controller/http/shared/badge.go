package shared

// Badge is a source's small colored square logo (e.g. "Go" on blue, "TC" on green). Shared
// between the article and source packages - both render source badges (feed rows, sidebar,
// sources table).
type Badge struct {
	Text  string
	Color string // Tailwind background color class, e.g. "bg-blue-500"
}
