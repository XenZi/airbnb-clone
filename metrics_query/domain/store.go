package domain

type AccommodationStore interface {
	Create(account Accommodation) error
	ReadAll() []Accommodation
	Read(accountNumber string) (Accommodation, error)
	Update(account Accommodation) error
}
