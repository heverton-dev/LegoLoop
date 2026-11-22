package runner

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/triagem"
)

func TestMutacaoAllowlistSliceAny(t *testing.T) {
	n := dag.Node{
		ID:      "n1",
		Tooling: map[string]any{"allowlist": []any{" read ", "   ", 42}},
	}
	if err := ValidateToolAllowlist(n, "read"); err != nil {
		t.Fatalf("esperava read permitido: %v", err)
	}
	if err := ValidateToolAllowlist(n, "write"); err == nil {
		t.Fatal("esperava write bloqueado")
	}
	if err := ValidateToolAllowlist(n, ""); err == nil {
		t.Fatal("esperava ferramenta vazia bloqueada")
	}
}

func TestMutacaoAgendaCabeLimiteExato(t *testing.T) {
	a := novaAgenda(3, func([]int) int { return 2 })
	a.ativos["n1"] = ativo{arquivos: []string{"a.go"}, pesoMB: 500}
	if !a.cabe([]string{"b.go"}, 500) {
		t.Fatal("esperava caber quando ativos+1 == workers")
	}
}

func TestMutacaoIsGitRepoECheckGitRepo(t *testing.T) {
	semGit := t.TempDir()
	if isGitRepo(semGit) {
		t.Fatal("diretorio sem .git nao eh repo git")
	}
	ok, err := checkGitRepo(semGit)
	if err != nil || ok {
		t.Fatalf("esperava false, nil; obtive %v, %v", ok, err)
	}
	comGit := t.TempDir()
	if err := os.Mkdir(filepath.Join(comGit, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	if !isGitRepo(comGit) {
		t.Fatal("diretorio com .git deve ser repo git")
	}
	ok, err = checkGitRepo(comGit)
	if err != nil || !ok {
		t.Fatalf("esperava true, nil; obtive %v, %v", ok, err)
	}
	ok, err = checkGitRepo("invalido\x00dir")
	if err == nil || ok {
		t.Fatalf("esperava false, erro para caminho invalido; obtive %v, %v", ok, err)
	}
}

func criarLinkOuJunction(t *testing.T, alvo, link string) {
	t.Helper()
	if err := os.Symlink(alvo, link); err == nil {
		return
	}
	cmd := exec.Command("cmd", "/c", "mklink", "/J", link, alvo)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("falha ao criar junction: %v (%s)", err, string(out))
	}
}

func TestMutacaoCheckGitRepoLinkEPrepareIsolation(t *testing.T) {
	base := t.TempDir()
	alvo := filepath.Join(base, "alvo")
	dir := filepath.Join(base, "repo")
	if err := os.MkdirAll(alvo, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	criarLinkOuJunction(t, alvo, filepath.Join(dir, ".git"))
	ok, err := checkGitRepo(dir)
	if err == nil || ok {
		t.Fatalf("esperava false, erro para .git link; obtive %v, %v", ok, err)
	}
	if _, _, err := createIsolationEnv(context.Background(), dir); err == nil || !strings.Contains(err.Error(), "em .git:") {
		t.Fatalf("esperava erro em createIsolationEnv com .git link, obtive: %v", err)
	}
	if _, _, err := prepareIsolation(context.Background(), dir, dag.Node{}); err == nil {
		t.Fatal("esperava erro em prepareIsolation com .git link")
	}
}

func TestMutacaoMkdirTempFalhaEmIsolation(t *testing.T) {
	arq := filepath.Join(t.TempDir(), "nao_eh_dir")
	if err := os.WriteFile(arq, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMP", arq)
	t.Setenv("TEMP", arq)
	t.Setenv("TMPDIR", arq)
	wd := t.TempDir()
	if _, _, err := createGitWorktree(context.Background(), wd); err == nil || strings.Contains(err.Error(), "git worktree add falhou") {
		t.Fatalf("esperava erro de MkdirTemp em createGitWorktree, obtive: %v", err)
	}
	if _, _, err := createTempCopy(wd); err == nil {
		t.Fatal("esperava erro em createTempCopy com TMP invalido")
	}
}

func TestMutacaoCopyTreeRaizEIgnorados(t *testing.T) {
	base := t.TempDir()
	src := filepath.Join(base, "bin")
	dst := filepath.Join(base, "dst")
	for _, sub := range []string{"node_modules", ".ssh"} {
		if err := os.MkdirAll(filepath.Join(src, sub), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	_ = os.WriteFile(filepath.Join(src, "app.txt"), []byte("ok"), 0o644)
	_ = os.WriteFile(filepath.Join(src, "node_modules", "x.js"), []byte("1"), 0o644)
	_ = os.WriteFile(filepath.Join(src, ".ssh", "k"), []byte("2"), 0o644)
	if err := copyTree(src, dst); err != nil {
		t.Fatalf("copyTree falhou: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "app.txt")); err != nil {
		t.Fatalf("app.txt deveria existir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dst, "node_modules")); !os.IsNotExist(err) {
		t.Fatal("node_modules deveria ser ignorado")
	}
	if _, err := os.Stat(filepath.Join(dst, ".ssh")); !os.IsNotExist(err) {
		t.Fatal(".ssh deveria ser ignorado")
	}
}

func TestMutacaoIsIgnoredDir(t *testing.T) {
	for _, nome := range []string{".git", "node_modules", "bin", ".gemini", ".claude"} {
		if !isIgnoredDir(nome) {
			t.Fatalf("esperava %s ignorado", nome)
		}
	}
	if isIgnoredDir("src") {
		t.Fatal("src nao deve ser ignorado")
	}
}

func TestMutacaoInjectContextInErrosEMultiplos(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	for _, f := range []string{"a.txt", "b.txt"} {
		_ = os.WriteFile(filepath.Join(src, f), []byte(f), 0o644)
	}
	if err := injectContextIn(src, dst, []string{"a.txt", "b.txt"}); err != nil {
		t.Fatalf("injectContextIn falhou: %v", err)
	}
	for _, f := range []string{"a.txt", "b.txt"} {
		dados, err := os.ReadFile(filepath.Join(dst, f))
		if err != nil || string(dados) != f {
			t.Fatalf("arquivo %s invalido: %q (%v)", f, string(dados), err)
		}
	}
	if err := checkNoLinks(filepath.Join(src, "inexistente"), "a.txt", false); err == nil {
		t.Fatal("esperava erro em checkNoLinks com base inexistente")
	}
	if err := injectContextIn(src, dst, []string{"../fora"}); err == nil || !strings.Contains(err.Error(), "path jail") {
		t.Fatalf("esperava erro path jail, obtive: %v", err)
	}
	if err := injectContextIn(src, filepath.Join(dst, "nao_existe"), []string{"a.txt"}); err == nil {
		t.Fatal("esperava erro quando dstWorkdir nao existe")
	}
	_ = os.Mkdir(filepath.Join(dst, "colide"), 0o755)
	_ = os.WriteFile(filepath.Join(src, "colide"), []byte("x"), 0o644)
	if err := injectContextIn(src, dst, []string{"colide"}); err == nil {
		t.Fatal("esperava erro quando destino eh diretorio")
	}
}

func TestMutacaoCopyFileECopyRegularBytes(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	_ = os.Mkdir(sub, 0o755)
	if err := copyFile(filepath.Join(dir, "inexistente"), filepath.Join(dir, "x")); err == nil {
		t.Fatal("esperava erro para origem inexistente em copyFile")
	}
	err := copyFile(sub, filepath.Join(dir, "out"))
	if err == nil || !strings.Contains(err.Error(), "origem não é arquivo regular") {
		t.Fatalf("esperava erro de nao regular, obtive: %v", err)
	}
	src := filepath.Join(dir, "in.txt")
	dst := filepath.Join(dir, "out.txt")
	_ = os.WriteFile(src, []byte("ola"), 0o644)
	if err := copyRegularBytes(src, dst); err != nil {
		t.Fatalf("copyRegularBytes falhou: %v", err)
	}
	if dados, err := os.ReadFile(dst); err != nil || string(dados) != "ola" {
		t.Fatalf("conteudo esperado ola, obtive %q (%v)", string(dados), err)
	}
	dstNovo := filepath.Join(dir, "nao_criar.txt")
	if err := copyRegularBytes(filepath.Join(dir, "nao_existe"), dstNovo); err == nil {
		t.Fatal("esperava erro para origem inexistente")
	}
	if _, err := os.Stat(dstNovo); !os.IsNotExist(err) {
		t.Fatal("destino nao deveria ser criado quando origem nao existe")
	}
	if err := copyRegularBytes(src, sub); err == nil || errors.Is(err, os.ErrInvalid) {
		t.Fatalf("esperava erro de os.Create (nao ErrInvalid), obtive: %v", err)
	}
	if err := copyRegularBytes(sub, filepath.Join(dir, "bad.txt")); err == nil || errors.Is(err, os.ErrClosed) {
		t.Fatalf("esperava erro de io.Copy (nao ErrClosed), obtive: %v", err)
	}
}

func TestMutacaoPromoteArtifactETransferRootFile(t *testing.T) {
	iso := t.TempDir()
	mainDir := t.TempDir()
	srcSub := filepath.Join(iso, "sub", "dir")
	_ = os.MkdirAll(srcSub, 0o755)
	_ = os.WriteFile(filepath.Join(srcSub, "out.txt"), []byte("dado"), 0o644)
	if err := promoteArtifact(iso, mainDir, "sub/dir/out.txt"); err != nil {
		t.Fatalf("promoteArtifact falhou: %v", err)
	}
	if dados, err := os.ReadFile(filepath.Join(mainDir, "sub", "dir", "out.txt")); err != nil || string(dados) != "dado" {
		t.Fatalf("conteudo esperado dado, obtive %q (%v)", string(dados), err)
	}
	if err := promoteArtifact(iso, mainDir, "../fora"); err == nil || !strings.Contains(err.Error(), "path jail") {
		t.Fatalf("esperava erro path jail, obtive: %v", err)
	}
	if err := promoteArtifact(iso, mainDir, "faltando.txt"); err == nil || !strings.Contains(err.Error(), "inválido na origem") {
		t.Fatalf("esperava erro invalido na origem, obtive: %v", err)
	}
}

func TestMutacaoTransferRootFileErros(t *testing.T) {
	iso := t.TempDir()
	mainDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(mainDir, "pai_arq"), []byte("x"), 0o644)
	_ = os.Mkdir(filepath.Join(mainDir, "colide"), 0o755)
	_ = os.WriteFile(filepath.Join(iso, "colide"), []byte("x"), 0o644)
	_ = os.Mkdir(filepath.Join(iso, "umdir"), 0o755)
	srcRoot, _ := os.OpenRoot(iso)
	defer srcRoot.Close()
	dstRoot, _ := os.OpenRoot(mainDir)
	defer dstRoot.Close()
	if err := transferRootFile(srcRoot, dstRoot, filepath.Join("pai_arq", "f.txt")); err == nil {
		t.Fatal("esperava erro em MkdirAll quando pai eh arquivo")
	}
	if err := transferRootFile(srcRoot, dstRoot, "inexistente.txt"); err == nil {
		t.Fatal("esperava erro em srcRoot.Open inexistente")
	}
	if _, err := os.Stat(filepath.Join(mainDir, "inexistente.txt")); !os.IsNotExist(err) {
		t.Fatal("destino nao deveria ser criado quando srcRoot.Open falha")
	}
	if err := transferRootFile(srcRoot, dstRoot, "colide"); err == nil || errors.Is(err, os.ErrInvalid) {
		t.Fatalf("esperava erro em dstRoot.Create (nao ErrInvalid), obtive: %v", err)
	}
	if err := transferRootFile(srcRoot, dstRoot, "umdir"); err == nil || errors.Is(err, os.ErrClosed) {
		t.Fatalf("esperava erro de io.Copy (nao ErrClosed), obtive: %v", err)
	}
}

func TestMutacaoMemoriaRuntime(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	oldOut := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOut)
	recordEtapaMemoria(dir, "", "et0", "fim")
	if buf.Len() != 0 {
		t.Fatalf("sessao vazia nao deveria tentar gravar: %s", buf.String())
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatalf("sessao vazia nao deveria criar arquivos em dir: %v", entries)
	}
	if err := triagem.SalvarMemoria(dir, triagem.MemoriaExecucao{
		SessaoID:    "s1",
		PID:         os.Getpid(),
		EtapaAtual:  "et1",
		TokenBudget: &triagem.TokenBudgetTelemetry{NodeID: "n1", EstimatedTokens: 77},
	}); err != nil {
		t.Fatal(err)
	}
	recordEtapaMemoria(dir, "s1", "et2", "fim")
	if strings.Contains(buf.String(), "memoria: falha ao gravar") {
		t.Fatalf("nao esperava log de erro em sucesso: %s", buf.String())
	}
	mem, err := triagem.LerMemoria(dir, "s1")
	if err != nil || mem.TokenBudget == nil || mem.TokenBudget.EstimatedTokens != 77 {
		t.Fatalf("esperava TokenBudget preservado: %+v (%v)", mem, err)
	}
	arq := filepath.Join(dir, "arq")
	_ = os.WriteFile(arq, []byte("x"), 0o644)
	recordEtapaMemoria(arq, "s2", "et3", "fim")
	if !strings.Contains(buf.String(), "memoria: falha ao gravar") {
		t.Fatal("esperava log de falha quando memDir eh arquivo")
	}
}

func TestMutacaoVerifyMutationEResolveSafePath(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a"), 0o644)
	noop := func(s string) (string, bool) { return s, false }
	if _, err := VerifyMutation(context.Background(), dir, nil, "a.go", noop); err == nil {
		t.Fatal("esperava erro quando testCmd vazio")
	}
	if _, err := VerifyMutation(context.Background(), dir, []string{"echo"}, "../fora.go", noop); err == nil || !strings.Contains(err.Error(), "Path Jail") {
		t.Fatalf("esperava erro Path Jail para targetFile fora, obtive: %v", err)
	}
	res, err := VerifyMutation(context.Background(), dir, []string{"echo"}, "a.go", noop)
	if err != nil || !res.Passed {
		t.Fatalf("esperava Passed=true quando mutateFn nao altera: %+v (%v)", res, err)
	}
	if _, err := resolveSafePath(dir, "../fora"); err == nil || !strings.Contains(err.Error(), "caminho fora do Path Jail") {
		t.Fatalf("esperava erro 'caminho fora do Path Jail', obtive: %v", err)
	}
	if got, err := resolveSafePath(dir, "."); err != nil || got == "" {
		t.Fatalf("esperava sucesso para '.' em resolveSafePath: %q (%v)", got, err)
	}
}

func TestMutacaoSlots(t *testing.T) {
	dir := t.TempDir()
	if err := executeSlotCommands(context.Background(), [][]string{{}}, dir); err != nil {
		t.Fatalf("argv vazio deveria ser ignorado: %v", err)
	}
	if err := executeSlotCommands(context.Background(), [][]string{{"cmd_inexistente_xyz_999"}}, dir); err == nil {
		t.Fatal("esperava erro para comando inexistente")
	}
	root, _ := os.OpenRoot(dir)
	defer root.Close()
	if ok, _, err := evalSlotGates(nil, root, 0); err != nil || !ok {
		t.Fatalf("esperava true, nil para slots nil: %v, %v", ok, err)
	}
	if ok, _, err := evalSlotGates(&dag.Slots{}, root, 0); err != nil || !ok {
		t.Fatalf("esperava true, nil para gates vazios: %v, %v", ok, err)
	}
	if ok, _, err := evalSlotGates(&dag.Slots{Gates: []string{"invalido(("}}, root, 0); err == nil || ok {
		t.Fatalf("esperava false, erro para gate invalido: %v, %v", ok, err)
	}
	closedRoot, _ := os.OpenRoot(dir)
	_ = closedRoot.Close()
	if ok, msg, err := evalSlotGates(&dag.Slots{Gates: []string{"file_exists('a.txt')"}}, closedRoot, 0); err == nil || ok || msg != "" {
		t.Fatalf("esperava false, '', erro em Eval com root fechado: %v, %q, %v", ok, msg, err)
	}
}

func TestMutacaoTokenBudget(t *testing.T) {
	id := novaSessaoID()
	if len(id) != 32 || strings.HasPrefix(id, "sess-") {
		t.Fatalf("sessaoID invalido: %q", id)
	}
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "in.txt"), []byte("0123456789abcdef"), 0o644)
	if got := estimatePromptTokens(dag.Node{ContextIn: []string{"in.txt"}}, dir); got != 4 {
		t.Fatalf("esperava 4 tokens, obtive %d", got)
	}
	if got := budgetPolicy(map[string]any{"policy": ""}); got != "abort" {
		t.Fatalf("esperava abort para policy vazia, obtive %q", got)
	}
	n := dag.Node{ID: "n1", Inputs: map[string]string{"k": "12345678"}, TokenBudget: map[string]any{"max_prompt_tokens": 2}}
	truncateNodeInputs(&n, 0)
	if n.Inputs["k"] != "12345678" {
		t.Fatalf("maxTokens=0 nao deveria truncar: %q", n.Inputs["k"])
	}
	if err := enforceTokenBudget(&n, dir, dir, "s-exato"); err != nil {
		t.Fatalf("est == maxTokens deveria passar: %v", err)
	}
}

func TestMutacaoRecordBudgetTelemetry(t *testing.T) {
	dir := t.TempDir()
	var buf bytes.Buffer
	oldOut := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOut)
	tel := triagem.TokenBudgetTelemetry{NodeID: "n-novo", Action: "pass", EstimatedTokens: 5, MaxPromptTokens: 10}
	recordBudgetTelemetry(dir, "s-nova", tel)
	if strings.Contains(buf.String(), "token_budget: falha") {
		t.Fatalf("nao esperava log de erro em sucesso: %s", buf.String())
	}
	mem, err := triagem.LerMemoria(dir, "s-nova")
	if err != nil || mem.EtapaAtual != "n-novo" {
		t.Fatalf("esperava EtapaAtual=n-novo: %+v (%v)", mem, err)
	}
	arq := filepath.Join(dir, "arq")
	_ = os.WriteFile(arq, []byte("x"), 0o644)
	recordBudgetTelemetry(arq, "s-err", tel)
	if !strings.Contains(buf.String(), "token_budget: falha") {
		t.Fatal("esperava log de falha quando memDir eh arquivo")
	}
}

func TestMutacaoTriggersCronEWatcher(t *testing.T) {
	if !matchCronField("7", 0, 0, 7) || matchCronField("7", 1, 0, 7) || matchCronField("7", 0, 0, 59) {
		t.Fatal("normalizacao de domingo (0/7) falhou")
	}
	if matchCronField("*/0", 0, 0, 59) || matchCronField("*/5", 3, 0, 59) || matchCronField("*/-1", 0, 0, 59) {
		t.Fatal("validacao de step cron falhou")
	}
	if !matchCronField("10-20", 10, 0, 59) || !matchCronField("10-20", 20, 0, 59) || !matchCronField("10-20", 15, 0, 59) || matchCronField("10-20", 5, 0, 59) {
		t.Fatal("validacao de faixa cron falhou")
	}
	if w := NewFileWatcher(t.TempDir(), "a.txt", 0); w.pollInterval != 250*time.Millisecond {
		t.Fatalf("esperava 250ms, obtive %v", w.pollInterval)
	}
	if err := listenCron(context.Background(), "invalido", func() error { return nil }); err == nil {
		t.Fatal("esperava erro para cron invalido em listenCron")
	}
}

func TestMutacaoFileWatcherCriaELogErro(t *testing.T) {
	dir := t.TempDir()
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()
	chamadas := 0
	go func() {
		time.Sleep(15 * time.Millisecond)
		_ = os.WriteFile(filepath.Join(dir, "novo.txt"), []byte("1"), 0o644)
	}()
	_ = NewFileWatcher(dir, "novo.txt", 10*time.Millisecond).Watch(ctx, func(string) { chamadas++ })
	if chamadas != 0 {
		t.Fatalf("criacao inicial nao deveria disparar onChange, chamadas=%d", chamadas)
	}
	var buf bytes.Buffer
	oldOut := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(oldOut)
	alvo := filepath.Join(dir, "watch.txt")
	_ = os.WriteFile(alvo, []byte("v1"), 0o644)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 260*time.Millisecond)
	defer cancel2()
	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = os.WriteFile(alvo, []byte("v2"), 0o644)
		fut := time.Now().Add(2 * time.Second)
		_ = os.Chtimes(alvo, fut, fut)
	}()
	_ = NewTriggerScheduler(dir).Listen(ctx2, dag.Trigger{Type: dag.TriggerFileChange, Path: "watch.txt"}, func() error {
		return errors.New("falha_trg")
	})
	if !strings.Contains(buf.String(), "triggers: falha ao disparar") {
		t.Fatalf("esperava log de erro em Listen, obtive: %s", buf.String())
	}
}

func TestMutacaoWorktreeMergeDuplicado(t *testing.T) {
	src := []byte("package p\n\nfunc Foo() {}\n")
	merged, err := MergeAST(src, src)
	if err != nil {
		t.Fatalf("MergeAST falhou: %v", err)
	}
	if strings.Count(string(merged), "func Foo") != 1 {
		t.Fatalf("func Foo nao deveria duplicar:\n%s", string(merged))
	}
}

func TestMutacaoRunnerHelpersETelemetria(t *testing.T) {
	soma := (TokenUsage{Input: 1, Output: 2, ReportedInput: 10, CostUSD: 0.5}).add(TokenUsage{Input: 3, Output: 4, ReportedInput: 20, CostUSD: 0.25})
	if soma.ReportedInput != 30 || soma.CostUSD != 0.75 {
		t.Fatalf("TokenUsage.add incorreto: %+v", soma)
	}
	if parallelism(0) != runtime.NumCPU() {
		t.Fatalf("parallelism(0) deveria ser NumCPU")
	}
	if errString(nil) != "" || errString(errors.New("boom")) != "boom" {
		t.Fatal("errString incorreto")
	}
	if err := checkRedBlueBoundary(dag.Node{Role: "blue_team", Output: dag.HandoffContract{Artifact: "x_test.go"}}); err == nil {
		t.Fatal("esperava erro para artefato _test.go em blue_team")
	}
	dir := t.TempDir()
	d := &dag.DAG{Name: "p1"}
	rep := Report{Nodes: []NodeReport{{Status: StatusPassed}, {Status: StatusError}}}
	saveRunTelemetry(d, Options{Workdir: dir}, &rep)
	frPath := filepath.Join(dir, ".gemini", "flight_record.json")
	dados, err := os.ReadFile(frPath)
	if err != nil || !strings.Contains(string(dados), `"nodes_failed": 1`) || !strings.Contains(string(dados), `"first_pass_yield": 0.5`) {
		t.Fatalf("flight_record invalido: %s (%v)", string(dados), err)
	}
	repVazio := Report{}
	saveRunTelemetry(d, Options{Workdir: dir}, &repVazio)
	if len(repVazio.Warnings) != 0 {
		t.Fatalf("nao esperava warnings em report vazio: %v", repVazio.Warnings)
	}
}

func TestMutacaoRunnerSlotsEMergeGo(t *testing.T) {
	dir := t.TempDir()
	cmdOK := []string{"go", "version"}
	cmdErr := []string{"cmd_inexistente_xyz_999"}
	r := &run{d: &dag.DAG{Nodes: []dag.Node{{ID: "n1"}}}, opts: Options{Workdir: dir, SessaoID: "s1"}}
	nPre := dag.Node{ID: "n1", Slots: &dag.Slots{Hooks: dag.Hooks{Pre: [][]string{cmdOK}}}}
	if err := r.preExecNode(context.Background(), &nPre, dir); err != nil {
		t.Fatalf("preExecNode com hook valido falhou: %v", err)
	}
	if msg := runPostSlots(context.Background(), dag.Node{Slots: &dag.Slots{Scripts: [][]string{cmdOK}, Hooks: dag.Hooks{Post: [][]string{cmdOK}}}}, dir); msg != "" {
		t.Fatalf("runPostSlots valido falhou: %s", msg)
	}
	if msg := runPostSlots(context.Background(), dag.Node{Slots: &dag.Slots{Scripts: [][]string{cmdErr}}}, dir); !strings.Contains(msg, "slot script falhou") {
		t.Fatalf("esperava erro script, obtive: %q", msg)
	}
	if msg := runPostSlots(context.Background(), dag.Node{Slots: &dag.Slots{Hooks: dag.Hooks{Post: [][]string{cmdErr}}}}, dir); !strings.Contains(msg, "slot hook post falhou") {
		t.Fatalf("esperava erro hook post, obtive: %q", msg)
	}
	iso := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "nota.txt"), []byte("package p\nfunc A() {}\n"), 0o644)
	_ = os.WriteFile(filepath.Join(iso, "nota.txt"), []byte("package p\nfunc B() {}\n"), 0o644)
	r.root, _ = os.OpenRoot(dir)
	defer r.root.Close()
	isoRoot, _ := os.OpenRoot(iso)
	defer isoRoot.Close()
	r.mergeIsolatedGoArtifact(dag.Node{Output: dag.HandoffContract{Artifact: "nota.txt"}}, isoRoot)
	if dados, _ := os.ReadFile(filepath.Join(iso, "nota.txt")); strings.Contains(string(dados), "func A") {
		t.Fatal("nao deveria mesclar arquivo nao .go")
	}
}

func TestMutacaoRunnerSTE100PrompterESchema(t *testing.T) {
	if hasSTE100Tag(nil) || hasSTE100Tag([]string{"outra"}) || !hasSTE100Tag([]string{"ste100"}) || !hasSTE100Tag([]string{"ste-100"}) || !hasSTE100Tag([]string{"asd-ste100"}) {
		t.Fatal("hasSTE100Tag incorreto")
	}
	dir := t.TempDir()
	root, _ := os.OpenRoot(dir)
	defer root.Close()
	if ok, _ := evalSTE100Artifact(dag.Node{}, root); !ok {
		t.Fatal("artefato vazio deveria passar STE100")
	}
	if ok, _ := evalSTE100Artifact(dag.Node{Output: dag.HandoffContract{Artifact: "inexistente.md"}}, root); !ok {
		t.Fatal("artefato inexistente deveria passar STE100")
	}
	if ok, err := (nonIntPrompter{}).Prompt("x"); ok || err == nil {
		t.Fatalf("nonIntPrompter deveria retornar false, err: %v, %v", ok, err)
	}
	_ = os.WriteFile(filepath.Join(dir, "s.json"), []byte(`{"type":"object"}`), 0o644)
	_ = os.WriteFile(filepath.Join(dir, "ste.md"), []byte("Isso provavelmente talvez poderia funcionar perfeitamente."), 0o644)
	d, err := dag.Parse([]byte(`{"name":"t","nodes":[{"id":"n1","type":"script","command":["x"],"routing":{"mode":"manual","model":"none"},"rule_tags":["ste100"],"output":{"artifact":"ste.md","gate":"exit_code == 0"}}]}`))
	if err != nil {
		t.Fatal(err)
	}
	r := &run{d: d, root: root}
	nr := NodeReport{ID: "n1"}
	r.evalNodeGates(&nr, *d.Node("n1"), root, 1)
	if strings.Contains(nr.Error, "ste100 gate reprovou") {
		t.Fatalf("nao deveria rodar STE100 quando gate principal falhou: %s", nr.Error)
	}
}

func TestMutacaoRunnerEvalSchemaCaminhos(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "s.json"), []byte(`{"type":"object"}`), 0o644)
	root, _ := os.OpenRoot(dir)
	defer root.Close()
	r := &run{root: root}
	if ok, _ := r.evalSchema(dag.Node{}, root); !ok {
		t.Fatal("schema vazio deveria retornar true")
	}
	casos := []dag.Node{
		{Output: dag.HandoffContract{Schema: "../fora.json"}},
		{Output: dag.HandoffContract{Schema: "nao_existe.json"}},
		{Output: dag.HandoffContract{Schema: "s.json", Artifact: "../fora.json"}},
		{Output: dag.HandoffContract{Schema: "s.json", Artifact: "nao_existe.json"}},
	}
	for i, n := range casos {
		if ok, _ := r.evalSchema(n, root); ok {
			t.Fatalf("caso %d deveria retornar false", i)
		}
	}
}
