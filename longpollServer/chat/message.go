package chat

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id        int
	From      uuid.UUID
	To        uuid.UUID
	Content   string
	CreatedAt time.Time
}
