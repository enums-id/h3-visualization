package serviceareacoverage

import (
	"fmt"

	h3 "github.com/uber/h3-go/v4"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(in Input) (Output, error) {
	if len(in.Polygon) < 3 {
		return Output{}, fmt.Errorf("polygon requires at least 3 points")
	}

	if in.Resolution < 0 || in.Resolution > 15 {
		return Output{}, fmt.Errorf("resolution must be between 0 and 15")
	}

	// Create H3 LatLng vertices for the outer polygon
	outer, err := toLatLngLoop(in.Polygon, "outer polygon")
	if err != nil {
		return Output{}, fmt.Errorf("invalid outer polygon: %w", err)
	}

	// Create the GeoPolygon with holes
	geoPolygon := h3.GeoPolygon{
        GeoLoop: outer,
        Holes:   make([]h3.GeoLoop, 0, len(in.InnerHoles)),
     }

	 // Add inner holes to the GeoPolygon
     for i, hole := range in.InnerHoles {
		holes, err := toLatLngLoop(hole, fmt.Sprintf("inner hole %d", i))
		if err != nil {
			return Output{}, fmt.Errorf("invalid inner holes: %w", err)
		}
        geoPolygon.Holes = append(geoPolygon.Holes, h3.GeoLoop(holes))
     }

	 // Get H3 cells covering the polygon
     cells, err := h3.PolygonToCells(geoPolygon, in.Resolution)
     if err != nil {
        return Output{}, fmt.Errorf("polygon to cells failed: %w", err)
     }

	 // Convert H3 cells to strings
	out := make([]string, 0, len(cells))
     for _, c := range cells {
        out = append(out, c.String())
     }

     return Output{Cells: out}, nil
}

func toLatLngLoop(points [][2]float64, label string) ([]h3.LatLng, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("%s requires at least 3 points", label)
	}

	loop := make([]h3.LatLng, 0, len(points))
	for i, p := range points {
		lat := p[0]
		lng := p[1]

		if lat < -90 || lat > 90 {
				return nil, fmt.Errorf("invalid latitude at %s[%d]: %f", label, i, lat)
		}
		if lng < -180 || lng > 180 {
				return nil, fmt.Errorf("invalid longitude at %s[%d]: %f", label, i, lng)
		}

		loop = append(loop, h3.LatLng{Lat: lat, Lng: lng})
	}

	return loop, nil
}
