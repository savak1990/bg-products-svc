package products

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

type Product struct {
	ID         string    `json:"id"` // Autoincremented integer represented as string
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
}

type Repo interface {
	List() []Product
	Create(Product) (Product, error)
}

type InMemoryStore struct {
	mu     sync.RWMutex
	data   map[string]Product
	nextID int64
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data:   make(map[string]Product),
		nextID: 1,
	}
}

func (s *InMemoryStore) List() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Product, 0, len(s.data))
	for _, p := range s.data {
		out = append(out, p)
	}
	return out
}

func (s *InMemoryStore) Create(p Product) (Product, error) {
	if p.Name == "" {
		return Product{}, errors.New("name is required")
	}

	if p.PriceCents < 0 {
		return Product{}, errors.New("price must be non-negative")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	p.ID = strconv.FormatInt(id, 10)
	p.CreatedAt = time.Now().UTC()
	s.data[p.ID] = p

	return p, nil
}

var _ Repo = (*InMemoryStore)(nil)
