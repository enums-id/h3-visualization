package surgesimulation

type DriverETA struct {
	DriverIndex       int     `json:"driver_index"`
	DriverCell        string  `json:"driver_cell"`
	ToCustomerMinutes float64 `json:"to_customer_minutes"`
	TripMinutes       float64 `json:"trip_minutes"`
	TotalMinutes      float64 `json:"total_minutes"`
}

type Output struct {
	CustomerCell    string      `json:"customer_cell"`
	DestinationCell string      `json:"destination_cell"`
	DriverGrids     []string    `json:"driver_grids"`
	ETAs            []DriverETA `json:"etas"`
	BestDriverIndex int         `json:"best_driver_index"`
	SurgeMultiplier float64     `json:"surge_multiplier"`
}
