package tomeStore

import (
	"context"
	"time"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/diconico07/TomeTracker/pkg/ent/book"
)

// BookService describes the Book service.
type BookService interface {
	ListBooks(ctx context.Context) ([]*ent.Book, error)
	GetBook(ctx context.Context, id string) (*ent.Book, error)
	AddBook(ctx context.Context, name string, url string, seriesID int, tomeNumber int, isbn string) (*ent.Book, error)
	UpsertBook(ctx context.Context, book *ent.Book) (*ent.Book, error)
	DeleteBook(ctx context.Context, id string) error
	UpdateBook(ctx context.Context, id string, book BookUpdate) (*ent.Book, error)
	ListPlannedBooks(ctx context.Context) ([]*ent.Book, error)
	ListMissingBooks(ctx context.Context) ([]*ent.Book, error)
}

// NewBookService creates a new Book service.
func NewBookService(client *ent.Client) BookService {
	return &tomeStoreBookService{
		Client: client,
	}
}

type tomeStoreBookService struct {
	Client *ent.Client
}

func (s *tomeStoreBookService) ListBooks(ctx context.Context) ([]*ent.Book, error) {
	books, err := s.Client.Book.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *tomeStoreBookService) ListPlannedBooks(ctx context.Context) ([]*ent.Book, error) {
	books, err := s.Client.Book.Query().Where(book.Owned(false), book.ReleasedAtGT(time.Now().AddDate(0, -1, 0))).Order(ent.Asc(book.FieldReleasedAt)).All(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *tomeStoreBookService) ListMissingBooks(ctx context.Context) ([]*ent.Book, error) {
	books, err := s.Client.Book.Query().Where(book.Owned(false), book.ReleasedAtLT(time.Now())).Order(ent.Asc(book.FieldReleasedAt)).All(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (s *tomeStoreBookService) GetBook(ctx context.Context, id string) (*ent.Book, error) {
	book, err := s.Client.Book.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *tomeStoreBookService) AddBook(ctx context.Context, title string, url string, seriesID int, tomeNumber int, isbn string) (*ent.Book, error) {
	book, err := s.Client.Book.Create().SetTitle(title).SetURL(url).SetSeriesID(seriesID).SetNumber(tomeNumber).SetID(isbn).SetOwned(false).Save(ctx)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *tomeStoreBookService) DeleteBook(ctx context.Context, id string) error {
	s.Client.Book.DeleteOneID(id).Exec(ctx)
	return nil
}

func (s *tomeStoreBookService) UpsertBook(ctx context.Context, book *ent.Book) (*ent.Book, error) {
	err := s.Client.Book.Create().
		SetID(book.ID).
		SetTitle(book.Title).
		SetURL(book.URL).
		SetNumber(book.Number).
		SetOwned(book.Owned).
		SetReleasedAt(book.ReleasedAt).
		SetSeriesID(book.SeriesID).
		SetCover(book.Cover).
		OnConflictColumns("id").
		UpdateNewValues().Exec(ctx)
	if err != nil {
		return nil, err
	}
	createdBook, err := s.Client.Book.Get(ctx, book.ID)
	if err != nil {
		return nil, err
	}
	return createdBook, nil
}

func (s *tomeStoreBookService) UpdateBook(ctx context.Context, id string, book BookUpdate) (*ent.Book, error) {
	req := s.Client.Book.UpdateOneID(id)
	if book.Title != nil {
		req.SetTitle(*book.Title)
	}
	if book.URL != nil {
		req.SetURL(*book.URL)
	}
	if book.Number != nil {
		req.SetNumber(*book.Number)
	}
	if book.Owned != nil {
		req.SetOwned(*book.Owned)
	}
	if book.ReleasedAt != nil {
		req.SetReleasedAt(*book.ReleasedAt)
	}
	if book.SeriesID != nil {
		req.SetSeriesID(*book.SeriesID)
	}
	if book.Cover != nil {
		req.SetCover(*book.Cover)
	}
	updatedBook, err := req.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updatedBook, nil
}

type BookUpdate struct {
	Title      *string
	URL        *string
	Number     *int
	Owned      *bool
	ReleasedAt *time.Time
	SeriesID   *int
	Cover      *string
}
