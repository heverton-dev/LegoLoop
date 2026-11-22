package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
	"github.com/heverton-dev/LegoLoop/internal/brain"
	"github.com/heverton-dev/LegoLoop/internal/canon"
	"github.com/heverton-dev/LegoLoop/internal/compor"
	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/orchestrix"
	"github.com/heverton-dev/LegoLoop/internal/refinar"
	"github.com/heverton-dev/LegoLoop/internal/runner"
	"github.com/heverton-dev/LegoLoop/internal/ste"
)

func TestVerbos(t *testing.T) {
	v := Verbos()
	if len(v) == 0 {
		t.Fatal("Verbos() retornou vazio")
	}
}

func TestAnticorposExecutarStatus(t *testing.T) {
	dir := t.TempDir()
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}

	if code := executarStatus(nil, dir, d); code != 2 {
		t.Fatalf("esperava exit 2 sem args, obtive %d", code)
	}
	if code := executarStatus([]string{"inexistente"}, dir, d); code != 1 {
		t.Fatalf("esperava exit 1 para ausente, obtive %d", code)
	}
	procDir := filepath.Join(dir, "meu-proc")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(procDir, brain.AntibodiesMarkdownFile), []byte("# Anticorpos\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if code := executarStatus([]string{"meu-proc"}, dir, d); code != 0 || !strings.Contains(out.String(), "Anticorpos") {
		t.Fatalf("esperava exit 0 com conteúdo, obtive %d (%q)", code, out.String())
	}
}

func TestSpeculativeBrainEDirExistsEReportImport(t *testing.T) {
	candFast := orchestrix.SpeculativeCandidate{
		Name:    "fast",
		Timeout: time.Second,
		Run: func(_ context.Context, prompt string) (string, runner.TokenUsage, error) {
			return "resp:" + prompt, runner.TokenUsage{Input: 10, Output: 5}, nil
		},
	}
	candHeavy := orchestrix.SpeculativeCandidate{
		Name:    "heavy",
		Timeout: time.Second,
		Run: func(_ context.Context, prompt string) (string, runner.TokenUsage, error) {
			return "heavy:" + prompt, runner.TokenUsage{Input: 20, Output: 10}, nil
		},
	}
	sb := &speculativeBrain{cand: candFast}
	if sb.Name() != "fast" || sb.Model() != "fast" || sb.Timeout() != time.Second || sb.Effort() != brain.EffortMedium {
		t.Fatalf("metadados inesperados: %s %s %v %v", sb.Name(), sb.Model(), sb.Timeout(), sb.Effort())
	}
	if _, err := sb.Chat(context.Background(), nil); err == nil {
		t.Fatal("esperava erro em Chat sem mensagens")
	}
	if txt, err := sb.Chat(context.Background(), []interview.Message{{Content: "oi"}}); err != nil || txt != "resp:oi" {
		t.Fatalf("Chat falhou: %q %v", txt, err)
	}
	if txt, used, tok, err := runSpeculativeWithBrain(context.Background(), candFast, candHeavy, "p1"); err != nil || txt != "resp:p1" || used != "fast" || tok.Input != 10 {
		t.Fatalf("runSpeculativeWithBrain fast falhou: %q %q %+v %v", txt, used, tok, err)
	}
	candFail := orchestrix.SpeculativeCandidate{
		Name:    "fast",
		Timeout: time.Second,
		Run: func(context.Context, string) (string, runner.TokenUsage, error) {
			return "", runner.TokenUsage{}, errors.New("falhou")
		},
	}
	if txt, used, tok, err := runSpeculativeWithBrain(context.Background(), candFail, candHeavy, "p2"); err != nil || !strings.HasPrefix(txt, "heavy:p2") || used != "heavy" || tok.Input != 20 {
		t.Fatalf("runSpeculativeWithBrain fallback falhou: %q %q %+v %v", txt, used, tok, err)
	}
	if _, _, _, err := runSpeculativeWithBrain(context.Background(), candFail, candFail, "p3"); err == nil {
		t.Fatal("esperava erro quando ambos falham")
	}

	dir := t.TempDir()
	if !dirExists(dir) || dirExists(filepath.Join(dir, "nao_existe")) {
		t.Fatal("dirExists incorreto")
	}

	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}
	res := canon.ImportResult{Name: "skill-x", TargetDir: dir, SHA256: "abc", FilesCount: 2}
	if code := reportImport(d, res, false); code != 0 || !strings.Contains(out.String(), "skill-x") {
		t.Fatalf("reportImport texto falhou: %d %q", code, out.String())
	}
	out.Reset()
	if code := reportImport(d, res, true); code != 0 || !strings.Contains(out.String(), `"abc"`) {
		t.Fatalf("reportImport json falhou: %d %q", code, out.String())
	}
}

func TestComporHelpers(t *testing.T) {
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}

	if _, err := confirmarCompor(comporOpts{}, d, &out); err == nil {
		t.Fatal("esperava erro sem stdin e sem --yes")
	}
	fnYes, err := confirmarCompor(comporOpts{yes: true}, d, &out)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := fnYes(compor.Plano{}, nil); err != nil || !ok {
		t.Fatalf("fnYes falhou: %v %v", ok, err)
	}

	if code := falhaCompor(d, compor.ErrNome); code != 2 {
		t.Fatalf("esperava 2 para ErrNome, obtive %d", code)
	}
	if code := falhaCompor(d, compor.ErrPedido); code != 2 {
		t.Fatalf("esperava 2 para ErrPedido, obtive %d", code)
	}
	if code := falhaCompor(d, compor.ErrCancelado); code != 1 {
		t.Fatalf("esperava 1 para ErrCancelado, obtive %d", code)
	}
	if code := falhaCompor(d, errors.New("boom")); code != 1 {
		t.Fatalf("esperava 1 para erro genérico, obtive %d", code)
	}

	res := &compor.Resultado{
		Nome: "proc",
		Arquivos: []compor.Arquivo{
			{Tipo: "dag", Estado: "atualizado", Caminho: "dag.yaml", Backup: "dag.bak"},
		},
	}
	out.Reset()
	if code := imprimirCompor(d, res, false); code != 0 || !strings.Contains(out.String(), "dag.yaml") {
		t.Fatalf("imprimirCompor texto falhou: %d %q", code, out.String())
	}
	out.Reset()
	if code := imprimirCompor(d, res, true); code != 0 || !strings.Contains(out.String(), `"dag.yaml"`) {
		t.Fatalf("imprimirCompor json falhou: %d %q", code, out.String())
	}

	// aprovarComposicao: entrada inválida -> personalizar (erro e sucesso) -> seguir
	in := strings.NewReader("xyz\n2\nNovo Titulo\na,b\n1\n")
	var pw bytes.Buffer
	p := newPrompter(in, &pw)
	pl := compor.Plano{Ligacao: compor.Ligacao{Composto: "legoloop-x@1.0.0"}}
	chamadas := 0
	replan := func(l compor.Ligacao) (compor.Plano, error) {
		chamadas++
		pl2 := pl
		pl2.Ligacao = l
		return pl2, nil
	}
	ok, err := p.aprovarComposicao(pl, replan)
	if err != nil || !ok || chamadas != 1 {
		t.Fatalf("aprovarComposicao falhou: ok=%v err=%v chamadas=%d", ok, err, chamadas)
	}

	// personalizarComposicao com erro no replan e cancelar
	ficha := ste.Cartao{
		Parametros: []ste.Parametro{
			{Nome: "titulo", Valor: "t1", Ajuda: "a1"},
			{Nome: "depende", Valor: "d1", Ajuda: "a2"},
		},
	}
	p2 := newPrompter(strings.NewReader("t2\nd2\n"), &pw)
	gotPl, err := p2.personalizarComposicao(ficha, pl, func(compor.Ligacao) (compor.Plano, error) {
		return compor.Plano{}, errors.New("invalido")
	})
	if err != nil || gotPl.Ligacao.Composto != pl.Ligacao.Composto {
		t.Fatalf("personalizarComposicao erro falhou: %+v %v", gotPl, err)
	}
	p3 := newPrompter(strings.NewReader("0\n"), &pw)
	if ok, err := p3.aprovarComposicao(pl, replan); err != nil || ok {
		t.Fatalf("esperava cancelar, obtive ok=%v err=%v", ok, err)
	}
}

func TestRefinarHelpers(t *testing.T) {
	if artefatoSai("") != nil {
		t.Fatal("artefatoSai vazio deve retornar nil")
	}
	if a := artefatoSai("rel.md"); a == nil || a.Tipo != "markdown" || !a.Promete {
		t.Fatalf("artefatoSai padrão incorreto: %+v", a)
	}
	if a := artefatoSai("out.json:json"); a == nil || a.Caminho != "out.json" || a.Tipo != "json" {
		t.Fatalf("artefatoSai json incorreto: %+v", a)
	}
	if a := artefatoSai("done:marcador"); a == nil || a.Promete {
		t.Fatalf("artefatoSai marcador incorreto: %+v", a)
	}

	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}
	if code := falhaRefinar(d, refinar.ErrNome); code != 2 {
		t.Fatalf("esperava 2, obtive %d", code)
	}
	if code := falhaRefinar(d, refinar.ErrMudanca); code != 2 {
		t.Fatalf("esperava 2, obtive %d", code)
	}
	if code := falhaRefinar(d, refinar.ErrCancelado); code != 1 {
		t.Fatalf("esperava 1, obtive %d", code)
	}
	if code := falhaRefinar(d, errors.New("boom")); code != 1 {
		t.Fatalf("esperava 1, obtive %d", code)
	}

	res := &refinar.Resultado{
		Nome: "proc",
		De:   "1.0.0",
		Para: "1.1.0",
		Arquivos: []refinar.Arquivo{
			{Tipo: "dag", Estado: "atualizado", Caminho: "dag.yaml"},
		},
	}
	if code := imprimirRefinar(d, res, false); code != 0 || !strings.Contains(out.String(), "dag.yaml") {
		t.Fatalf("imprimirRefinar texto falhou: %d %q", code, out.String())
	}
	out.Reset()
	if code := imprimirRefinar(d, res, true); code != 0 || !strings.Contains(out.String(), `"1.1.0"`) {
		t.Fatalf("imprimirRefinar json falhou: %d %q", code, out.String())
	}

	var pw bytes.Buffer
	p := newPrompter(strings.NewReader("xyz\n2\np2\nc2\n1\n"), &pw)
	pl := refinar.Plano{Mudanca: refinar.Mudanca{Etapa: "e1"}}
	ok, err := p.aprovarRefino(pl, func(m refinar.Mudanca) (refinar.Plano, error) {
		return refinar.Plano{Mudanca: m}, nil
	})
	if err != nil || !ok {
		t.Fatalf("aprovarRefino falhou: ok=%v err=%v", ok, err)
	}
	ficha := ste.Cartao{
		Parametros: []ste.Parametro{
			{Nome: "peca", Valor: "p1", Ajuda: "a1"},
			{Nome: "casos", Valor: "c1", Ajuda: "a2"},
		},
	}
	p2 := newPrompter(strings.NewReader("p2\nc2,c3\n"), &pw)
	if _, err := p2.personalizarRefino(ficha, pl, func(refinar.Mudanca) (refinar.Plano, error) {
		return refinar.Plano{}, errors.New("erro replan")
	}); err != nil {
		t.Fatalf("personalizarRefino erro falhou: %v", err)
	}
	p3 := newPrompter(strings.NewReader("0\n"), &pw)
	if ok, err := p3.aprovarRefino(pl, nil); err != nil || ok {
		t.Fatalf("esperava cancelar refino, obtive ok=%v err=%v", ok, err)
	}
}

func TestSwarmComandosUnit(t *testing.T) {
	t.Setenv("LEGOLOOP_SWARM_WORKER", "")
	if workerTargetFromEnv() != "" {
		t.Fatal("esperava vazio sem env")
	}
	t.Setenv("LEGOLOOP_SWARM_WORKER", "dev1")
	if workerTargetFromEnv() != "@dev1" {
		t.Fatalf("esperava @dev1, obtive %q", workerTargetFromEnv())
	}

	dir := t.TempDir()
	sb, err := board.OpenSwarmBoard(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sb.AddWorker(board.AddWorkerRequest{Name: "w1", Role: "engenheiro-tdd", WorkspaceMode: board.WorkspaceInherit, RepoDir: dir}); err != nil {
		t.Fatal(err)
	}
	if _, err := sb.PostMessage("coord", "@dev1", "oi"); err != nil {
		t.Fatal(err)
	}

	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}
	ctx := context.Background()

	if code := cmdSwarmWait(ctx, []string{"-flag-invalida"}, d); code != 2 {
		t.Fatalf("esperava 2, obtive %d", code)
	}
	if code := cmdSwarmWait(ctx, []string{"-dir", dir, "-target", "@dev1", "-timeout", "1s"}, d); code != 0 || !strings.Contains(out.String(), "oi") {
		t.Fatalf("cmdSwarmWait falhou: %d %q %q", code, out.String(), errBuf.String())
	}
	if code := cmdSwarmWait(ctx, []string{"-dir", dir, "-target", "@ninguem", "-timeout", "50ms"}, d); code != 1 {
		t.Fatalf("esperava 1 em timeout, obtive %d", code)
	}

	out.Reset()
	if code := cmdSwarmMonitor(ctx, []string{"-flag-invalida"}, d); code != 2 {
		t.Fatalf("esperava 2, obtive %d", code)
	}
	if code := cmdSwarmMonitor(ctx, []string{"-dir", dir, "-timeout", "1m"}, d); code != 0 || !strings.Contains(out.String(), "stalled") {
		t.Fatalf("cmdSwarmMonitor falhou: %d %q", code, out.String())
	}

	out.Reset()
	if code := cmdSwarmKill(ctx, []string{"-flag-invalida"}, d); code != 2 {
		t.Fatalf("esperava 2, obtive %d", code)
	}
	if code := cmdSwarmKill(ctx, []string{"-dir", dir, "-name", "nao_existe", "-repo", dir}, d); code != 1 {
		t.Fatalf("esperava 1 para worker inexistente, obtive %d", code)
	}
	if code := cmdSwarmKill(ctx, []string{"-dir", dir, "-name", "w1", "-repo", dir}, d); code != 0 || !strings.Contains(out.String(), "killed") {
		t.Fatalf("cmdSwarmKill falhou: %d %q %q", code, out.String(), errBuf.String())
	}

	if code := cmdSwarmSpawn(ctx, []string{"-flag-invalida"}, d); code != 2 {
		t.Fatalf("esperava 2, obtive %d", code)
	}
	if code := cmdSwarmSpawn(ctx, []string{"-dir", dir, "-name", "w1"}, d); code != 2 {
		t.Fatalf("esperava 2 sem comando, obtive %d", code)
	}
	if code := cmdSwarmSpawn(ctx, []string{"-dir", dir, "-name", "nao_existe", "--", "echo", "oi"}, d); code != 1 {
		t.Fatalf("esperava 1 para worker inexistente, obtive %d", code)
	}
}

func TestAuditDoubleSurfaceBookECmdCommitFlags(t *testing.T) {
	dir := t.TempDir()
	var out, errBuf bytes.Buffer
	d := Deps{Stdout: &out, Stderr: &errBuf}

	if code := auditDoubleSurfaceBook(dir, "", d); code != 1 {
		t.Fatalf("esperava 1 sem PDF, obtive %d", code)
	}
	sub := filepath.Join(dir, "edicao")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "livro_executivo.pdf"), []byte("%PDF"), 0o600); err != nil {
		t.Fatal(err)
	}
	outJSON := filepath.Join(dir, "audit.json")
	if code := auditDoubleSurfaceBook(dir, outJSON, d); code != 0 {
		t.Fatalf("esperava 0 com PDF em subpasta, obtive %d (%q)", code, errBuf.String())
	}
	if _, err := os.Stat(outJSON); err != nil {
		t.Fatalf("audit.json não criado: %v", err)
	}

	if code := cmdCommit(context.Background(), []string{"-flag-invalida"}, d); code != 2 {
		t.Fatalf("cmdCommit flag inválida: esperava 2, obtive %d", code)
	}
}
