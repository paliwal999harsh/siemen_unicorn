package storage

import (
	"unicorn/model"
)

type UnicornStorage interface {
	SaveUnicorn(model.Unicorn)
	GetUnicorns(int) []model.Unicorn
}

type InMemoryUnicornStorage struct {
	unicorns chan model.Unicorn
}

func NewInMemoryUnicornStorage() *InMemoryUnicornStorage {
	return &InMemoryUnicornStorage{unicorns: make(chan model.Unicorn, 10)}
}

func (us *InMemoryUnicornStorage) SaveUnicorn(u model.Unicorn) {
	us.unicorns <- u
}

func (us *InMemoryUnicornStorage) GetUnicorns(amount int) []model.Unicorn {
	if amount > len(us.unicorns) {
		amount = len(us.unicorns)
	}

	result := make([]model.Unicorn, amount)
	for i := 0; i < amount; i++ {
		u := <-us.unicorns
		result[i] = u
	}
	return result
}
