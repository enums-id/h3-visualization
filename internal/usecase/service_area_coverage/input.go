package serviceareacoverage

type Input struct {
	Polygon   [][2]float64 `json:"polygon"`
	InnerHoles [][][2]float64 `json:"inner_holes"`
	Resolution int         `json:"resolution"`
}
