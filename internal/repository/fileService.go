package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type Service struct {
	mu       sync.Mutex
	fileName string
}

func NewFileRepo(fileName string) *Service {
	s := &Service{
		fileName: fileName,
	}
	return s
}

func (s *Service) Save(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ioutil.WriteFile(s.fileName, data, 0644); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (s *Service) Load() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.fileName); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := ioutil.ReadFile(s.fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return data, nil
}
