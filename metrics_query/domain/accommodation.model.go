package domain

type Accommodation struct {
	Id                                 string
	OnScreenTime                       float64
	NumberOfVisits                     uint32
	NotClosedEventTimeStamps           map[string]string
	LastAppliedUserJoinedEventNumber   int64
	LastAppliedUserLeftEventNumber     int64
	NumberOfReservations               uint32
	LastAppliedUserReservedEventNumber int64
	NumberOfRatings                    uint32
	LastAppliedUserRatedEventNumber    int64
}
