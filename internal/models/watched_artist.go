package models

const (
	MatchTypeExact = "exact"
	MatchTypeFuzzy = "fuzzy"
)

type WatchedArtist struct {
	Id        int64
	Name      string
	MatchType string
	ArtistId  *int64
	OwnerId   *int64 // nil for now (single-user); will reference an owner/user id once multi-user support exists
}

func NewWatchedArtistFuzzyName(name string) *WatchedArtist {
	return &WatchedArtist{Name: name, MatchType: MatchTypeFuzzy}
}

func NewWatchedArtistExactName(name string, artistId int64) *WatchedArtist {
	return &WatchedArtist{Name: name, MatchType: MatchTypeExact, ArtistId: &artistId}
}
