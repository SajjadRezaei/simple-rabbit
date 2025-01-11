package storage

import (
	"encoding/json"
	"log"
	"time"

	"go.etcd.io/bbolt"

	"simple-rabbit/internal/entities"
)

type Storage struct {
	db *bbolt.DB
}

func NewStorage(path string) *Storage {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{db: db}
}

func (s *Storage) SaveMessage(queue string, message entities.Message) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(queue))
		if err != nil {
			return err
		}
		id, _ := bucket.NextSequence()
		msg, err := json.Marshal(message)
		if err != nil {
			return err
		}
		return bucket.Put(itob(int(id)), msg)
	})
}

func (s *Storage) LoadMessages(queue string) ([]entities.Message, error) {
	var messages []entities.Message
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(queue))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			msg := entities.Message{}
			err := json.Unmarshal(v, &msg)

			if err != nil {
				return err
			}

			msg.Key = k

			if msg.Expiration == nil || msg.Expiration.After(time.Now()) {
				messages = append(messages, msg)
			} else {
				s.DeleteMessage(queue, k)
			}

			return nil
		})
	})

	return messages, err
}

func (s *Storage) DeleteMessage(queue string, id []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(queue))
		if bucket == nil {
			return nil
		}
		return bucket.Delete(id)
	})
}
func itob(v int) []byte {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

func btoi(b []byte) int {
	if len(b) != 4 {
		panic("btoi requires a 4-byte slice")
	}
	return int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])
}
