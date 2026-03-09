package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"h3-visualization/internal/core"
	"h3-visualization/internal/runner"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").ParseFiles(
		"templates/layout.html",
		"templates/index.html",
	))
}

func renderPage(w http.ResponseWriter, data PageData) {
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// IndexHandler renders the initial form page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	useCase := r.URL.Query().Get("usecase")
	if useCase == "" {
		useCase = "point_indexing"
	}
	renderPage(w, newPageData(useCase))
}

// RunHandler parses the submitted form, builds an engine payload, runs the engine, and re-renders the page with results.
func RunHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic in RunHandler: %v", rec)
			data := newPageData("point_indexing")
			data.ErrorMessage = fmt.Sprintf("internal error: %v", rec)
			renderPage(w, data)
		}
	}()

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		r.ParseForm() //nolint - fallback for non-multipart
	}

	useCase := r.FormValue("use_case")
	if useCase == "" {
		useCase = "point_indexing"
	}

	// Collect all submitted values for form re-population.
	formValues := make(map[string]string)
	for k, v := range r.Form {
		if len(v) > 0 {
			formValues[k] = v[0]
		}
	}

	payload, err := buildPayload(useCase, r)

	data := newPageData(useCase)
	// Overlay submitted values over defaults so every section retains values.
	for k, v := range formValues {
		data.FormValues[k] = v
	}

	if err != nil {
		data.ErrorMessage = "Payload error: " + err.Error()
		renderPage(w, data)
		return
	}

	out, runErr := runner.New().Run(core.RunRequest{
		UseCase: core.UseCase(useCase),
		Payload: payload,
	})
	if runErr != nil {
		data.ErrorMessage = runErr.Error()
	} else {
		b, _ := json.MarshalIndent(out, "", "  ")
		data.ResultJSON = string(b)
	}

	renderPage(w, data)
}

// buildPayload converts use-case-specific form fields into a JSON payload for the engine.
func buildPayload(useCase string, r *http.Request) ([]byte, error) {
	switch core.UseCase(useCase) {

	case core.UseCasePointIndexing:
		lat, err := strconv.ParseFloat(r.FormValue("pi_lat"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %w", err)
		}
		lng, err := strconv.ParseFloat(r.FormValue("pi_lng"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %w", err)
		}
		res, err := strconv.Atoi(r.FormValue("pi_resolution"))
		if err != nil {
			return nil, fmt.Errorf("invalid resolution: %w", err)
		}
		return json.Marshal(map[string]any{
			"lat": lat, "lng": lng, "resolution": res,
		})

	case core.UseCaseNeighborhood:
		k, err := strconv.Atoi(r.FormValue("na_k"))
		if err != nil {
			return nil, fmt.Errorf("invalid k: %w", err)
		}
		includeCenter := r.FormValue("na_include_center") == "true"
		return json.Marshal(map[string]any{
			"center_cell":   r.FormValue("na_center_cell"),
			"k":             k,
			"include_center": includeCenter,
		})

	case core.UseCaseServiceAreaCoverage:
		polygonJSON := r.FormValue("sac_polygon_json")
		if polygonJSON == "" {
			return nil, fmt.Errorf("polygon is empty – click the map to add at least 3 points")
		}
		var polygon [][2]float64
		if err := json.Unmarshal([]byte(polygonJSON), &polygon); err != nil {
			return nil, fmt.Errorf("invalid polygon data: %w", err)
		}
		if len(polygon) < 3 {
			return nil, fmt.Errorf("polygon requires at least 3 points, got %d", len(polygon))
		}
		res, err := strconv.Atoi(r.FormValue("sac_resolution"))
		if err != nil {
			return nil, fmt.Errorf("invalid resolution: %w", err)
		}
		return json.Marshal(map[string]any{
			"polygon":     polygon,
			"inner_holes": [][][2]float64{},
			"resolution":  res,
		})

	case core.UseCaseRouteCells:
		origin := r.FormValue("rc_origin_cell")
		dest := r.FormValue("rc_destination_cell")
		if origin == "" || dest == "" {
			return nil, fmt.Errorf("both origin and destination cells are required")
		}
		return json.Marshal(map[string]any{
			"origin_cell":      origin,
			"destination_cell": dest,
		})

	case core.UseCaseSurgeSimulation:
		custLat, err := strconv.ParseFloat(r.FormValue("ss_customer_lat"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid customer latitude: %w", err)
		}
		custLng, err := strconv.ParseFloat(r.FormValue("ss_customer_lng"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid customer longitude: %w", err)
		}
		destLat, err := strconv.ParseFloat(r.FormValue("ss_destination_lat"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid destination latitude: %w", err)
		}
		destLng, err := strconv.ParseFloat(r.FormValue("ss_destination_lng"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid destination longitude: %w", err)
		}
		driversJSON := r.FormValue("ss_drivers_json")
		if driversJSON == "" {
			return nil, fmt.Errorf("at least one driver is required")
		}
		var drivers []map[string]float64
		if err := json.Unmarshal([]byte(driversJSON), &drivers); err != nil {
			return nil, fmt.Errorf("invalid drivers data: %w", err)
		}
		res, err := strconv.Atoi(r.FormValue("ss_resolution"))
		if err != nil {
			return nil, fmt.Errorf("invalid resolution: %w", err)
		}
		avgSpeed, err := strconv.ParseFloat(r.FormValue("ss_avg_speed_kmh"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid average speed: %w", err)
		}
		pickupMin, err := strconv.ParseFloat(r.FormValue("ss_pickup_service_minutes"), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid pickup service minutes: %w", err)
		}
		nk, err := strconv.Atoi(r.FormValue("ss_neighborhood_k"))
		if err != nil {
			return nil, fmt.Errorf("invalid neighborhood k: %w", err)
		}
		return json.Marshal(map[string]any{
			"customer":               map[string]float64{"lat": custLat, "lng": custLng},
			"drivers":                drivers,
			"destination":            map[string]float64{"lat": destLat, "lng": destLng},
			"resolution":             res,
			"avg_speed_kmh":          avgSpeed,
			"pickup_service_minutes": pickupMin,
			"neighborhood_k":         nk,
		})

	default:
		return nil, core.ErrUnknownUseCase
	}
}
