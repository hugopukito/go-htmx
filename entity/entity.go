package entity

import (
	"time"

	"github.com/google/uuid"
)

type Dog struct {
	ID           uuid.UUID
	Name         string
	Score        int
	DateCreation time.Time
}

type DogTmpl struct {
	Dog
	Date string
}
