# Implementation Plan: Bridge Go Engine to HTML (No REST API, No JS Framework)

## Goal

Build a simple server-side rendered (SSR) web app where HTML forms call Go engine modules directly, then render results back on the page.

---

## 1) Target Architecture

```text
cmd/
  engine/main.go                 # existing CLI
  web/main.go                    # new web entrypoint

internal/
  web/
    handler.go                   # HTTP handlers + form parsing + engine invocation
    viewmodel.go                 # UI payload structs (optional but clean)
    parse.go                     # parse helpers for coords, arrays, json
    routes.go                    # route registration (optional split)

templates/
  layout.html                    # base layout
  index.html                     # form + result sections
  partial_result_*.html          # optional partials per use case

static/
  styles.css                     # modern basic CSS
```

---

## 2) UX Scope (Simple but Useful)

Single page `/`:

- Use case selector (`point_indexing`, `neighborhood_analysis`, `service_area_coverage`, `route_cells`, `surge_simulation`)
- Dynamic form sections (show/hide fields with tiny vanilla JS, optional)
- Submit button
- Rendered result card:
  - success JSON (pretty printed)
  - or inline error message
- “Example payload” buttons (auto-fill textarea) to help testing

No SPA behavior required.

---

## 3) Routing Plan

- `GET /`  
  Render initial form page with default use case and example payload.
- `POST /run`  
  Parse submitted fields, build `core.RunRequest`, call `runner.New().Run(...)`, render same page with result.
- `GET /static/styles.css`  
  Serve CSS via `http.FileServer`.

No REST endpoints.

---

## 4) Data Input Strategy

### Option A (fastest, recommended for learning)

Use one payload textarea:

- user chooses use case
- user pastes JSON payload
- server validates JSON and runs engine

### Option B (better UX, more code)

Provide individual HTML input fields per use case and compose payload server-side.

Recommended now: **A first**, later upgrade to B.

---

## 5) Templating Strategy

Use `html/template`:

- `layout.html` (header, main container, footer)
- `index.html` includes:
  - `<form method="POST" action="/run">`
  - usecase `<select>`
  - payload `<textarea>`
  - submit button
  - result `<pre>` block (if present)
  - error alert block (if present)

Template data struct:

```go
type PageData struct {
    Title          string
    UseCase        string
    Payload        string
    ResultJSON     string
    ErrorMessage   string
    ExamplePayload map[string]string
}
```

---

## 6) Web Handler Responsibilities

1. Read `usecase`, `payload` from form.
2. Validate non-empty + valid JSON.
3. Call engine:
   - `r := runner.New()`
   - `out, err := r.Run(core.RunRequest{...})`
4. Marshal output to indented JSON for display.
5. Render same template with `PageData`.

---

## 7) Styling Plan (Modern Basic CSS)

`static/styles.css`:

- centered container with max width
- card-based layout
- clean typography, muted background
- input/select/textarea modern spacing + borders
- responsive button states
- error (red tinted) / success (green tinted) panels
- code/pre block styling for JSON output

No CSS framework dependency.

---

## 8) Files to Create

1. `cmd/web/main.go`
   - setup mux, static files, handler registration, listen on `:8080`

2. `internal/web/handler.go`
   - `IndexHandler`, `RunHandler`, template parse/render helpers

3. `templates/layout.html`
4. `templates/index.html`
5. `static/styles.css`

Optional:

- `internal/web/examples.go` for use case sample payload strings

---

## 9) Example Payload Set (preloaded in UI)

- point_indexing

```json
{ "lat": -6.2, "lng": 106.816666, "resolution": 9 }
```

- neighborhood_analysis

```json
{ "center_cell": "8928308280fffff", "k": 1, "include_center": false }
```

- service_area_coverage

```json
{
  "polygon": [
    [-6.2, 106.81],
    [-6.19, 106.82],
    [-6.21, 106.83]
  ],
  "inner_holes": [],
  "resolution": 9
}
```

- route_cells

```json
{ "origin_cell": "8928308280fffff", "destination_cell": "8928308280bffff" }
```

- surge_simulation

```json
{
  "customer": { "lat": -6.2, "lng": 106.816666 },
  "drivers": [
    { "lat": -6.21, "lng": 106.81 },
    { "lat": -6.199, "lng": 106.82 }
  ],
  "destination": { "lat": -6.18, "lng": 106.84 },
  "resolution": 9,
  "avg_speed_kmh": 30,
  "pickup_service_minutes": 2,
  "neighborhood_k": 1
}
```

---

## 10) Validation Plan

1. `go test ./...`
2. `go run ./cmd/web`
3. manual browser test:
   - open `http://localhost:8080`
   - run each use case with provided examples
   - verify output/error rendering
4. check invalid JSON handling and unknown use case handling.

---

## 11) Incremental Delivery Steps (execution order)

1. Create web entrypoint + routes.
2. Add templates and CSS.
3. Implement `GET /` default rendering.
4. Implement `POST /run` engine invocation.
5. Add examples map + default payload switching. For best UX, don't let user input the JSON directly. Instead, use a picker/input/select for non geographic input and use map pin on top of map for geographic input. Dont forget to serve default values so user can directly run if have no idea about the data
6. Run tests and fix compile issues.
7. Manual smoke test in browser.
8. (Optional) add minimal JS for auto-fill payload when use case changes.

---

## 12) Risks / Notes

- `route_cells` and any partially implemented use cases may still panic depending on current module status; handler should recover and show safe error.
- For learning, keep logic in server handlers explicit and simple (avoid premature abstractions).
- Avoid adding persistence/database for now.

---

## 13) Future Enhancements (after MVP)

- Add H3 boundary coordinate output for direct polygon rendering
- Add lightweight map rendering (Leaflet + no framework)
- Add benchmark page for use case performance snapshots
