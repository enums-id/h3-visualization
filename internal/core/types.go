package core

type UseCase string

const (
	UseCasePointIndexing       UseCase = "point_indexing"
	UseCaseNeighborhood        UseCase = "neighborhood_analysis"
	UseCaseServiceAreaCoverage UseCase = "service_area_coverage"
	UseCaseRouteCells          UseCase = "route_cells"
	UseCaseSurgeSimulation     UseCase = "surge_simulation"
)

type RunRequest struct {
	UseCase UseCase `json:"use_case"`
	Payload []byte  `json:"payload"`
}
