package tomeStore

import (
	"context"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/diconico07/TomeTracker/pkg/ent/book"
	"github.com/diconico07/TomeTracker/pkg/ent/series"
)

// SeriesService describes the Greeter service.
type SeriesService interface {
	ListSeries(ctx context.Context) ([]*ent.Series, error)
	GetSeries(ctx context.Context, id int) (*ent.Series, error)
	AddSeries(ctx context.Context, name string, url string) (*ent.Series, error)
	DeleteSeries(ctx context.Context, id int) error
	UpdateSeries(ctx context.Context, id int, series SeriesUpdate) (*ent.Series, error)
	CountBooksInSeries(ctx context.Context, id int) (int, int, error)
}

// NewService creates a new Greeter service.
func NewService(client *ent.Client) SeriesService {
	return &tomeStoreService{
		Client: client,
	}
}

type tomeStoreService struct {
	Client *ent.Client
}

func (s *tomeStoreService) ListSeries(ctx context.Context) ([]*ent.Series, error) {
	series, err := s.Client.Series.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func (s *tomeStoreService) GetSeries(ctx context.Context, id int) (*ent.Series, error) {
	series, err := s.Client.Series.Query().Where(series.ID(id)).WithBooks(func(bq *ent.BookQuery) {
		bq.Order(book.ByNumber())
	}).Only(ctx)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func (s *tomeStoreService) AddSeries(ctx context.Context, name string, url string) (*ent.Series, error) {
	series, err := s.Client.Series.Create().SetName(name).SetURL(url).Save(ctx)
	if err != nil {
		return nil, err
	}
	return series, nil
}

func (s *tomeStoreService) DeleteSeries(ctx context.Context, id int) error {
	s.Client.Series.DeleteOneID(id).Exec(ctx)
	return nil
}

func (s *tomeStoreService) UpdateSeries(ctx context.Context, id int, series SeriesUpdate) (*ent.Series, error) {
	req := s.Client.Series.UpdateOneID(id)
	if series.Name != nil {
		req.SetName(*series.Name)
	}
	if series.URL != nil {
		req.SetURL(*series.URL)
	}
	if series.Cover != nil {
		req.SetCover(*series.Cover)
	}
	if series.Author != nil {
		req.SetAuthor(*series.Author)
	}
	if series.Description != nil {
		req.SetDescription(*series.Description)
	}
	updatedSeries, err := req.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedSeries, nil
}

func (s *tomeStoreService) CountBooksInSeries(ctx context.Context, id int) (int, int, error) {
	series, err := s.Client.Series.Query().Where(series.ID(id)).WithBooks().Only(ctx)
	if err != nil {
		return 0, 0, err
	}

	totalBooks := len(series.Edges.Books)
	ownedBooks := 0
	for _, book := range series.Edges.Books {
		if book.Owned {
			ownedBooks++
		}
	}
	return totalBooks, ownedBooks, nil
}

type SeriesUpdate struct {
	Name        *string
	URL         *string
	Cover       *string
	Author      *string
	Description *string
}
