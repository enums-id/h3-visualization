package pointindexing

import "h3-visualization/internal/h3ops"

type Service struct {
	h3 h3ops.Client
}

func NewService(h3Client h3ops.Client) *Service {
	return &Service{h3: h3Client}
}

func (s *Service) Execute(in Input) (Output, error) {
	cell, err := s.h3.LatLngToCell(in.Lat, in.Lng, in.Resolution)
	if err != nil {
		return Output{}, err
	}
	return Output{Cell: cell.String()}, nil
}
