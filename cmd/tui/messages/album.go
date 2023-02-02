package messages

import (
	"context"
	"github.com/bunkr-cli/bunkr/internal/scrape"
	tea "github.com/charmbracelet/bubbletea"
)

type AlbumsReadyMsg struct {
	Albums []*scrape.Album
}

func ListAlbums(ctx context.Context, force bool) tea.Cmd {
	return func() tea.Msg {
		albums, err := scrape.DefaultScraper.Albums(ctx, force)
		if err != nil {
			return NewErrMsg("Failed to fetch albums", err)
		}

		return AlbumsReadyMsg{Albums: albums}
	}
}

type AlbumHydratedMsg struct {
	Album *scrape.Album
}

func HydrateAlbum(ctx context.Context, album *scrape.Album) tea.Cmd {
	return func() tea.Msg {
		if err := scrape.DefaultScraper.HydrateAlbum(ctx, album); err != nil {
			return nil
		}

		return AlbumHydratedMsg{Album: album}
	}
}
