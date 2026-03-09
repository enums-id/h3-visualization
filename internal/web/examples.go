package web

// DefaultFormValues holds example inputs for each use case, keyed by HTML form field name.
var DefaultFormValues = map[string]map[string]string{
	"point_indexing": {
		"pi_lat":        "-6.2",
		"pi_lng":        "106.816666",
		"pi_resolution": "9",
	},
	"neighborhood_analysis": {
		"na_center_cell":    "8928308280fffff",
		"na_k":              "1",
		"na_include_center": "",
	},
	"service_area_coverage": {
		"sac_polygon_json": `[[-6.2,106.81],[-6.19,106.82],[-6.21,106.83]]`,
		"sac_resolution":   "9",
	},
	"route_cells": {
		"rc_origin_cell":      "8928308280fffff",
		"rc_destination_cell": "8928308280bffff",
	},
	"surge_simulation": {
		"ss_customer_lat":           "-6.2",
		"ss_customer_lng":           "106.816666",
		"ss_drivers_json":           `[{"lat":-6.21,"lng":106.81},{"lat":-6.199,"lng":106.82}]`,
		"ss_destination_lat":        "-6.18",
		"ss_destination_lng":        "106.84",
		"ss_resolution":             "9",
		"ss_avg_speed_kmh":          "30",
		"ss_pickup_service_minutes": "2",
		"ss_neighborhood_k":         "1",
	},
}
