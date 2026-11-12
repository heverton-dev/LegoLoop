package roteamento

import (
	"reflect"
	"testing"
)

func storePerfil() map[string]Homologacao {
	return map[string]Homologacao{
		"claude|oficial":   {Estado: EstadoHomologado},
		"claude|9router":   {Estado: EstadoHomologado},
		"opencode|oficial": {Estado: EstadoPendente},
	}
}

func TestPerfilDoisNiveis(t *testing.T) {
	global := []byte(`{"schema_version":1,"candidatos":{"micro":[{"harness":"claude","via":"oficial","modelo":"opus","esforco":"medium"}]}}`)
	projeto := []byte(`{"schema_version":1,"candidatos":{"nano":[{"harness":"claude","via":"9router","modelo":"combo-a","esforco":"low"}]}}`)
	got, err := CarregarPerfil(global, projeto, storePerfil())
	if err != nil {
		t.Fatal(err)
	}
	quer := map[string][]CandidatoVia{
		ClasseNano:    {{Harness: "claude", Via: Via9Router, Modelo: "combo-a", Esforco: "low"}},
		ClasseMicro:   {{Harness: "claude", Via: ViaOficial, Modelo: "opus", Esforco: "medium"}},
		ClasseGrande:  {{Harness: "claude", Via: ViaOficial, Modelo: "sonnet", Esforco: "high"}},
		ClasseExtreme: {{Harness: "claude", Via: ViaOficial, Modelo: "opus", Esforco: "high"}},
	}
	if !reflect.DeepEqual(got, quer) {
		t.Fatalf("perfil = %+v, quer %+v", got, quer)
	}
}

func TestPerfilGatewayDesabilitado(t *testing.T) {
	projeto := []byte(`{"gateways":{"9router":{"habilitado":false}},"candidatos":{"nano":[{"harness":"claude","via":"9router","modelo":"combo-a","esforco":"low"}]}}`)
	got, err := CarregarPerfil(nil, projeto, storePerfil())
	if err != nil {
		t.Fatal(err)
	}
	if len(got[ClasseNano]) != 0 {
		t.Fatalf("nano = %+v, quer vazio", got[ClasseNano])
	}
}

func TestPerfilDescartaNaoHomologado(t *testing.T) {
	global := []byte(`{"candidatos":{"micro":[{"harness":"opencode","via":"oficial","modelo":"x","esforco":"low"}]}}`)
	got, err := CarregarPerfil(global, nil, storePerfil())
	if err != nil {
		t.Fatal(err)
	}
	if len(got[ClasseMicro]) != 0 {
		t.Fatalf("micro = %+v, quer vazio", got[ClasseMicro])
	}
}

func TestPerfilInvalido(t *testing.T) {
	if _, err := CarregarPerfil([]byte("{"), nil, storePerfil()); err == nil {
		t.Fatal("esperava erro")
	}
}
