package stations

type Response []Station

type Station struct {
	GTFS_ID     string `json:"gtfs_stop_id"`
	StationID   string `json:"station_id"`
	StationName string `json:"stop_name"`
}
