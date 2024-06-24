package nyct

import (
	"encoding/json"
	"time"

	format "github.com/ExtraTypical/nyct-arrival-times/internal/formatduration"
	nyctapi "github.com/ExtraTypical/nyct-arrival-times/internal/nyct"
	"github.com/ExtraTypical/nyct-arrival-times/internal/stations"
)

type Response struct {
	Trains []Train
}

type Train struct {
	RouteID    string
	Direction  string
	ArrivingIn string
}

func CheckArrivalTimes(nyctStopId string, nyctDirection string, trainsToReturn int) (Response, error) {

	stationsResponse, err := stations.LoadStations()
	if err != nil {
		return Response{}, err
	}

	var stationsData stations.Response
	if err := json.Unmarshal(stationsResponse, &stationsData); err != nil {
		return Response{}, err
	}

	localStation, err := stations.LoadLocalStation(nyctStopId, stationsData)
	if err != nil {
		return Response{}, err
	}

	nyctData, err := nyctapi.CallNYCT()
	if err != nil {
		return Response{}, err
	}

	trips := nyctapi.ProcessTrips(nyctData, localStation.GTFS_ID, nyctDirection)

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
