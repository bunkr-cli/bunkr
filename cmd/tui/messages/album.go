package messages

import (
	"context"
	"github.com/bunkr-cli/bunkr/internal/scrape"
	tea "github.com/charmbracelet/bubbletea"
	"time"
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

var hydrateTries = 5

func HydrateAlbum(ctx context.Context, album *scrape.Album) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()

		for i := 1; i <= hydrateTries; i += 1 {
			select {
			case <-ctx.Done():
				return nil
			default:
				err := scrape.DefaultScraper.HydrateAlbum(ctx, album)
				if err == nil {
					return AlbumHydratedMsg{Album: album}
				}
				time.Sleep(time.Duration(i) * 2 * time.Second)
			}
		}

		return nil
	}
}
