package scrape

import (
	"github.com/dustin/go-humanize"
	zone "github.com/lrstanley/bubblezone"
	"github.com/skratchdot/open-golang/open"
	"net/url"
	"strconv"
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

	Hydrated  bool   `json:"-"`
	Enabled   uint8  `json:"enabled"`
	Public    uint8  `json:"public"`
	Desc      string `json:"description"`
	Timestamp uint64 `json:"timestamp"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	NotFound  bool   `json:"not_found"`
	Files     []File `json:"files"`
}

func (a *Album) Title() string {
	return zone.Mark(a.Identifier, a.Name)
}
func (a *Album) URL() *url.URL { return BaseUrl.JoinPath("a", a.Identifier) }
func (a *Album) Description() string {
	if !a.Hydrated {
		return ""
	}

	var totalSize uint64
	for _, file := range a.Files {
		size, _ := strconv.Atoi(file.Size)
		totalSize += uint64(size)
	}
	return strconv.Itoa(len(a.Files)) + " files (" + humanize.Bytes(totalSize) + ")"

}
func (a *Album) FilterValue() string { return zone.Mark(a.Identifier, a.Name) }

func (a *Album) Open() error {
	return open.Run(a.URL().String())
}
