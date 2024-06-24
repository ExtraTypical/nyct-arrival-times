package nyct

import (
	"io"
	"log"
	"net/http"
	"sort"
	"testapi/internal/protobuf"
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

func ProcessTrips(feedMessage *protobuf.FeedMessage, baseStationID string) []TripInfo {
	var trips []TripInfo
	now := time.Now()

	northID := baseStationID + "N"
	southID := baseStationID + "S"

	for _, entity := range feedMessage.GetEntity() {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate == nil {
			continue
		}

		for _, stopTimeUpdate := range tripUpdate.GetStopTimeUpdate() {
			stopID := stopTimeUpdate.GetStopId()
			if stopID == northID || stopID == southID {
				arrivalTime := time.Unix(stopTimeUpdate.GetArrival().GetTime(), 0)
				if arrivalTime.After(now) {
					direction := "Northbound"
					if stopID == southID {
						direction = "Southbound"
					}
					trips = append(trips, TripInfo{
						RouteID:     tripUpdate.GetTrip().GetRouteId(),
						TripID:      tripUpdate.GetTrip().GetTripId(),
						ArrivalTime: arrivalTime,
						Direction:   direction,
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
