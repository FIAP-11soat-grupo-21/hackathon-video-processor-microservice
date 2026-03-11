package domain

import "time"

type VideoChunk struct {
	VideoID   string
	ChunkPart int
	Status    string
	UserID    string
	UserName  string
	UserEmail string
	UpdatedAt time.Time
}

type User struct {
	ID    string
	Name  string
	Email string
}
