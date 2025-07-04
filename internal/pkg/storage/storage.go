package storage

import (
	"encoding/json"
	"htmx-example/internal/pkg/models"
	"os"
	"sync"
)

type JsonStorage struct {
	filePath string
	mutex    sync.Mutex
}

func NewJsonStorage(filePath string) *JsonStorage {
	return &JsonStorage{filePath: filePath}
}

func (s *JsonStorage) Read() (*models.Companies, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.NewCompanies(), nil
		}
		return nil, err
	}

	var companies models.Companies
	if err := json.Unmarshal(data, &companies); err != nil {
		return nil, err
	}
	return &companies, nil
}

func (s *JsonStorage) Write(companies *models.Companies) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.MarshalIndent(companies, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
