package surgesimulation

import (
	"testing"

	"h3-visualization/internal/h3ops"
)

func TestExecute(t *testing.T) {
	svc := NewService(h3ops.New())
	out, err := svc.Execute(Input{
		Customer:             Coordinate{Lat: -6.2, Lng: 106.816666},
		Drivers:              []Coordinate{{Lat: -6.21, Lng: 106.81}, {Lat: -6.199, Lng: 106.82}},
		Destination:          Coordinate{Lat: -6.18, Lng: 106.84},
		Resolution:           9,
		AvgSpeedKMH:          30,
		PickupServiceMinutes: 2,
		NeighborhoodK:        1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.CustomerCell == "" || out.DestinationCell == "" {
		t.Fatalf("expected customer and destination cells")
	}
	if len(out.DriverGrids) != 2 || len(out.ETAs) != 2 {
		t.Fatalf("unexpected drivers/etas length")
	}
	if out.BestDriverIndex < 0 {
		t.Fatalf("expected best driver index")
	}
}
