package cmd

import "strings"

type Storage interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Get(key string) (string, bool) {
	val, ok := s.data[strings.ToLower(key)]
	return val, ok
}

func (s *InMemoryStorage) Set(key, value string) {
	s.data[strings.ToLower(key)] = value
}
