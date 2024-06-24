package arriving

import (
	"encoding/json"
	"fmt"
	"strconv"
	format "testapi/internal/formatduration"
	"testapi/internal/nyct"
	"testapi/internal/stations"
	"time"
)

type Response struct {
	Trains []Train
}

type Train struct {
	RouteID    string
	Direction  string
	ArrivingIn string
}

func Arriving(nyctStopId int, nyctDirection string, trainsToReturn int) (Response, error) {

	// nyctStopId := 354

	// Check and convert type
	var stopIdType interface{} = nyctStopId
	var stopId string
	switch v := stopIdType.(type) {
	case int:
		stopId = strconv.Itoa(nyctStopId)
	case string:
		stopId = v
	default:
		return Response{}, fmt.Errorf("%v is not typeof int or string. ", nyctStopId)
	}

	// nyctDirection := "N"
	// trainsToReturn := 2

	stationsResponse, err := stations.LoadStations()
	if err != nil {
		return Response{}, err
	}

	var stationsData stations.Response
	if err := json.Unmarshal(stationsResponse, &stationsData); err != nil {
		return Response{}, err
	}

	localStation, err := stations.LoadLocalStation(stopId, stationsData)
	if err != nil {
		return Response{}, err
	}

	nyctData, err := nyct.CallNYCT()
	if err != nil {
		return Response{}, err
	}

	trips := nyct.ProcessTrips(nyctData, localStation.GTFS_ID, nyctDirection)

	// fmt.Printf("Upcoming trains for %s:\n", localStation.StationName)

	response := Response{
		Trains: []Train{},
	}

	for i, trip := range trips {
		if i >= trainsToReturn {
			break
		}
		timeUntilArrival := time.Until(trip.ArrivalTime)
		// fmt.Printf("Route %s %s arriving in %v\n", trip.RouteID, trip.Direction, timeUntilArrival.Round(time.Minute))

		train := Train{
			RouteID:    trip.RouteID,
			Direction:  trip.Direction,
			ArrivingIn: format.Duration(timeUntilArrival),
		}

		response.Trains = append(response.Trains, train)
	}

	return response, nil
}
