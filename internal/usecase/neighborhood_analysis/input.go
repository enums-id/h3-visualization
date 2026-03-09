package neighborhoodanalysis

type Input struct {
	CenterCell string `json:"center_cell"`
	K          int    `json:"k"`
	IncludeCenter bool `json:"include_center"`
}
