package surgesimulation

import (
	"fmt"
	"math"

	h3 "github.com/uber/h3-go/v4"
	"h3-visualization/internal/h3ops"
)

type Service struct {
	h3 h3ops.Client
}

func NewService(h3Client h3ops.Client) *Service {
	return &Service{h3: h3Client}
}

func (s *Service) Execute(in Input) (Output, error) {
	if in.Resolution < 0 || in.Resolution > 15 {
		return Output{}, fmt.Errorf("resolution must be between 0 and 15")
	}
	if len(in.Drivers) == 0 {
		return Output{}, fmt.Errorf("drivers is required")
	}
	if in.AvgSpeedKMH <= 0 {
		return Output{}, fmt.Errorf("avg_speed_kmh must be > 0")
	}
	if in.NeighborhoodK < 0 {
		return Output{}, fmt.Errorf("neighborhood_k must be >= 0")
	}

	customerCell, err := s.h3.LatLngToCell(in.Customer.Lat, in.Customer.Lng, in.Resolution)
	if err != nil {
		return Output{}, err
	}
	destinationCell, err := s.h3.LatLngToCell(in.Destination.Lat, in.Destination.Lng, in.Resolution)
	if err != nil {
		return Output{}, err
	}

	driverGrids := make([]string, 0, len(in.Drivers))
	etas := make([]DriverETA, 0, len(in.Drivers))
	bestIndex := -1
	bestTotal := math.MaxFloat64

	for i, d := range in.Drivers {
		driverCell, err := s.h3.LatLngToCell(d.Lat, d.Lng, in.Resolution)
		if err != nil {
			return Output{}, fmt.Errorf("invalid driver[%d]: %w", i, err)
		}
		driverGrids = append(driverGrids, driverCell.String())

		toCustomerSteps, err := h3.GridDistance(driverCell, customerCell)
		if err != nil {
			return Output{}, fmt.Errorf("cannot compute driver[%d] to customer distance: %w", i, err)
		}
		tripSteps, err := h3.GridDistance(customerCell, destinationCell)
		if err != nil {
			return Output{}, fmt.Errorf("cannot compute customer to destination distance: %w", err)
		}

		// Approximation for learning flow:
		// 1 grid step ~= average edge length at current resolution.
		edgeKm := 1.0
		toCustomerMinutes := (float64(toCustomerSteps) * edgeKm / in.AvgSpeedKMH) * 60
		tripMinutes := (float64(tripSteps) * edgeKm / in.AvgSpeedKMH) * 60
		total := toCustomerMinutes + in.PickupServiceMinutes + tripMinutes

		etas = append(etas, DriverETA{
			DriverIndex:       i,
			DriverCell:        driverCell.String(),
			ToCustomerMinutes: toCustomerMinutes,
			TripMinutes:       tripMinutes,
			TotalMinutes:      total,
		})

		if total < bestTotal {
			bestTotal = total
			bestIndex = i
		}
	}

	surge, err := s.computeSurge(customerCell, driverGrids, in.NeighborhoodK)
	if err != nil {
		return Output{}, err
	}

	return Output{
		CustomerCell:    customerCell.String(),
		DestinationCell: destinationCell.String(),
		DriverGrids:     driverGrids,
		ETAs:            etas,
		BestDriverIndex: bestIndex,
		SurgeMultiplier: surge,
	}, nil
}

func (s *Service) computeSurge(customerCell h3.Cell, driverGrids []string, k int) (float64, error) {
	neighborhood, err := s.h3.GridDisk(customerCell, k)
	if err != nil {
		return 0, err
	}

	available := map[h3.Cell]struct{}{}
	for _, c := range neighborhood {
		available[c] = struct{}{}
	}

	supply := 0
	for _, d := range driverGrids {
		cell := h3.CellFromString(d)
		if _, ok := available[cell]; ok {
			supply++
		}
	}

	switch {
	case supply >= 5:
		return 1.0, nil
	case supply >= 3:
		return 1.2, nil
	case supply >= 1:
		return 1.5, nil
	default:
		return 2.0, nil
	}
}
