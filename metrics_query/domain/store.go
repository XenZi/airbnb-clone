package domain

type AccommodationStore interface {
	Create(acc Accommodation, collection string) error
	Read(id, collection string) (*Accommodation, error)
	Update(acc Accommodation, collection string) error
}
