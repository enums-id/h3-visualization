package routecells

import (
	"fmt"

	"github.com/uber/h3-go/v4"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(in Input) (Output, error) {
	if in.OriginCell == "" || in.DestinationCell == "" {
		return Output{}, fmt.Errorf("origin_cell and destination_cell are required")
	}

	origin := h3.CellFromString(in.OriginCell)

	if !origin.IsValid() {
		return Output{}, fmt.Errorf("invalid origin_cell: %s", in.OriginCell)
	}

	destination := h3.CellFromString(in.DestinationCell)

	if !destination.IsValid() {
		return Output{}, fmt.Errorf("invalid destination_cell: %s", in.DestinationCell)
	}

	pathCells, err := h3.GridPath(origin, destination)

	if err != nil {
		return Output{}, fmt.Errorf("grid path cells failed: %w", err)
	}

	out := make([]string, 0, len(pathCells))

	for _, c := range pathCells {
		out = append(out, c.String())
	}

	return Output{out}, nil
}
