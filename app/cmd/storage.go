package cmd

import (
	"log"
	"strings"
	"time"
)

type Storage interface {
	Get(key string) (string, bool)
	Set(key, value string, exp *time.Time)
}

type InMemoryStorage struct {
	data map[string]*CacheRecord
}

type CacheRecord struct {
	K, V      string
	CreatedAt time.Time
	ExpireAt  time.Time
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]*CacheRecord),
	}
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	val, ok := s.data[strings.ToLower(key)]
	if ok && val.ExpireAt.Before(time.Now()) {
		go func() {
			log.Printf("DEBUG: DELETE key='%s', value='%s'", val.K, val.V)
			s.deleteRecord(val)
		}()
		return "", false
	}
	if ok {
		return val.V, true
	}
	return "", false
}

func (s *InMemoryStorage) Set(key, value string, exp *time.Time) {
	if exp == nil {
		s.data[strings.ToLower(key)] = &CacheRecord{K: key, V: value, CreatedAt: time.Now(), ExpireAt: time.Now().AddDate(1, 0, 0)}
	} else {
		s.data[strings.ToLower(key)] = &CacheRecord{K: key, V: value, CreatedAt: time.Now(), ExpireAt: *exp}
	}
}

func (s *InMemoryStorage) deleteRecord(cacheRecord *CacheRecord) {
	delete(s.data, cacheRecord.K)
}
