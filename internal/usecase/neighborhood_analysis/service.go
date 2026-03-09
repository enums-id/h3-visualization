package neighborhoodanalysis

import (
	"fmt"

	"h3-visualization/internal/h3ops"

	h3 "github.com/uber/h3-go/v4"
)

type Service struct {
	h3 h3ops.Client
}

func NewService(h3Client h3ops.Client) *Service {
	return &Service{h3: h3Client}
}

func (s *Service) Execute(in Input) (Output, error) {
	if in.CenterCell == "" {
		return Output{}, fmt.Errorf("center_cell is required")
	}

	if in.K < 0 {
		return Output{}, fmt.Errorf("k must be non-negative")
	}

	center := h3.CellFromString(in.CenterCell)
	if center.IsValid() && center == 0 {
		return Output{}, fmt.Errorf("invalid center_cell")
	}

	cells, err := s.h3.GridDisk(center, in.K)
	if err != nil {
		return Output{}, err
	}

	result := make([]string, 0, len(cells))
	for _, c := range cells {
		if !in.IncludeCenter && c == center {
			continue
		}
		result = append(result, c.String())
	}

	return Output{Cells: result}, nil
}
