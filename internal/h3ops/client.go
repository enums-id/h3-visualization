package h3ops

import h3 "github.com/uber/h3-go/v4"

type Client interface {
	LatLngToCell(lat, lng float64, res int) (h3.Cell, error)
	GridDisk(cell h3.Cell, k int) ([]h3.Cell, error)
}
