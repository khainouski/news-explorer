package metrics

// buckets covers a wide range on purpose (10ms to 2min) since this is a generic HTTP-duration
// histogram, not tuned to any one endpoint's expected latency yet.
var buckets = []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10, 20, 30, 40, 50, 60, 90, 120}
