package performancefinder

import (
	"github.com/alexnel24/concurrency-opry/internal/models"
	"github.com/alexnel24/concurrency-opry/internal/store"
)

type PerformanceFinder struct {
	stores *store.Stores
}

func NewPerformanceFinder(stores *store.Stores) *PerformanceFinder {
	return &PerformanceFinder{stores: stores}
}

func (pf *PerformanceFinder) FindArtistPerformances(artistNames []string, filter string) ([]*models.Performance, error) {
	artists := make([]*models.Artist, 0, len(artistNames))
	for _, name := range artistNames {
		if artist, exists := pf.stores.ArtistStore.GetArtist(name); exists {
			artists = append(artists, artist)
		}
	}

	switch filter {
	case models.FilterUpcoming:
		return pf.stores.PerformanceStore.GetArtistPerformancesByUpcoming(pf.stores.DB(), artists, true)
	case models.FilterPast:
		return pf.stores.PerformanceStore.GetArtistPerformancesByUpcoming(pf.stores.DB(), artists, false)
	default:
		return pf.stores.PerformanceStore.GetAllArtistPerformances(pf.stores.DB(), artists)
	}
}
