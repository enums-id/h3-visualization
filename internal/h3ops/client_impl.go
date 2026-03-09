package h3ops

import (
	"fmt"

	h3 "github.com/uber/h3-go/v4"
)

type H3Client struct{}

func New() *H3Client {
	return &H3Client{}
}

func (c *H3Client) LatLngToCell(lat, lng float64, res int) (h3.Cell, error) {
	if res < 0 || res > 15 {
		return 0, fmt.Errorf("resolution must be between 0 and 15")
	}
	return h3.LatLngToCell(h3.NewLatLng(lat, lng), res)
}

func (c *H3Client) GridDisk(cell h3.Cell, k int) ([]h3.Cell, error) {
	if k < 0 {
		return nil, fmt.Errorf("k must be >= 0")
	}
	return h3.GridDisk(cell, k)
}
