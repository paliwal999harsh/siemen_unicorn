package storage

import (
	"sync"
	"unicorn/model"
	"unicorn/pkg/collection"
)

type UnicornStore interface {
	SaveUnicorn(model.Unicorn)
	GetUnicorns(int) []model.Unicorn
	AvailableUnicorns() int
	Capacity() int
	DecreaseCapacity(int)
	IsAtCapacity() bool
}

const MaxStoreCapacity = 100

type InMemoryUnicornStore struct {
	unicorns collection.Stack[model.Unicorn]
	capacity int
	sync.Mutex
}

func (us *InMemoryUnicornStore) SaveUnicorn(u model.Unicorn) {
	us.Lock()
	defer us.Unlock()
	_, _ = us.unicorns.Push(u)
	us.capacity++
}

func (us *InMemoryUnicornStore) GetUnicorns(amount int) []model.Unicorn {
	us.Lock()
	defer us.Unlock()

	if us.unicorns.Empty() {
		return nil
	}
	var unicorns []model.Unicorn
	for range amount {
		unicorn, _ := us.unicorns.Pop()
		unicorns = append(unicorns, unicorn)
	}
	return unicorns
}

func (us *InMemoryUnicornStore) AvailableUnicorns() int {
	us.Lock()
	defer us.Unlock()

	return us.unicorns.Size()
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

func (us *InMemoryUnicornStore) IsAtCapacity() bool {
	us.Lock()
	defer us.Unlock()

	return us.unicorns.Size() >= MaxStoreCapacity
}

func NewInMemoryUnicornStore() UnicornStore {
	return &InMemoryUnicornStore{unicorns: collection.NewSliceStack[model.Unicorn](), capacity: 0}
}
