package cli

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"legoloop/internal/fabricar"
)

// cerebroPorPapel é um 9Router falso que responde conforme o papel pedido no prompt do sistema:
// necessidades, explicação da Fusão, SPEC ou entrevista. Conta as chamadas de cada papel.
type cerebroPorPapel struct{ needs, explain, spec, outros atomic.Int32 }

const explicacaoFalsa = "testa antes de construir"

func (c *cerebroPorPapel) responder(system string) string {
	switch {
	case strings.Contains(system, "vocabulário fechado do LegoLoop"):
		c.needs.Add(1)
		return `{"natureza": "software", "etapas": ["construir", "verificar"], "casos": []}`
	case strings.Contains(system, "explicador da Fusão"):
		c.explain.Add(1)
		return `{"explicacoes": {"d1": "` + explicacaoFalsa + `", "d2": "` + explicacaoFalsa + `", "d3": "` + explicacaoFalsa + `", "d4": "` + explicacaoFalsa + `"}}`
	case strings.Contains(system, "Arquiteto de Especificações"):
		c.spec.Add(1)
		return specFake
	}
	c.outros.Add(1)
	return "HIPÓTESE: x\nCONFIANÇA: 96%\nPERGUNTA: ok?\nPALPITE: sim"
}

func (c *cerebroPorPapel) servidor(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/models":
			_ = json.NewEncoder(w).Encode(map[string]any{"data": []map[string]string{{"id": "code-fast"}}})
		case "/v1/chat/completions":
			var req struct {
				Messages []struct{ Role, Content string } `json:"messages"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			system := ""
			if len(req.Messages) > 0 {
				system = req.Messages[0].Content
			}
			reply := c.responder(system)
			_ = json.NewEncoder(w).Encode(map[string]any{"choices": []map[string]any{{"message": map[string]string{"content": reply}}}})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

// matrizDeTeste grava no cofre uma matriz v2 com três skills: aidd-tdd (apta, Score baixo), aidd-site-dominator
// (Score alto, mas exclui software) e sem-perfil (Score mais alto, sem perfil de uso), e um MCP. Os ref seguem
// o formato real do registry_matrix.json: relativos ao HOME ("~/...") e, no MCP, "arquivo#nome".
func matrizDeTeste(t *testing.T, vault string) {
	t.Helper()
	ref := func(n string) string { return "~/.claude/skills/" + n + "/SKILL.md" }
	type perfil = map[string]any
	ent := func(n string, p perfil) map[string]any {
		e := map[string]any{"ref": ref(n), "nome": n, "tipo": "skill", "hash": "h-" + n}
		if p != nil {
			e["perfil"] = p
		}
		return e
	}
	reg := map[string]any{
		"generated_at": "2026-11-01T10:00:00Z", "matrix_version": "2.0.0", "total_pieces": 4,
		"ratings": []map[string]any{
			{"ref": ref("aidd-tdd"), "name": "aidd-tdd", "total_score": 0.58},
			{"ref": ref("aidd-site-dominator"), "name": "aidd-site-dominator", "total_score": 0.90},
			{"ref": ref("sem-perfil"), "name": "sem-perfil", "total_score": 0.95},
		},
		"entradas": []map[string]any{
			ent("aidd-tdd", perfil{"etapas": []string{"construir", "verificar"}, "naturezas": []string{"software"}, "casos_de_uso": []string{}, "origem": "declarado"}),
			ent("aidd-site-dominator", perfil{"etapas": []string{"publicar"}, "naturezas": []string{"midia"}, "casos_de_uso": []string{}, "quando_nao_usar": []string{"software"}, "origem": "declarado"}),
			ent("sem-perfil", nil),
			{"ref": "~/.claude.json#designi-mcp", "nome": "designi-mcp", "tipo": "mcp", "hash": "h-mcp"},
		},
	}
	b, err := json.Marshal(reg)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(vault, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vault, "registry_matrix.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
}

type saidaMotor struct {
	Report struct {
		Varredura struct {
			PecasVistas map[string]int `json:"pecas_vistas"`
			MatrizIdade int64          `json:"matriz_idade"`
		} `json:"varredura"`
		Fusao struct {
			Necessidades struct {
				Natureza string
				Etapas   []string
			} `json:"necessidades"`
			Decisoes []fabricar.Decisao `json:"decisoes"`
		} `json:"fusao"`
	} `json:"report"`
}

func decisaoDe(ds []fabricar.Decisao, peca string) (fabricar.Decisao, bool) {
	for _, d := range ds {
		if d.Peca == peca {
			return d, true
		}
	}
	return fabricar.Decisao{}, false
}

// TestMotorRealBinarioEscolhePelaAptidaoDaMatriz prova pelo binário que as Camadas 2 e 3 são as da v2:
// o fabricar lê a matriz do cofre (não varre drives), o cérebro dá as necessidades e explica, e a escolha
// segue a aptidão, não o Score.
func TestMotorRealBinarioEscolhePelaAptidaoDaMatriz(t *testing.T) {
	bin := buildBinary(t)
	cer := &cerebroPorPapel{}
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeTeste(t, vault)
	r := runBinaryEm(t, bin, cer.servidor(t), root, proj, "fabricar", "extrair audio", "--goal", "construir e testar um extrator",
		"--vault", vault, "--docs", proj, "--yes", "--sem-entrevista", "--dry-run", "--json", "--no-link")
	if r.code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", r.code, r.out, r.errb)
	}
	if strings.Contains(r.out+r.errb, "mídia removível") {
		t.Errorf("peça do HOME (ref \"~/...\") não é mídia removível:\n%s", r.errb)
	}
	if strings.Contains(r.out+r.errb, "Drives detectados") || strings.Contains(r.out+r.errb, "Varredura...") {
		t.Errorf("o fabricar não deve varrer drives:\n%s", r.errb)
	}
	var s saidaMotor
	if err := json.Unmarshal([]byte(r.out), &s); err != nil {
		t.Fatalf("stdout não é JSON: %v\n%s", err, r.out)
	}
	f := s.Report.Fusao
	if f.Necessidades.Natureza != "software" || len(f.Necessidades.Etapas) != 2 {
		t.Errorf("necessidades devem vir do cérebro: %+v", f.Necessidades)
	}
	tdd, ok := decisaoDe(f.Decisoes, "aidd-tdd")
	if !ok || !tdd.Entrou || tdd.Explicacao != explicacaoFalsa {
		t.Errorf("aidd-tdd deve entrar, com a explicação do cérebro: %+v", f.Decisoes)
	}
	if d, ok := decisaoDe(f.Decisoes, "aidd-site-dominator"); !ok || d.Entrou || d.Motivo == "" {
		t.Errorf("aidd-site-dominator deve sair com motivo, apesar do Score 0.90: %+v", f.Decisoes)
	}
	if d, ok := decisaoDe(f.Decisoes, "sem-perfil"); !ok || d.Entrou || !strings.Contains(d.Motivo, "aptidão desconhecida") {
		t.Errorf("sem-perfil deve aparecer como aptidão desconhecida: %+v", f.Decisoes)
	}
	for _, d := range f.Decisoes {
		if d.Motivo == "selecionada na fusão" {
			t.Errorf("decisão inventada no relatório: %+v", d)
		}
	}
	if s.Report.Varredura.PecasVistas["skill"] != 3 || s.Report.Varredura.MatrizIdade <= 0 {
		t.Errorf("a Camada 2 deve contar as peças da matriz e a idade dela: %+v", s.Report.Varredura)
	}
	if cer.needs.Load() != 1 || cer.explain.Load() != 1 {
		t.Errorf("o cérebro deve ser chamado 1 vez para as necessidades e 1 para explicar: needs=%d explain=%d", cer.needs.Load(), cer.explain.Load())
	}
}

// TestMotorRealBinarioSemMatrizParaComAviso: sem registry_matrix.json o fabricar sai com 1 e diz como gerar.
func TestMotorRealBinarioSemMatrizParaComAviso(t *testing.T) {
	bin := buildBinary(t)
	cer := &cerebroPorPapel{}
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	r := runBinaryEm(t, bin, cer.servidor(t), root, proj, "fabricar", "extrair audio", "--vault", vault, "--docs", proj,
		"--yes", "--sem-entrevista", "--dry-run", "--no-link")
	if r.code != 1 || !strings.Contains(r.out+r.errb, "legoloop matriz gerar") {
		t.Fatalf("esperava exit 1 e o aviso para gerar a matriz; exit %d\n%s\n%s", r.code, r.out, r.errb)
	}
	if cer.spec.Load() != 0 {
		t.Errorf("sem matriz, a SPEC não deve ser pedida ao cérebro")
	}
}
