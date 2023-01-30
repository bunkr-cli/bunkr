package messages

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/bunkr-cli/bunkr/internal/scrape"
)

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
