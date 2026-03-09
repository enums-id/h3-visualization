# H3 Engine Use Cases

This project focuses on backend H3 engine capabilities before HTML visualization.

## 1) `point_indexing`

### Real-world case
Assigning a courier request (single GPS point) into one H3 cell bucket for fast regional lookup.

### What it does
Converts `lat/lng + resolution` into a single H3 cell.

### Example command
```bash
go run ./cmd/engine -usecase point_indexing -payload '{"lat":-6.200000,"lng":106.816666,"resolution":9}'
```

### Expected output shape
```json
{
  "cell": "<h3-cell-id>"
}
```

### Expected result notes
- `cell` is non-empty when input is valid.
- Higher resolution means smaller cell area.

---

## 2) `neighborhood_analysis`

### Real-world case
Finding nearby service coverage cells around a selected operational zone center.

### What it does
Returns cells within H3 grid distance `k` from a center cell (`GridDisk`).

### Example command (include center)
```bash
go run ./cmd/engine -usecase neighborhood_analysis -payload '{"center_cell":"8928308280fffff","k":1,"include_center":true}'
```

### Example command (exclude center)
```bash
go run ./cmd/engine -usecase neighborhood_analysis -payload '{"center_cell":"8928308280fffff","k":1,"include_center":false}'
```

### Expected output shape
```json
{
  "cells": ["<h3-cell-id>", "<h3-cell-id>"]
}
```

### Expected result notes
- `k=0` returns only center (unless excluded).
- `include_center=false` removes the center cell from the result.

---

## 3) `service_area_coverage`

### Real-world case
Modeling a delivery service area polygon (with optional restricted subzones/holes) into H3 cells.

### What it does
Converts polygon boundaries into a cell set covering the area (`PolygonToCells`).

### Example command
```bash
go run ./cmd/engine -usecase service_area_coverage -payload '{"polygon":[[-6.20,106.81],[-6.19,106.82],[-6.21,106.83]],"inner_holes":[],"resolution":9}'
```

### Example command (with one hole)
```bash
go run ./cmd/engine -usecase service_area_coverage -payload '{"polygon":[[-6.210,106.800],[-6.180,106.820],[-6.220,106.850]],"inner_holes":[[[-6.202,106.815],[-6.198,106.820],[-6.204,106.826]]],"resolution":9}'
```

### Expected output shape
```json
{
  "cells": ["<h3-cell-id>", "<h3-cell-id>"]
}
```

### Expected result notes
- Output should be non-empty for a valid polygon.
- Holes remove cells inside excluded sub-polygons.

---

## 4) `route_cells`

### Real-world case
Estimating hex-grid traversal between origin and destination zones for route simulation.

### What it does
Computes a cell path between two H3 cells (`GridPathCells`).

### Example command
```bash
go run ./cmd/engine -usecase route_cells -payload '{"origin_cell":"8928308280fffff","destination_cell":"8928308280bffff"}'
```

### Expected output shape
```json
{
  "path": ["<h3-cell-id>", "<h3-cell-id>"]
}
```

### Expected result notes
- First item should be origin cell.
- Last item should be destination cell.
- Path may fail for invalid/non-compatible cells.

---

## 5) `surge_simulation`

### Real-world case
Customer requests a ride, multiple drivers exist nearby, and the system estimates ETA and surge level from local supply.

### What it does
- Maps customer, drivers, and destination into H3 grids.
- Estimates total time for each driver:
  - driver -> customer
  - pickup service minutes
  - customer -> destination
- Selects the best (lowest total) driver.
- Computes a simple surge multiplier from nearby available drivers in customer `k` neighborhood.

### Example command
```bash
go run ./cmd/engine -usecase surge_simulation -payload '{"customer":{"lat":-6.2,"lng":106.816666},"drivers":[{"lat":-6.21,"lng":106.81},{"lat":-6.199,"lng":106.82},{"lat":-6.205,"lng":106.814}],"destination":{"lat":-6.18,"lng":106.84},"resolution":9,"avg_speed_kmh":30,"pickup_service_minutes":2,"neighborhood_k":1}'
```

### Expected output shape
```json
{
  "customer_cell": "<h3-cell-id>",
  "destination_cell": "<h3-cell-id>",
  "driver_grids": ["<h3-cell-id>"],
  "etas": [
    {
      "driver_index": 0,
      "driver_cell": "<h3-cell-id>",
      "to_customer_minutes": 1.2,
      "trip_minutes": 8.1,
      "total_minutes": 11.3
    }
  ],
  "best_driver_index": 0,
  "surge_multiplier": 1.2
}
```

### Expected result notes
- `best_driver_index` points to the minimum `total_minutes` driver.
- `surge_multiplier` is higher when fewer drivers are available near customer.

---

## Validation commands

```bash
go test ./internal/usecase/service_area_coverage -v
go test ./...
```
