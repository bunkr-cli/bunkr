package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bunkr-cli/bunkr/internal/scrape"
)

type AlbumsReadyMsg struct {
	Albums []*scrape.Album
}

func ListAlbums(force bool) tea.Cmd {
	return func() tea.Msg {
		albums, err := scrape.DefaultScraper.Albums(force)
		if err != nil {
			return NewErrMsg("Failed to fetch albums", err)
		}

		return AlbumsReadyMsg{Albums: albums}
	}
}

type AlbumHydratedMsg struct {
	Album *scrape.Album
}

func HydrateAlbum(album *scrape.Album) tea.Cmd {
	return func() tea.Msg {
		if err := scrape.DefaultScraper.HydrateAlbum(album); err != nil {
			return nil
		}

		return AlbumHydratedMsg{Album: album}
	}
}
