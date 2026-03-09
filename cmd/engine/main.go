package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"h3-visualization/internal/core"
	"h3-visualization/internal/runner"
)

func main() {
	useCase := flag.String("usecase", "point_indexing", "use case name")
	payload := flag.String("payload", `{}`, "json payload")
	flag.Parse()

	r := runner.New()
	out, err := r.Run(core.RunRequest{
		UseCase: core.UseCase(*useCase),
		Payload: []byte(*payload),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}
