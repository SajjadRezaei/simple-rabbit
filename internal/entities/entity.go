package entities

import "time"

type Message struct {
	Content    string
	Expiration *time.Time
	Key        []byte
	Priority   int
}
