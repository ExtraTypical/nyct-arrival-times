package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testapi/internal/getenv"
	"testapi/internal/nyct"
	"testapi/internal/stations"
	"time"
)

func main() {

	nyctStopId, err := getenv.GetEnvVariable("STOP_ID")
	if err != nil {
		log.Fatal(err)
	}

	stationsResponse, err := stations.LoadStations()
	if err != nil {
		log.Fatal(err)
	}

	var stationsData stations.Response
	if err := json.Unmarshal(stationsResponse, &stationsData); err != nil {
		log.Fatal(err)
	}

	localStation, err := stations.LoadLocalStation(nyctStopId, stationsData)
	if err != nil {
		log.Fatal(err)
	}

	nyctData, err := nyct.CallNYCT()
	if err != nil {
		log.Fatal(err)
	}

	trips := nyct.ProcessTrips(nyctData, localStation.GTFS_ID)

	fmt.Printf("Upcoming trains for %s:\n", localStation.StationName)
	for i, trip := range trips {
		if i >= 2 {
			break
		}
		timeUntilArrival := time.Until(trip.ArrivalTime)
		fmt.Printf("Route %s arriving in %v\n", trip.RouteID, timeUntilArrival.Round(time.Minute))
	}
}
