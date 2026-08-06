package article

import (
	"context"
	"fmt"
	"strings"

	"github.com/khainouski/news-explorer/internal/domain"
	"github.com/khainouski/news-explorer/pkg/otel/tracer"
)

const insertBatchColumns = 7 // id, source_id, external_id, title, summary, url, published_at

// InsertBatch inserts every article in one round trip, skipping ones that already exist for their
// source (ON CONFLICT DO NOTHING). Returns how many rows were actually new.
func (r *Repo) InsertBatch(ctx context.Context, articles []domain.Article) (int, error) {
	ctx, span := tracer.Start(ctx, "adapter postgres article InsertBatch")
	defer span.End()

	if len(articles) == 0 {
		return 0, nil
	}

	placeholders := make([]string, len(articles))
	args := make([]any, 0, len(articles)*insertBatchColumns)

	for i, a := range articles {
		base := i*insertBatchColumns + 1
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			base, base+1, base+2, base+3, base+4, base+5, base+6)
		args = append(args, a.ID, a.SourceID, a.ExternalID, a.Title, a.Summary, a.URL, a.PublishedAt)
	}

	q := `
		INSERT INTO articles (id, source_id, external_id, title, summary, url, published_at)
		VALUES ` + strings.Join(placeholders, ", ") + `
		ON CONFLICT (source_id, external_id) DO NOTHING`

	cmdTag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("pool.Exec: %w", err)
	}

	return int(cmdTag.RowsAffected()), nil
}
