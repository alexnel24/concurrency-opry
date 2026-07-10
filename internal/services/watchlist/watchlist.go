package watchlist

import (
	"strings"

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
	// TODO: implement internal/store/watched_artist_store.go, then uncomment below
	// return wl.stores.WatchedArtistStore.GetWatchedArtists(ownerId)

	return []*models.WatchedArtist{}, nil // placeholder until watchlist storage exists
}

func (wl *Watchlist) RemoveWatchedArtist(name string, ownerId *int64) (found bool, err error) {
	// TODO: implement internal/store/watched_artist_store.go, then uncomment below
	// return wl.stores.WatchedArtistStore.RemoveWatchedArtist(name, ownerId)

	return false, nil // placeholder until watchlist storage exists
}

func (wl *Watchlist) AddWatchedArtist(name string, ownerId *int64, specificity string) ([]AddedWatchedArtist, error) {
	var newEntries []*models.WatchedArtist
	if artist, exists := wl.stores.ArtistStore.GetArtist(name); exists {
		newEntries = append(newEntries, models.NewWatchedArtistExactName(name, artist.Id))
	} else if specificity == "individual" {
		for _, word := range strings.Fields(name) {
			newEntries = append(newEntries, models.NewWatchedArtistFuzzyName(word))
		}
	} else {
		newEntries = append(newEntries, models.NewWatchedArtistFuzzyName(name))
	}

	for _, wa := range newEntries {
		wa.OwnerId = ownerId
	}

	results := make([]AddedWatchedArtist, 0, len(newEntries))
	for _, wa := range newEntries {
		// TODO: implement internal/store/watched_artist_store.go, then uncomment below
		// alreadyExists, err := wl.stores.WatchedArtistStore.AddWatchedArtist(wa)
		// if err != nil {
		// 	return nil, err
		// }
		alreadyExists := false // placeholder until watchlist storage exists

		results = append(results, AddedWatchedArtist{Entry: wa, AlreadyExists: alreadyExists})
	}

	return results, nil
}
