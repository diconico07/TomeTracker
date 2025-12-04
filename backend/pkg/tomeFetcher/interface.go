package tomefetcher

import (
	"context"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/diconico07/TomeTracker/pkg/tomeStore"
)

type Fetcher interface {
	IsURLSupported(url string) bool
	ProcessSeries(ctx context.Context, series *ent.Series) ([]*ent.Book, tomeStore.SeriesUpdate, error)
}
