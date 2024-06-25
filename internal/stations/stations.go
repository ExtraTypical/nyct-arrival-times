package stations

import (
	"fmt"
	"io"
	"net/http"
)

func LoadStations() ([]byte, error) {

	client := &http.Client{}

	nyctDataUrl := "https://data.ny.gov/resource/39hk-dx4f.json"
	request, err := http.NewRequest("GET", nyctDataUrl, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	response.Header.Add("Accept", "application/json")
	response.Header.Add("Content-Type", "application/json")

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return data, nil
}

func LoadLocalStation(stop_id string, stations Response) (Station, error) {

	for _, stop := range stations {
		if stop_id == stop.GTFS_ID {
			return stop, nil
		}
	}
	return Station{}, fmt.Errorf("no stop found for that stop_id")
}
