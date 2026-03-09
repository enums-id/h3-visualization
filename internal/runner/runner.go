package runner

import (
	"encoding/json"
	"fmt"

	"h3-visualization/internal/core"
	"h3-visualization/internal/h3ops"
	neighborhoodanalysis "h3-visualization/internal/usecase/neighborhood_analysis"
	pointindexing "h3-visualization/internal/usecase/point_indexing"
	routecells "h3-visualization/internal/usecase/route_cells"
	serviceareacoverage "h3-visualization/internal/usecase/service_area_coverage"
	surgesimulation "h3-visualization/internal/usecase/surge_simulation"
)

type Runner struct {
	h3 *h3ops.H3Client
}

func New() *Runner {
	return &Runner{h3: h3ops.New()}
}

func (r *Runner) Run(req core.RunRequest) (any, error) {
	switch req.UseCase {
	case core.UseCasePointIndexing:
		var in pointindexing.Input
		if err := json.Unmarshal(req.Payload, &in); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		return pointindexing.NewService(r.h3).Execute(in)
	case core.UseCaseNeighborhood:
		var in neighborhoodanalysis.Input
		if err := json.Unmarshal(req.Payload, &in); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		return neighborhoodanalysis.NewService(r.h3).Execute(in)
	case core.UseCaseServiceAreaCoverage:
		var in serviceareacoverage.Input
		if err := json.Unmarshal(req.Payload, &in); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		return serviceareacoverage.NewService().Execute(in)
	case core.UseCaseRouteCells:
		var in routecells.Input
		if err := json.Unmarshal(req.Payload, &in); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		return routecells.NewService().Execute(in)
	case core.UseCaseSurgeSimulation:
		var in surgesimulation.Input
		if err := json.Unmarshal(req.Payload, &in); err != nil {
			return nil, fmt.Errorf("invalid payload: %w", err)
		}
		return surgesimulation.NewService(r.h3).Execute(in)
	default:
		return nil, core.ErrUnknownUseCase
	}
}
