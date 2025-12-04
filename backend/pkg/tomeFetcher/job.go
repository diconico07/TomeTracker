package tomefetcher

import (
	"context"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/rs/zerolog/log"
)

func (t *TomeFetcher) ProcessAll(ctx context.Context) error {
	series, err := t.ListSeries()
	if err != nil {
		return err
	}
	booksToUpsert := make([]*ent.Book, 0)
	for _, s := range series {
		if t.Next.IsURLSupported(s.URL) {
			log.Debug().Str("URL", s.URL).Int("ID", s.ID).Msg("Process Series")
			detailedSeries, err := t.GetSeries(s.ID)
			if err != nil {
				return err
			}
			books, updatedSeries, err := t.Next.ProcessSeries(ctx, detailedSeries)
			if err != nil {
				return err
			}
			_, err = t.UpdateSeries(s.ID, &updatedSeries)
			if err != nil {
				return err
			}
			booksToUpsert = append(booksToUpsert, books...)
		} else {
			log.Debug().Str("URL", s.URL).Int("ID", s.ID).Msg("Skip Series")
		}
	}

	for _, book := range booksToUpsert {
		_, err := t.UpsertBook(book)
		if err != nil {
			return err
		}
	}

	return nil
}
