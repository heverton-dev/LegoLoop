package refinar_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/processo"
	"legoloop/internal/refinar"
)

// noRota e noDag imitam os nós curados à mão: roteamento 9router, timeout e gate com '&&' sem escape HTML.
func noRota(n int, art string) string {
	return fmt.Sprintf(`    {
      "id": "TICK-01.%d",
      "label": "`+"`e0%d`"+` — Etapa antiga — peça `+"`velha`"+`",
      "ref": null,
      "timeout_sec": 180,
      "routing": {"harness": "9router", "model": "t-codigo-128k"},
      "output": {
        "artifacts": ["%s"],
        "gate": "file_exists('%s') && word_count('%s') >= 30"
      }
    }`, n, n, art, art, art)
}

func noDag(n int, art string) string {
	return fmt.Sprintf(`    {
      "id": "TICK-01.%d",
      "name": "Etapa antiga — peça velha",
      "ref": "skills/velha",
      "timeout_sec": 180,
      "routing": {"harness": "9router", "model": "t-codigo-128k"},
      "output": {
        "artifact": "%s",
        "gate": "file_exists('%s') && word_count('%s') >= 30"
      }
    }`, n, art, art, art)
}

const stopGate = `    {
      "id": "SG-FINAL",
      "name": "Homologação Final",
      "timeout_sec": 180,
      "output": {"artifact": "bin/legoloop.exe", "gate": "exit_code == 0"}
    }`

func docV2(cab string, nos []string) string {
	return "{\n  \"schema_version\": 2,\n" + cab + "  \"nodes\": [\n" + strings.Join(nos, ",\n") + "\n  ]\n}\n"
}

func plantarV2(t *testing.T, c cofre, nTicks int) (rota, espelho, projecao string) {
	t.Helper()
	pr, err := processo.Projetar(processoDeTeste(), inicio)
	if err != nil {
		t.Fatal(err)
	}
	var nosR, nosD []string
	for i := 1; i <= nTicks; i++ {
		nosR = append(nosR, noRota(i, fmt.Sprintf("out/0%d_etapa.md", i)))
		nosD = append(nosD, noDag(i, fmt.Sprintf("out/0%d_etapa.md", i)))
	}
	nosR, nosD = append(nosR, stopGate), append(nosD, stopGate)
	rota = filepath.Join(c.docs, "docs", "routes", pr.RouteNome)
	espelho = filepath.Join(c.docs, "docs", "dags", pr.DAGNome)
	projecao = filepath.Join(c.docs, "processos", nome, caminhos.ArquivoDAG)
	cabDag := "  \"upstream_signatures\": {\n    \"route\": {\n      \"path\": \"docs/routes/" + pr.RouteNome +
		"\",\n      \"sha256\": \"" + strings.Repeat("0", 64) + "\"\n    }\n  },\n"
	for p, conteudo := range map[string]string{
		rota:     docV2("  \"template_version\": \"2.2.0\",\n", nosR),
		espelho:  docV2(cabDag, nosD),
		projecao: docV2(cabDag, nosD),
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(conteudo), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return rota, espelho, projecao
}

func rodarRefinar(t *testing.T, c cofre) refinar.Plano {
	t.Helper()
	var plano refinar.Plano
	confirm := func(pl refinar.Plano, _ func(refinar.Mudanca) (refinar.Plano, error)) (bool, error) {
		plano = pl
		return true, nil
	}
	if _, err := refinar.Run(context.Background(), refinar.Request{Nome: nome, Mudanca: mudanca}, depsDe(c, nil, confirm, inicio)); err != nil {
		t.Fatal(err)
	}
	return plano
}

// O refinar não pode trocar ROUTE e DAGs v2 curados pelo formato v1 gerado: o roteamento e o stop-gate ficam.
func TestRefinarPreservaV2Curado(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	rota, espelho, projecao := plantarV2(t, c, 3)
	rodarRefinar(t, c)
	p := lerProcesso(t, c)
	soma := sha256.Sum256([]byte(ler(t, rota)))
	shaRota := hex.EncodeToString(soma[:])
	for _, caminho := range []string{rota, espelho, projecao} {
		txt := ler(t, caminho)
		for _, quer := range []string{`"schema_version": 2`, `"harness": "9router"`, `"timeout_sec": 180`, `SG-FINAL`} {
			if !strings.Contains(txt, quer) {
				t.Errorf("%s perdeu %q:\n%s", caminho, quer, txt)
			}
		}
		if strings.Contains(txt, `\u0026`) {
			t.Errorf("%s escapou o && do gate", caminho)
		}
		for i, e := range p.Etapas {
			if !strings.Contains(txt, e.Gate) {
				t.Errorf("%s: o gate da etapa %s não entrou no arquivo", caminho, e.ID)
			}
			if !strings.Contains(txt, e.Produz.Caminho) || strings.Contains(txt, fmt.Sprintf("out/0%d_etapa.md", i+1)) {
				t.Errorf("%s: etapa %s não trocou o artefato", caminho, e.ID)
			}
		}
	}
	if !strings.Contains(ler(t, rota), "`e02` — Etapa e02 — peça `aidd-qa`") {
		t.Errorf("o label do ROUTE não refletiu a peça nova:\n%s", ler(t, rota))
	}
	for _, caminho := range []string{espelho, projecao} {
		txt := ler(t, caminho)
		if !strings.Contains(txt, `"ref": "skills/aidd-qa"`) || !strings.Contains(txt, shaRota) {
			t.Errorf("%s: ref ou sha do ROUTE desatualizados:\n%s", caminho, txt)
		}
	}
}

// Número de nós TICK diferente do número de etapas: o arquivo fica como está e o plano avisa.
func TestRefinarMantemV2ComContagemDiferente(t *testing.T) {
	c := semear(t, processoDeTeste(), inicio, true)
	rota, _, _ := plantarV2(t, c, 2)
	antes := ler(t, rota)
	plano := rodarRefinar(t, c)
	if ler(t, rota) != antes {
		t.Error("ROUTE v2 com contagem diferente não pode mudar")
	}
	achou := false
	for _, a := range plano.Avisos {
		if strings.Contains(a, "v2 curado") {
			achou = true
		}
	}
	if !achou {
		t.Errorf("faltou o aviso de v2 mantido: %v", plano.Avisos)
	}
}
