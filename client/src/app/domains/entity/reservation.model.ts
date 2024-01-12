/*
	Id                gocql.UUID `json:"id"`
	UserID            string     `json:"userId"`
	AccommodationID   string     `json:"accommodationId"`
	StartDate         string     `json:"startDate"`
	EndDate           string     `json:"endDate"`
	Username          string     `json:"username"`
	AccommodationName string     `json:"accommodationName"`
	Location          string     `json:"location"`
	Price             int        `json:"price"`
	NumberOfDays      int        `json:"numOfDays"`
	Continent         string     `json:"continent"`
	DateRange         []string   `json:"dateRange"`
	IsActive          bool       `json:"isActive"`
	Country           string     `json:"country"`
	HostID            string     `json:"hostId"`
*/

export interface Reservation {
    id: string
    userID: string
    accommodationID: string
    startDate: string
    endDate: string
    username: string
    accommodationName: string
    location: string
    price: number
    numberOfDays: number
    continent: string
    dateRange: string[]
    isActive: boolean
    country: string
    hostID: string
}