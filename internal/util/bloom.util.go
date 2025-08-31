package util

import (
	"os"
	"sync"

	"github.com/willf/bloom"
)

type BloomFilterManager struct {
	filter   *bloom.BloomFilter
	filePath string
	mu       sync.RWMutex
}

func NewBloomFilterManager(filePath string, n uint, fp float64) (*BloomFilterManager, error) {
	manager := &BloomFilterManager{
		filePath: filePath,
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			manager.filter = bloom.NewWithEstimates(n, fp)
			return manager, nil
		}
		return nil, err
	}
	defer file.Close()

	newFilter := bloom.NewWithEstimates(n, fp)
	if _, err := newFilter.ReadFrom(file); err != nil {
		return nil, err
	}
	manager.filter = newFilter
	return manager, nil
}

func (m *BloomFilterManager) Add(email string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.filter.Add([]byte(email))
}

func (m *BloomFilterManager) Test(email string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.filter.Test([]byte(email))
}

func (m *BloomFilterManager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.Create(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = m.filter.WriteTo(file)
	return err
}
