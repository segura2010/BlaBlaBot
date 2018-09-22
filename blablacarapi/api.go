package blablacarapi

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"

	"BlaBlaBot/config"
)

var (
	API_HOST = "https://public-api.blablacar.com/api"
)


func GetTrips(from, to, locale, currency, date string) (TripsResponse, error){

	url := "https://public-api.blablacar.com/api/v2/trips?fn="+ from +"&tn="+ to +"&locale="+ locale +"&_format=json&cur="+ currency +"&db="+ date +"&sort=trip_date&order=asc&limit=100"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("key", config.GetInstance().ApiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil{
		return TripsResponse{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		return TripsResponse{}, err
	}

	fmt.Printf(string(body))

	var trip TripsResponse
	err = json.Unmarshal(body, &trip)
	if err != nil{
		return TripsResponse{}, err
	}

	return trip, nil
}
