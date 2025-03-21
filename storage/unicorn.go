package storage

import (
	"sync"
	"unicorn/model"
)

type UnicornStore interface {
	SaveUnicorn(model.Unicorn)
	GetUnicorns(int) []model.Unicorn
	AvailableUnicorns() int
	Capacity() int
	DecreaseCapacity(int)
}

type InMemoryUnicornStore struct {
	unicorns []model.Unicorn
	capacity int
	sync.Mutex
}

func NewInMemoryUnicornStore() UnicornStore {
	return &InMemoryUnicornStore{unicorns: make([]model.Unicorn, 0, 100), capacity: 0}
}

func (us *InMemoryUnicornStore) SaveUnicorn(u model.Unicorn) {
	us.Lock()
	defer us.Unlock()
	us.unicorns = append([]model.Unicorn{u}, us.unicorns...)
	us.capacity++
}

func (us *InMemoryUnicornStore) GetUnicorns(amount int) []model.Unicorn {
	us.Lock()
	defer us.Unlock()

	if len(us.unicorns) == 0 {
		return nil
	}
	unicorns := us.unicorns[len(us.unicorns)-amount:]
	us.unicorns = us.unicorns[:len(us.unicorns)-amount]
	return unicorns
}

func (us *InMemoryUnicornStore) AvailableUnicorns() int {
	us.Lock()
	defer us.Unlock()

	return len(us.unicorns)
}

func (us *InMemoryUnicornStore) Capacity() int {
	us.Lock()
	defer us.Unlock()

	return us.capacity
}

func (us *InMemoryUnicornStore) DecreaseCapacity(capacity int) {
	us.Lock()
	defer us.Unlock()

	us.capacity -= capacity
}
