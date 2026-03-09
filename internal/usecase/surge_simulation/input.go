package surgesimulation

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Input struct {
	Customer             Coordinate   `json:"customer"`
	Drivers              []Coordinate `json:"drivers"`
	Destination          Coordinate   `json:"destination"`
	Resolution           int          `json:"resolution"`
	AvgSpeedKMH          float64      `json:"avg_speed_kmh"`
	PickupServiceMinutes float64      `json:"pickup_service_minutes"`
	NeighborhoodK        int          `json:"neighborhood_k"`
}
