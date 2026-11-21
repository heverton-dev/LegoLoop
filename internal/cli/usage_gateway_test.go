package cli

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

// TestUsoDoGatewayDizAFerramentaReal: o modo gateway expõe só call_mcp_tool; o texto de uso não pode dizer "4 ferramentas".
func TestUsoDoGatewayDizAFerramentaReal(t *testing.T) {
	var out, errb bytes.Buffer
	Run(context.Background(), nil, Deps{Stdout: &out, Stderr: &errb})
	uso := out.String() + errb.String()
	if strings.Contains(uso, "4 ferramentas") || !strings.Contains(uso, "call_mcp_tool") {
		t.Fatalf("uso do mcp-router --gateway desatualizado: %s", uso)
	}
}
