package bloom

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/willf/bloom"
)

// BloomFilterManager adalah wrapper untuk bloom filter dan file persistence
type BloomFilterManager struct {
	filter   *bloom.BloomFilter // filter utama
	filePath string             // path file untuk persistence
	mu       sync.RWMutex       // mutex untuk concurrency
	n        uint               // estimasi jumlah item
	fp       float64            // false positive rate
}

// ensureDir memastikan direktori filePath sudah ada
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// NewBloomFilterManager membuat instance baru BloomFilterManager
func NewBloomFilterManager(filePath string, n uint, fp float64) (*BloomFilterManager, error) {
	if err := ensureDir(filePath); err != nil {
		return nil, err
	}
	manager := &BloomFilterManager{
		filePath: filePath,
		n:        n,
		fp:       fp,
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

// Add menambahkan email ke bloom filter
func (m *BloomFilterManager) Add(email string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.filter.Add([]byte(email))
}

// Test mengecek apakah email ada di bloom filter
func (m *BloomFilterManager) Test(email string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.filter.Test([]byte(email))
}

// Save menyimpan bloom filter ke file
func (m *BloomFilterManager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if err := ensureDir(m.filePath); err != nil {
		return err
	}
	file, err := os.Create(m.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = m.filter.WriteTo(file)
	return err
}

// Regenerate membuat ulang bloom filter dari daftar email
func (m *BloomFilterManager) Regenerate(emails []string) {
	newFilter := bloom.NewWithEstimates(m.n, m.fp)
	for _, email := range emails {
		newFilter.Add([]byte(email))
	}
	m.mu.Lock()
	m.filter = newFilter
	m.mu.Unlock()
}
