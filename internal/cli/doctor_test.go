package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDoctorListaPythonOpcional(t *testing.T) {
	d, out, _ := deps(t)
	if code := run(t, d, "doctor"); code != 0 {
		t.Fatalf("doctor falhou: code=%d", code)
	}
	texto := out.String()
	encontrou := false
	for _, linha := range strings.Split(texto, "\n") {
		baixa := strings.ToLower(linha)
		if strings.Contains(baixa, "python") && strings.Contains(baixa, "opcional") {
			encontrou = true
			break
		}
	}
	if !encontrou {
		t.Errorf("doctor não listou Python como opcional:\n%s", texto)
	}

	d2, out2, _ := deps(t)
	if code := run(t, d2, "doctor", "--json"); code != 0 {
		t.Fatalf("doctor --json falhou: code=%d", code)
	}
	var payload map[string]any
	if err := json.Unmarshal(out2.Bytes(), &payload); err != nil {
		t.Fatalf("--json inválido: %v", err)
	}
	if got, ok := payload["python"].(string); !ok || got != "opcional" {
		t.Errorf("doctor --json campo python = %v, esperado %q", payload["python"], "opcional")
	}
}
