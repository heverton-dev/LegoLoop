package dag_test

import (
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/dag"
)

func unitDoc(tipo, extra string) []byte {
	return []byte(`{"name":"x","version":"1","trigger":{"type":"manual"},"nodes":[{` +
		`"id":"a","type":"` + tipo + `",` + extra + `"routing":{"mode":"manual","model":"m"},` +
		`"output":{"artifact":"A.md","gate":"true"}}]}`)
}

func TestUnitNoWriteBoundaryOnlyStrict(t *testing.T) {
	if _, err := dag.Parse(unitDoc("unit_no", `"isolation":{"write_boundary":"strict"},`)); err != nil {
		t.Fatalf("strict deve valer: %v", err)
	}
	_, err := dag.Parse(unitDoc("unit_no", `"isolation":{"write_boundary":"none"},`))
	if err == nil || !strings.Contains(err.Error(), "write_boundary") {
		t.Fatalf("unit_no com write_boundary none deve ser recusado: %v", err)
	}
}

func TestStopGateKeepsWriteBoundaryNone(t *testing.T) {
	if _, err := dag.Parse(unitDoc("stop_gate", `"isolation":{"write_boundary":"none"},`)); err != nil {
		t.Fatalf("stop_gate com none deve valer: %v", err)
	}
}

func TestRetryPolicyRules(t *testing.T) {
	ok := []string{`"retry_policy":{"max_retries":0},`, `"retry_policy":{"max_retries":1},`, `"retry_policy":{"max_retries":2},`}
	for _, e := range ok {
		d, err := dag.Parse(unitDoc("unit_no", e))
		if err != nil {
			t.Errorf("%s deve valer: %v", e, err)
			continue
		}
		if d.Nodes[0].RetryPolicy == nil {
			t.Errorf("%s: retry_policy não carregado", e)
		}
	}
	bad := map[string]string{
		"acima do teto": `"retry_policy":{"max_retries":3},`,
		"negativo":      `"retry_policy":{"max_retries":-1},`,
	}
	for nome, e := range bad {
		if _, err := dag.Parse(unitDoc("unit_no", e)); err == nil || !strings.Contains(err.Error(), "max_retries") {
			t.Errorf("%s deve ser recusado: %v", nome, err)
		}
	}
	if _, err := dag.Parse(unitDoc("script", `"command":["echo","x"],"retry_policy":{"max_retries":1},`)); err == nil ||
		!strings.Contains(err.Error(), "retry_policy") {
		t.Errorf("retry_policy fora do unit_no deve ser recusado: %v", err)
	}
}
