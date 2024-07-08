package nyctapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/ExtraTypical/nyct-arrival-times/internal/protobuf"

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

	filterByDirection := direction == "N" || direction == "S" || direction == "W" || direction == "E"

	for _, entity := range feedMessage.GetEntity() {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate == nil {
			continue
		}

		for _, stopTimeUpdate := range tripUpdate.GetStopTimeUpdate() {
			stopID := stopTimeUpdate.GetStopId()
			stopIDBase := strings.TrimRight(stopID, "NSWE")

			if stopIDBase == baseStationID {

				stopDirection := stopID[len(stopID)-1:]
				isTargetDirection := !filterByDirection || stopDirection == direction

				if isTargetDirection {
					arrivalTime := time.Unix(stopTimeUpdate.GetArrival().GetTime(), 0)

					if arrivalTime.After(now) {
						tripID := tripUpdate.GetTrip().GetTripId()

						lastStop, err := getLastStop(tripID, feedMessage)
						if err != nil {
							log.Printf("Warning: Last stop not found for trip %s", tripID)
							lastStop = "Unknown"
						}

						trip := TripInfo{
							RouteID:     tripUpdate.GetTrip().GetRouteId(),
							TripID:      tripID,
							ArrivalTime: arrivalTime,
							Direction:   getDirectionFromStopID(stopID),
							LastStop:    lastStop,
						}
						trips = append(trips, trip)
					}
				}
				break
			}
		}
	}

	sort.Slice(trips, func(i, j int) bool {
		return trips[i].ArrivalTime.Before(trips[j].ArrivalTime)
	})

	return trips
}

func getDirectionFromStopID(stopID string) string {
	directionSuffix := stopID[len(stopID)-1:]
	switch directionSuffix {
	case "N":
		return "Northbound"
	case "S":
		return "Southbound"
	default:
		return "Unknown"
	}
}

func getLastStop(tripID string, feedMessage *protobuf.FeedMessage) (string, error) {
	for _, entity := range feedMessage.GetEntity() {
		tripUpdate := entity.GetTripUpdate()
		if tripUpdate == nil {
			continue
		}

		if tripUpdate.GetTrip().GetTripId() == tripID {
			stopTimeUpdate := tripUpdate.GetStopTimeUpdate()
			if len(stopTimeUpdate) == 0 {
				return "", fmt.Errorf("no stop time updates found for train %s", tripID)
			}

			lastStop := stopTimeUpdate[len(stopTimeUpdate)-1]
			return lastStop.GetStopId(), nil
		}
	}
	return "", fmt.Errorf("train with ID %s not found", tripID)
}
