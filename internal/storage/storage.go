package storage

import (
	"go.etcd.io/bbolt"
	"log"
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

func (s *Storage) SaveMessage(queue string, message string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(queue))
		if err != nil {
			return err
		}
		id, _ := bucket.NextSequence()
		return bucket.Put(itob(int(id)), []byte(message))
	})
}

func (s *Storage) LoadMessage(queue string) ([]string, error) {
	var messages []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(queue))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(k, v []byte) error {
			messages = append(messages, string(k))
			return nil
		})
	})

	return messages, err
}

func (s *Storage) DeleteMessage(queue string, id int) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(queue))
		if bucket == nil {
			return nil
		}
		return bucket.Delete(itob(id))
	})
}
func itob(v int) []byte {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}

}
