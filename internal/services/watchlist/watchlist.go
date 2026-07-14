package watchlist

import (
	"github.com/alexnel24/concurrency-opry/internal/models"
	"github.com/alexnel24/concurrency-opry/internal/store"
)

type AddedWatchedArtist struct {
	Entry         *models.WatchedArtist
	AlreadyExists bool
}

type Watchlist struct {
	stores *store.Stores
}

func NewWatchlist(stores *store.Stores) *Watchlist {
	return &Watchlist{stores: stores}
}

func (wl *Watchlist) ListWatchedArtists(ownerId *int64) ([]*models.WatchedArtist, error) {
	return wl.stores.WatchlistStore.GetWatchedArtists(wl.stores.DB(), ownerId)
}

func (wl *Watchlist) RemoveWatchedArtist(name string, ownerId *int64) (found bool, err error) {
	return wl.stores.WatchlistStore.RemoveWatchedArtist(wl.stores.DB(), name, ownerId)
}

func (wl *Watchlist) AddWatchedArtist(name string, ownerId *int64) (AddedWatchedArtist, error) {
	var artistId *int64
	if artist, exists := wl.stores.ArtistStore.GetArtist(name); exists {
		artistId = &artist.Id
	}

	entry := models.NewWatchedArtist(name, artistId, ownerId)
	alreadyExists, err := wl.stores.WatchlistStore.AddWatchedArtist(wl.stores.DB(), entry)
	if err != nil {
		return AddedWatchedArtist{}, err
	}

	return AddedWatchedArtist{Entry: entry, AlreadyExists: alreadyExists}, nil
}

func (wl *Watchlist) GetWatchedArtistPerformances(ownerId *int64, filter string) ([]*models.Performance, error) {
	watched, err := wl.stores.WatchlistStore.GetWatchedArtists(wl.stores.DB(), ownerId)
	if err != nil {
		return nil, err
	}

	var artists []*models.Artist
	for _, wa := range watched {
		if artist, exists := wl.stores.ArtistStore.GetArtist(wa.Name); exists {
			artists = append(artists, artist)
		}
	}

	if filter == models.FilterUpcoming {
		return wl.stores.PerformanceStore.GetArtistPerformancesByUpcoming(wl.stores.DB(), artists, true)
	}
	if filter == models.FilterPast {
		return wl.stores.PerformanceStore.GetArtistPerformancesByUpcoming(wl.stores.DB(), artists, false)
	}
	return wl.stores.PerformanceStore.GetAllArtistPerformances(wl.stores.DB(), artists)
}
