package main

import (
	"bytes"
	"io"
	"log"
	"os"

	profilev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"
	"github.com/grafana/pyroscope/pkg/pprof"
)

// Install:
// go install ./cmd/trimstrings

// Clean:
// cat profile.pb.gz | trimstrings > profile_clean.pb.gz

// Validate:
// cat profile.pb.gz | gunzip | strings | tail -n 50
// cat profile_clean.pb.gz | gunzip | strings | tail -n 50

func main() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln("Failed to read input:", err)
	}

	var p profilev1.Profile
	if err = pprof.Unmarshal(data, &p); err != nil {
		log.Fatalln("Failed to parse profile:", err)
	}
	for i := range p.StringTable {
		p.StringTable[i] = ""
	}

	clean, err := pprof.Marshal(&p, true)
	if err != nil {
		log.Fatalln("Failed to marshal:", err)
	}

	if _, err = io.Copy(os.Stdout, bytes.NewBuffer(clean)); err != nil {
		log.Fatalln("Failed to write profile:", err)
	}
}
