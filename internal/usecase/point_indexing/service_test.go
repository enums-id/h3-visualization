package pointindexing

import (
	"testing"

	"h3-visualization/internal/h3ops"
)

func TestExecute(t *testing.T) {
	svc := NewService(h3ops.New())
	out, err := svc.Execute(Input{Lat: -6.2, Lng: 106.816666, Resolution: 9})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Cell == "" {
		t.Fatalf("expected non-empty cell")
	}
}
