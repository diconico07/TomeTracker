package yenpressfetcher

import (
	"context"
	"fmt"
	"net/http"
	nurl "net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"

	"github.com/diconico07/TomeTracker/pkg/ent"
	"github.com/diconico07/TomeTracker/pkg/tomeStore"
)

// YenPressFetcher is an implementation of the Fetcher interface for Yen Press.
type YenPressFetcher struct {
}

// New creates a new YenPressFetcher.
func New() *YenPressFetcher {
	return &YenPressFetcher{}
}

// IsURLSupported checks if the given URL is supported by this fetcher.
func (f *YenPressFetcher) IsURLSupported(url string) bool {
	// Implement logic to check if the URL belongs to Yen Press
	u, err := nurl.Parse(url)
	if err != nil {
		return false
	}
	if u.Host != "yenpress.com" {
		return false
	}
	if !strings.HasPrefix(u.Path, "/series/") {
		return false
	}
	return true
}

// ProcessSeries processes a series from Yen Press and returns a list of books.
func (f *YenPressFetcher) ProcessSeries(ctx context.Context, series *ent.Series) ([]*ent.Book, tomeStore.SeriesUpdate, error) {
	resp, err := http.Get(series.URL)
	if err != nil {
		return nil, tomeStore.SeriesUpdate{}, fmt.Errorf("failed to fetch URL %s: %w", series.URL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, tomeStore.SeriesUpdate{}, fmt.Errorf("failed to fetch URL %s: %w", series.URL, err)
	}

	title := doc.Find(".series-heading h1.heading").Contents().First().Text()
	cover := doc.Find(".series-heading .heading-bg").AttrOr("data-src", "")

	authors := doc.Find(".series-heading .story-details p").FilterFunction(func(i int, s *goquery.Selection) bool {
		return strings.HasPrefix(s.Text(), "Author")
	}).Find("span").Map(func(i int, s *goquery.Selection) string {
		return s.Text()
	})
	author := strings.Join(authors, ", ")
	log.Debug().Str("Title", title).Str("Cover", cover).Str("Author", author).Msg("Parsed Series")

	savedBooks := make(map[string]*ent.Book)
	for _, book := range series.Edges.Books {
		savedBooks[book.ID] = book
	}

	var books []*ent.Book
	page := doc.Find("#volumes-list")
	newBooks, nextPage, err := parsePage(page, savedBooks, series.ID)
	if err != nil {
		return nil, tomeStore.SeriesUpdate{}, err
	}
	books = append(books, newBooks...)

	for nextPage != "" {
		log.Debug().Str("URL", nextPage).Msg("Fetching Next Page")
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextPage, nil)
		if err != nil {
			return nil, tomeStore.SeriesUpdate{}, fmt.Errorf("failed to create request for URL %s: %w", nextPage, err)
		}
		req.Header.Set("x-requested-with", "XMLHttpRequest")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, tomeStore.SeriesUpdate{}, fmt.Errorf("failed to fetch URL %s: %w", nextPage, err)
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, tomeStore.SeriesUpdate{}, err
		}
		newBooks, nextPage, err = parsePage(doc.Selection, savedBooks, series.ID)
		if err != nil {
			return nil, tomeStore.SeriesUpdate{}, err
		}
		books = append(books, newBooks...)
	}
	log.Debug().Int("Books", len(books)).Msg("Parsed Books")
	return books, tomeStore.SeriesUpdate{Name: &title, Cover: &cover, Author: &author}, nil
}

func parsePage(page *goquery.Selection, savedBooks map[string]*ent.Book, seriesID int) ([]*ent.Book, string, error) {
	var books []*ent.Book
	var err error
	var nextPage string

	page.Find("a").Each(func(i int, s *goquery.Selection) {
		url, exists := s.Attr("href")
		if !exists {
			log.Error().Msg("Invalid section")
			return
		}

		if url == "" {
			nextPage, _ = s.Attr("data-url")
			nextPage = "https://yenpress.com" + nextPage
			return
		}

		title := s.Find("span").Text()
		cover := s.Find("img").AttrOr("src", "")

		fullURL := "https://yenpress.com" + url

		log.Debug().Str("URL", fullURL).Msg("Parsing Volume")

		urlTitleParts := strings.Split(url[len("/titles/"):], "-")
		isbn := urlTitleParts[0]

		// Extract volume number from the URL
		var volume int
		for i, part := range urlTitleParts {
			if part == "vol" || part == "part" {
				volume, err = strconv.Atoi(urlTitleParts[i+1])
				if err != nil {
					log.Error().Err(err).Str("ISBN", isbn).Msg("Invalid volume number")
					return
				}
				break
			}
		}
		if saved := savedBooks[isbn]; saved != nil {
			if saved.Owned {
				log.Debug().Str("ISBN", isbn).Msg("Skipping Already Owned")
				return
			}
			if saved.UpdatedAt.After(saved.ReleasedAt.AddDate(0, 0, 1)) {
				log.Debug().Str("ISBN", isbn).Msg("Skipping Already Released")
				return
			}
		}
		releaseDate, err := getReleaseDate(fullURL)
		if err != nil {
			log.Error().Err(err).Str("ISBN", isbn).Msg("Failed to get release date")
			return
		}

		books = append(books, &ent.Book{
			ID:         isbn,
			Title:      title,
			URL:        fullURL,
			Number:     volume,
			Owned:      false,
			SeriesID:   seriesID,
			UpdatedAt:  time.Now(),
			ReleasedAt: releaseDate,
			Cover:      cover,
		})
	})
	log.Debug().Int("Books", len(books)).Msg("Parsed Books in Page")
	return books, nextPage, nil
}

func getReleaseDate(url string) (time.Time, error) {
	log.Debug().Str("URL", url).Msg("Getting Release Date")
	resp, err := http.Get(url)
	if err != nil {
		return time.Time{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return time.Time{}, err
	}

	releaseDate := doc.Find(".book-details .detail-info .detail-box").FilterFunction(func(i int, s *goquery.Selection) bool {
		return s.Find("span").Text() == "Release Date"
	}).Find(".info").First().Text()
	parsedReleaseDate, err := time.Parse("Jan 02, 2006", releaseDate)
	if err != nil {
		return time.Time{}, err
	}
	return parsedReleaseDate, nil
}
