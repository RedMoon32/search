package storage

import (
	"encoding/json"
	"log"

	"go.etcd.io/bbolt"
)

type BoltKVStore struct {
	db *bbolt.DB
}

func NewKVStore(path string) *BoltKVStore {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatalf("Failed to open bbolt database: %v", err)
	}
	return &BoltKVStore{db: db}
}

func (s *BoltKVStore) Save(key string, value map[string]interface{}) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("kvstore"))
		if err != nil {
			return err
		}
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), data)
	})
}

func (s *BoltKVStore) Load(key string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("kvstore"))
		if b == nil {
			return bbolt.ErrBucketNotFound
		}
		data := b.Get([]byte(key))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &result)
	})
	return result, err
}
