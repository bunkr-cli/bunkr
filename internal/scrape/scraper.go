package scrape

import (
	"compress/gzip"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

func NewScraper() *Scraper {
	return &Scraper{}
}

var DefaultScraper = NewScraper()

type Scraper struct {
	FetchedAt time.Time
	AlbumList []*Album
}

var UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
var BaseUrl = &url.URL{
	Scheme: "https",
	Host:   "bunkr.ru",
}

var ErrInvalidStatusCode = errors.New("invalid status code")

func (s *Scraper) Fetch(ctx context.Context, path string) (*http.Response, error) {
	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req.Header.Set("User-Agent", UserAgent)
		return nil
	}

	u, err := BaseUrl.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%w: %d", ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp, nil
}

func (s *Scraper) Albums(ctx context.Context, force bool) ([]*Album, error) {
	if !force {
		if err := s.Load(); err == nil {
			if time.Since(s.FetchedAt) < time.Hour {
				return s.AlbumList, nil
			}
		}
	}

	resp, err := s.Fetch(ctx, "/albums")
	if err != nil {
		return nil, err
	}

	nextData, err := GetNextData(resp.Body)
	if err != nil {
		return nil, err
	}
	CloseReader(resp.Body)

	var page AlbumsPage
	if err := json.Unmarshal([]byte(nextData), &page); err != nil {
		return page.Props.PageProps.Albums, err
	}

	s.AlbumList = page.Props.PageProps.Albums
	s.FetchedAt = time.Now()

	if err := s.Save(); err != nil {
		return s.AlbumList, err
	}

	return s.AlbumList, nil
}

func (s *Scraper) HydrateAlbum(ctx context.Context, album *Album) error {
	if album.Hydrated {
		return nil
	}

	resp, err := s.Fetch(ctx, path.Join("a", album.Identifier))
	if err != nil {
		return err
	}

	nextData, err := GetNextData(resp.Body)
	if err != nil {
		return err
	}
	CloseReader(resp.Body)

	var page AlbumPage
	page.Props.PageProps.Album = album
	if err := json.Unmarshal([]byte(nextData), &page); err != nil {
		return err
	}

	album.Hydrated = true
	return nil
}

func (s *Scraper) Save() error {
	f, err := os.Create("scraper.gob.gz")
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(f)

	if err := gob.NewEncoder(gz).Encode(s); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	return f.Close()
}

func (s *Scraper) Load() error {
	f, err := os.Open("scraper.gob.gz")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	if err := gob.NewDecoder(gz).Decode(s); err != nil {
		return err
	}

	return gz.Close()
}

var ErrNoNextData = errors.New("__NEXT_DATA__ was empty")

func GetNextData(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	text := doc.Find("#__NEXT_DATA__").First().Text()
	if text == "" {
		return "", ErrNoNextData
	}

	return text, nil
}

func CloseReader(r io.ReadCloser) {
	_, _ = io.Copy(io.Discard, r)
	_ = r.Close()
}
