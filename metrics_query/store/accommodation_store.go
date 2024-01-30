package store

import (
	"errors"
	"metrics_query/domain"
)

type AccommodationStore struct {
	accommodations map[string]*domain.Accommodation
}

func NewAccommodationStore() domain.AccommodationStore {
	return AccommodationStore{make(map[string]*domain.Accommodation)}
}
func (a AccommodationStore) Create(accommodation domain.Accommodation) error {
	_, ok := a.accommodations[accommodation.Id]
	if ok {
		return errors.New("accommodation already exists")
	}
	a.accommodations[accommodation.Id] = &accommodation
	return nil
}

func (a AccommodationStore) ReadAll() []domain.Accommodation {
	accommodations := make([]domain.Accommodation, 0, len(a.accommodations))
	for _, val := range a.accommodations {
		accommodations = append(accommodations, *val)
	}
	return accommodations
}

func (a AccommodationStore) Read(id string) (domain.Accommodation, error) {
	acc, ok := a.accommodations[id]
	if !ok {
		return domain.Accommodation{}, errors.New("accommodation not found")
	}
	return *acc, nil
}

func (a AccommodationStore) Update(accommodation domain.Accommodation) error {
	_, ok := a.accommodations[accommodation.Id]
	if !ok {
		return errors.New("accommodation not found")
	}
	a.accommodations[accommodation.Id] = &accommodation
	return nil
}
