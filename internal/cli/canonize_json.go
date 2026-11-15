package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
)

// cmdCanonizeJSONSafe envolve cmdCanonize garantindo que, quando --json é passado,
// stdout receba apenas o JSON puro sem texto extra de ganchos ou avisos.
func cmdCanonizeJSONSafe(ctx context.Context, args []string, d Deps) int {
	if !hasJSONFlag(args) {
		return cmdCanonize(ctx, args, d)
	}
	var rawOut bytes.Buffer
	dCap := d
	dCap.Stdout = &rawOut
	code := cmdCanonize(ctx, args, dCap)
	emitPureJSONOrFallback(&rawOut, d)
	return code
}

func hasJSONFlag(args []string) bool {
	for _, a := range args {
		if a == "--json" || a == "-json" || strings.HasPrefix(a, "--json=") || strings.HasPrefix(a, "-json=") {
			return true
		}
	}
	return false
}

func emitPureJSONOrFallback(buf *bytes.Buffer, d Deps) {
	data := bytes.TrimSpace(buf.Bytes())
	if len(data) == 0 {
		return
	}
	if json.Valid(data) {
		_, _ = d.Stdout.Write(append(data, '\n'))
		return
	}
	start := bytes.IndexByte(data, '{')
	end := bytes.LastIndexByte(data, '}')
	if start >= 0 && end > start && json.Valid(data[start:end+1]) {
		if start > 0 && d.Stderr != nil {
			_, _ = d.Stderr.Write(data[:start])
		}
		_, _ = d.Stdout.Write(append(data[start:end+1], '\n'))
		if end+1 < len(data) && d.Stderr != nil {
			_, _ = d.Stderr.Write(data[end+1:])
		}
		return
	}
	_, _ = d.Stdout.Write(buf.Bytes())
}
