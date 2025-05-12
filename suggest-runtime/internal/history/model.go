package history

import (
	"time"
)

type QueryTimestamp struct {
	Query string    `json:"query"`
	Time  time.Time `json:"time"`
}
