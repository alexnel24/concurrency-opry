package models

type Artist struct {
	Id   int64
	Name string
}

func NewArtist(name string) *Artist {
	return &Artist{Name: name}	
}
