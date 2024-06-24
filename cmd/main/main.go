package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testapi/internal/nyct"
	"testapi/internal/stations"
	"time"
)

func main() {

	nyctStopId := "354"
	nyctDirection := "N"
	trainsToReturn := 2

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

	trips := nyct.ProcessTrips(nyctData, localStation.GTFS_ID, nyctDirection)

	fmt.Printf("Upcoming trains for %s:\n", localStation.StationName)
	for i, trip := range trips {
		if i >= trainsToReturn {
			break
		}
		timeUntilArrival := time.Until(trip.ArrivalTime)
		fmt.Printf("Route %s %s arriving in %v\n", trip.RouteID, trip.Direction, timeUntilArrival.Round(time.Minute))
	}
}
