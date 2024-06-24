package nyctapi

import (
	"go-nyct/internal/protobuf"
	"io"
	"log"
	"net/http"
	"sort"
	"time"

	"google.golang.org/protobuf/proto"
)

func CallNYCT() (*protobuf.FeedMessage, error) {
	client := &http.Client{}

	NYCTUrl := "https://api-endpoint.mta.info/Dataservice/mtagtfsfeeds/nyct%2Fgtfs"
	request, err := http.NewRequest("GET", NYCTUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var nyctData protobuf.FeedMessage
	if err := proto.Unmarshal(data, &nyctData); err != nil {
		log.Fatal(err)
	}

	return &nyctData, nil
}

func ProcessTrips(feedMessage *protobuf.FeedMessage, baseStationID string, direction string) []TripInfo {
	var trips []TripInfo
	now := time.Now()

	var targetID string
	var targetDirection string

	switch direction {
	case "N":
		targetID = baseStationID + "N"
		targetDirection = "Northbound"
	case "S":
		targetID = baseStationID + "S"
		targetDirection = "Southbound"
	default:
		return trips // Return empty slice if invalid direction
	}

	for _, entity := range feedMessage.GetEntity() {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate == nil {
			continue
		}

		for _, stopTimeUpdate := range tripUpdate.GetStopTimeUpdate() {
			stopID := stopTimeUpdate.GetStopId()
			if stopID == targetID {
				arrivalTime := time.Unix(stopTimeUpdate.GetArrival().GetTime(), 0)
				if arrivalTime.After(now) {
					trips = append(trips, TripInfo{
						RouteID:     tripUpdate.GetTrip().GetRouteId(),
						TripID:      tripUpdate.GetTrip().GetTripId(),
						ArrivalTime: arrivalTime,
						Direction:   targetDirection,
					})
				}
				break // We found the update for this station, no need to check further stops
			}
		}
	}

	// Sort trips by arrival time
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].ArrivalTime.Before(trips[j].ArrivalTime)
	})

	return trips
}
