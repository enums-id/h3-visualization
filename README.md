# H3 Visualization

An interactive Go application for exploring [Uber's H3](https://h3geo.org/) hexagonal geospatial indexing system. It provides both a CLI engine and a web UI with live map visualization.

## Features

Five H3 use cases, each runnable via CLI or the web interface:

| Use Case | H3 Operation | Description |
|---|---|---|
| **Point Indexing** | `LatLngToCell` | Convert a GPS coordinate to an H3 cell at a given resolution |
| **Neighborhood Analysis** | `GridDisk` | Find all cells within k rings of a center cell |
| **Service Area Coverage** | `PolygonToCells` | Fill a polygon with H3 cells |
| **Route Cells** | `GridPathCells` | Find cells along the grid path between two cells |
| **Surge Simulation** | Composite | Simulate ride-hailing surge: ETA estimation, best driver selection, surge multiplier |

## Project Structure

```
.
├── cmd/
│   ├── engine/        # CLI entry point
│   └── web/           # Web server entry point
├── internal/
│   ├── core/          # Shared types and errors
│   ├── h3ops/         # H3 operations client (interface + implementation)
│   ├── runner/        # Use case dispatcher
│   ├── usecase/       # Individual use case services and tests
│   │   ├── point_indexing/
│   │   ├── neighborhood_analysis/
│   │   ├── service_area_coverage/
│   │   ├── route_cells/
│   │   └── surge_simulation/
│   └── web/           # HTTP handlers, view models, examples
├── templates/         # HTML templates (layout + index)
├── static/            # Static assets (CSS, JS)
├── go.mod
└── USECASES.md        # Detailed use case documentation
```

## Requirements

- Go 1.23.5+

## Running

### Web UI

```bash
go run ./cmd/web
```

Opens at [http://localhost:8080](http://localhost:8080). The interactive UI lets you place points and polygons on a map, run each use case, and view results as JSON or rendered H3 cells on the result map.

> **Note:** Run this command from the repository root directory so the server can locate `templates/` and `static/`.

### CLI Engine

```bash
go run ./cmd/engine -usecase <name> -payload '<json>'
```

#### Examples

**Point Indexing**
```bash
go run ./cmd/engine -usecase point_indexing \
  -payload '{"lat":-6.200000,"lng":106.816666,"resolution":9}'
```

**Neighborhood Analysis**
```bash
go run ./cmd/engine -usecase neighborhood_analysis \
  -payload '{"center_cell":"8928308280fffff","k":1,"include_center":true}'
```

**Service Area Coverage**
```bash
go run ./cmd/engine -usecase service_area_coverage \
  -payload '{"polygon":[[-6.20,106.81],[-6.19,106.82],[-6.21,106.83]],"inner_holes":[],"resolution":9}'
```

**Route Cells**
```bash
go run ./cmd/engine -usecase route_cells \
  -payload '{"origin_cell":"8928308280fffff","destination_cell":"8928308280bffff"}'
```

**Surge Simulation**
```bash
go run ./cmd/engine -usecase surge_simulation \
  -payload '{"customer":{"lat":-6.2,"lng":106.816666},"drivers":[{"lat":-6.21,"lng":106.81},{"lat":-6.199,"lng":106.82}],"destination":{"lat":-6.18,"lng":106.84},"resolution":9,"avg_speed_kmh":30,"pickup_service_minutes":2,"neighborhood_k":1}'
```

## Testing

```bash
go test ./...
```

## Dependencies

- [uber/h3-go v4](https://github.com/uber/h3-go) — Go bindings for the H3 library
- [Leaflet.js](https://leafletjs.com/) — Map rendering in the web UI
- [h3-js](https://github.com/uber/h3-js) — Client-side H3 cell boundary rendering
- [OpenStreetMap](https://www.openstreetmap.org/) — Map tiles
