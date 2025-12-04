package tomefetcher

import (
	"fmt"

	"resty.dev/v3"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/diconico07/TomeTracker/pkg/tomeStore"
)

type TomeFetcher struct {
	Client *resty.Client
	Next   Fetcher
}

func New(base_url string, next Fetcher) *TomeFetcher {
	return &TomeFetcher{
		Client: resty.New().SetBaseURL(base_url),
		Next:   next,
	}
}

func (t *TomeFetcher) ListSeries() ([]*ent.Series, error) {
	var series []*ent.Series

	resp, err := t.Client.R().
		SetResult(&series).
		Get("/api/series") // Replace with your actual API endpoint

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return series, nil
}

func (t *TomeFetcher) GetSeries(id int) (*ent.Series, error) {
	var series ent.Series

	resp, err := t.Client.R().
		SetResult(&series).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		Get("/api/series/{id}")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &series, nil
}

func (t *TomeFetcher) UpsertBook(book *ent.Book) (*ent.Book, error) {
	var createdBook ent.Book

	resp, err := t.Client.R().
		SetBody(book).
		SetResult(&createdBook).
		SetPathParam("isbn", book.ID).
		Post("/api/books/{isbn}")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &createdBook, nil
}

func (t *TomeFetcher) UpdateSeries(id int, series *tomeStore.SeriesUpdate) (*ent.Series, error) {
	var updatedSeries ent.Series

	resp, err := t.Client.R().
		SetBody(series).
		SetResult(&updatedSeries).
		SetPathParam("id", fmt.Sprintf("%d", id)).
		Patch("/api/series/{id}")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return &updatedSeries, nil
}
