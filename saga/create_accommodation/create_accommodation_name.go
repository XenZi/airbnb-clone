package create_accommodation

type SagaCommandType int8

// KOMANDE

type SagaReplyType int8

type AvailableAccommodationDates struct { // Pomocna struktura za SendCreateAccommodationAvailability
	AccommodationId string
	DateRange       []string
	Location        string
	Price           int
}

type SendCreateAccommodationAvailability struct { // ekvivalent sa Order Details
	AccommodationID string
	Location        string
	DateRange       []AvailableAccommodationDates
}

type CreateAccommodationCommandType int8

// komande
const (
	CreateAvailability CreateAccommodationCommandType = iota
	DenyAccommodation
	UpdateAccommodation
	RollbackAccommodation
	UnknownCommand
)

type CreateAccommodationReplyType int8

// REPLY
const (
	AvailabilityCreated CreateAccommodationReplyType = iota
	AvailabilityNotCreated
	UnknownReply
)

type CreateAccommodationCommand struct {
	Type    CreateAccommodationCommandType
	Payload SendCreateAccommodationAvailability
}

type CreateAccommodationReply struct {
	Type    CreateAccommodationReplyType
	Payload SendCreateAccommodationAvailability
}
