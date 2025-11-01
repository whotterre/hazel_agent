package store

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Birthday struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Month     int       `json:"month"`
	Day       int       `json:"day"`
	CreatedAt time.Time `json:"created_at"`
}

type BirthdayStore struct {
	mu        sync.RWMutex
	birthdays map[string]Birthday
	file      string
}

func NewBirthdayStore(filename string) *BirthdayStore {
	store := &BirthdayStore{
		birthdays: make(map[string]Birthday),
		file:      filename,
	}
	store.load()
	return store
}

func (bs *BirthdayStore) AddBirthday(name, date string) (string, error) {
	var t time.Time
	var err error

	if len(date) == 5 {
		t, err = time.Parse("01-02", date)
		if err != nil {
			return "", err
		}
	} else {
		t, err = time.Parse("2006-01-02", date)
		if err != nil {
			return "", err
		}
	}

	birthday := Birthday{
		ID:        uuid.New().String(),
		Name:      name,
		Month:     int(t.Month()),
		Day:       t.Day(),
		CreatedAt: time.Now(),
	}

	bs.mu.Lock()
	bs.birthdays[birthday.ID] = birthday
	bs.mu.Unlock()

	bs.save()
	return birthday.ID, nil
}

func (bs *BirthdayStore) List() []Birthday {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	birthdays := make([]Birthday, 0, len(bs.birthdays))
	for _, b := range bs.birthdays {
		birthdays = append(birthdays, b)
	}
	return birthdays
}

func (bs *BirthdayStore) save() {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	data, _ := json.MarshalIndent(bs.birthdays, "", "  ")
	os.WriteFile(bs.file, data, 0644)
}

func (bs *BirthdayStore) load() {
	data, err := os.ReadFile(bs.file)
	if err != nil {
		return
	}
	json.Unmarshal(data, &bs.birthdays)
}
