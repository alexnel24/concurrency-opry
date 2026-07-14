package handlers

type watchedArtistResponse struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	ArtistId      *int64 `json:"artist_id,omitempty"`
	OwnerId       *int64 `json:"owner_id,omitempty"`
	AlreadyExists bool   `json:"already_exists,omitempty"`
}

type artistPerformanceResponse struct {
	ArtistName string `json:"artist_name"`
	EventLink  string `json:"event_link"`
	EventTitle string `json:"event_title"`
	EventTime  string `json:"event_time"`
	Upcoming   bool   `json:"upcoming"`
}
