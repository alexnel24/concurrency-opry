package models

type WatchedArtist struct {
	Id       int64
	Name     string
	ArtistId *int64
	OwnerId  *int64 // nil for now (single-user); will reference an owner/user id once multi-user support exists
}

func NewWatchedArtist(name string, artistId *int64, ownerId *int64) *WatchedArtist {
	return &WatchedArtist{Name: name, ArtistId: artistId, OwnerId: ownerId}
}
