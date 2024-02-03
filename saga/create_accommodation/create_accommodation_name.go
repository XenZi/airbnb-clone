package create_accommodation

type SagaCommandType int8

const (
	CreateAccommodation SagaCommandType = iota
	CreateAvailability
	ApproveAccommodation
	DenyAccommodation
	UpdateAccommodation
	RollbackAccommodation
	UnknownCommand
)

type SagaReplyType int8

type AvailableAccommodationDates struct {
	AccommodationId string
	DateRange       []string
	Location        string
	Price           int
}

type SendCreateAccommodationAvailability struct {
	AccommodationID string
	Location        string
	DateRange       []AvailableAccommodationDates
}

const (
	AccommodationCreated SagaReplyType = iota
	AvailabilityCreated
	AccommodationRolledBack
	AvailabilityNotCreated
	UnknownReply
)

type SagaCommand struct {
	Type    SagaCommandType
	Payload interface{}
}

type SagaReply struct {
	Type    SagaReplyType
	Payload interface{}
}
