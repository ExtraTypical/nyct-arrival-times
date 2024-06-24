package nyct

import "time"

type TripInfo struct {
	RouteID     string
	TripID      string
	ArrivalTime time.Time
	Direction   string
}
