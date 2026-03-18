package web

import "os"

type PageData struct {
	Title        string
	UseCase      string
	ResultJSON   string
	ErrorMessage string
	FormValues   map[string]string
	Resolutions  []int
	UseCases     []UseCaseInfo
	MapTilerKey  string
}

type UseCaseInfo struct {
	ID    string
	Label string
}

var AllUseCases = []UseCaseInfo{
	{ID: "point_indexing", Label: "Point Indexing"},
	{ID: "neighborhood_analysis", Label: "Neighborhood Analysis"},
	{ID: "service_area_coverage", Label: "Service Area Coverage"},
	{ID: "route_cells", Label: "Route Cells"},
	{ID: "surge_simulation", Label: "Surge Simulation"},
}

var allResolutions []int

func init() {
	for i := 0; i <= 15; i++ {
		allResolutions = append(allResolutions, i)
	}
}

// newPageData creates a PageData seeded with defaults for every use case,
// so hidden form sections always carry their default values on first render.
func newPageData(useCase string) PageData {
	merged := make(map[string]string)
	for _, defaults := range DefaultFormValues {
		for k, v := range defaults {
			merged[k] = v
		}
	}
	return PageData{
		Title:       "H3 Visualization",
		UseCase:     useCase,
		FormValues:  merged,
		Resolutions: allResolutions,
		UseCases:    AllUseCases,
		MapTilerKey: os.Getenv("MAPTILER_KEY"),
	}
}
