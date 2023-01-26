package scrape

import (
	zone "github.com/lrstanley/bubblezone"
	"github.com/skratchdot/open-golang/open"
	"net/url"
)

type AlbumsPage struct {
	Props struct {
		PageProps struct {
			Albums []*Album `json:"albums"`
		} `json:"pageProps"`
	} `json:"props"`
}

type AlbumPage struct {
	Props struct {
		PageProps struct {
			Album *Album `json:"album"`
		} `json:"pageProps"`
	} `json:"props"`
}

type Album struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`

	Enabled   uint8  `json:"enabled"`
	Public    uint8  `json:"public"`
	Desc      string `json:"description"`
	Timestamp uint64 `json:"timestamp"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	NotFound  bool   `json:"not_found"`
	Files     []File `json:"files"`
}

func (a *Album) Title() string       { return zone.Mark(a.Name, a.Name) }
func (a *Album) URL() *url.URL       { return BaseUrl.JoinPath("a", a.Identifier) }
func (a *Album) Description() string { return a.URL().String() }
func (a *Album) FilterValue() string { return zone.Mark(a.Name, a.Name) }

func (a *Album) Open() error {
	return open.Run(a.URL().String())
}
