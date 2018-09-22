package blablacarapi

import (
)

type TripPrice struct{
	Value float64
	Currency string
	Symbol string
	String_Value string
}

type Place struct{
	City_Name string
	Address string
	Latitude float64
	Longitude float64
	Country_Code string
}

type Trip struct{
	Permanent_Id string
	Links map[string]string
	Departure_Place Place
	Arrival_Place Place
	Departure_Date string
	Price TripPrice
	Price_With_Commission TripPrice
	SeatsLeft int64
	Seats int64
}

type TripsResponse struct{
	Trips []Trip
}
