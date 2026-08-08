package shared

// Badge is a source's small colored square logo (e.g. "Go" on blue, "TC" on green).
type Badge struct {
	Text  string
	Color string // Tailwind background color class, e.g. "bg-blue-500"
}
