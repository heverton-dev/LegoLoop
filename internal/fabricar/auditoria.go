// Package fabricar — auditoria universal agnóstica de alvo. O motor disseca verbo, processo, peça,
// função, documento, binário ou repo externo, aplica o critério binário e emite o veredito com evidência.
package fabricar

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/interview"
	"github.com/heverton-dev/LegoLoop/internal/processo"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

// TipoAlvo classifica o alvo da auditoria universal.
type TipoAlvo string

// Os sete tipos de alvo da auditoria universal.
const (
	AlvoVerbo       TipoAlvo = "verbo"
	AlvoProcesso    TipoAlvo = "processo"
	AlvoPeca        TipoAlvo = "peca"
	AlvoFuncao      TipoAlvo = "funcao"
	AlvoDocumento   TipoAlvo = "documento"
	AlvoBinario     TipoAlvo = "binario"
	AlvoRepoExterno TipoAlvo = "repo-externo"
)

// Vereditos binários da auditoria.
const (
	VereditoAprovado  = "APROVADO"
	VereditoReprovado = "REPROVADO"
)

// Códigos de saída da auditoria: 0 aprovado, 1 reprovado, 2 alvo inexistente.
const (
	ExitAprovado        = 0
	ExitReprovado       = 1
	ExitAlvoInexistente = 2
)

// Identidade do Processo auditoria-atomica: slug do diretório e nome canônico do contrato.
const (
	SlugAuditoriaAtomica = "auditoria-atomica"
	NomeAuditoriaAtomica = "legoloop-" + SlugAuditoriaAtomica
)

// ErrAlvoInexistente: o alvo da auditoria não existe ou está inacessível.
var ErrAlvoInexistente = errors.New("auditoria: alvo inexistente")

// Alvo é o que a auditoria universal disseca.
type Alvo struct {
	Tipo    TipoAlvo
	Nome    string
	Caminho string
}

// ResultadoTeste é a evidência mecânica do critério binário: o comando, a saída e o código de saída.
type ResultadoTeste struct {
	Comando string
	Saida   string
	Exit    int
}

// TesteFunc avalia o alvo no sistema real e devolve a evidência.
type TesteFunc func(Alvo) ResultadoTeste

// Lacuna é o que o alvo não cumpre, com o ticket que fecha a falta.
type Lacuna struct {
	Descricao string
	Ticket    string
}

// Ticket é a unidade de planejamento que fecha uma lacuna da auditoria.
type Ticket struct {
	ID        string
	Descricao string
}

// Relatorio é a entrega da auditoria de um alvo: disseção, lacunas, tickets, refatoração e veredito.
type Relatorio struct {
	Alvo        Alvo
	Dissecao    []string
	Lacunas     []Lacuna
	Tickets     []Ticket
	Refatoracao string
	Veredito    string
	Evidencias  []ResultadoTeste
}

// Auditar disseca o alvo, aplica o critério binário e devolve o relatório com o código de saída.
// Alvo inexistente devolve ErrAlvoInexistente com o código 2 e relatório vazio.
func Auditar(a Alvo, teste TesteFunc) (Relatorio, int, error) {
	info, err := os.Stat(a.Caminho)
	if err != nil {
		return Relatorio{}, ExitAlvoInexistente, fmt.Errorf("%w: %s", ErrAlvoInexistente, a.Caminho)
	}
	rel := Relatorio{Alvo: a}
	rel.Dissecao = dissecaoDe(a, info)
	rel.Lacunas = criteriosDe(a, info, teste)
	if teste != nil {
		ev := teste(a)
		rel.Evidencias = []ResultadoTeste{ev}
		if ev.Exit != 0 {
			rel.Lacunas = append(rel.Lacunas, Lacuna{Descricao: fmt.Sprintf("o critério binário reprovou o alvo (exit %d)", ev.Exit)})
		}
	}
	return rel, fecharRelatorio(&rel), nil
}

// AuditarEGravar audita o alvo e grava o relatório renderizado. Alvo inexistente não grava nada:
// zero artefato órfão.
func AuditarEGravar(a Alvo, teste TesteFunc, saida string) (Relatorio, int, error) {
	rel, code, err := Auditar(a, teste)
	if err != nil {
		return rel, code, err
	}
	if err := EscreverRelatorio(saida, Renderizar(rel)); err != nil {
		return rel, code, err
	}
	return rel, code, nil
}

// EscreverRelatorio grava o conteúdo do relatório, criando os diretórios do caminho.
func EscreverRelatorio(caminho, conteudo string) error {
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	return os.WriteFile(caminho, []byte(conteudo), 0o644)
}

// dissecaoDe descreve o alvo: tipo, caminho e o que a leitura do item mostra.
func dissecaoDe(a Alvo, info os.FileInfo) []string {
	linhas := []string{fmt.Sprintf("Tipo: %s", a.Tipo), fmt.Sprintf("Item: %s", a.Nome), fmt.Sprintf("Caminho: %s", a.Caminho)}
	if info.IsDir() {
		return append(linhas, dissecaoDir(a.Caminho)...)
	}
	return append(linhas, dissecaoArquivo(a.Caminho, info)...)
}

// dissecaoDir lista a contagem e os nomes dos arquivos do diretório do alvo.
func dissecaoDir(dir string) []string {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return []string{"Leitura: indisponível"}
	}
	nomes := make([]string, 0, len(entradas))
	for _, e := range entradas {
		if !e.IsDir() {
			nomes = append(nomes, e.Name())
		}
	}
	return []string{fmt.Sprintf("Arquivos: %d", len(nomes)), fmt.Sprintf("Itens: %s", strings.Join(nomes, ", "))}
}

// dissecaoArquivo mede o arquivo: tamanho, linhas, palavras e SHA-256.
func dissecaoArquivo(caminho string, info os.FileInfo) []string {
	linhas := []string{fmt.Sprintf("Tamanho: %d bytes", info.Size())}
	if b, err := os.ReadFile(caminho); err == nil {
		texto := string(b)
		linhas = append(linhas,
			fmt.Sprintf("Linhas: %d", strings.Count(texto, "\n")+1),
			fmt.Sprintf("Palavras: %d", len(strings.Fields(texto))),
			fmt.Sprintf("SHA-256: %s", sha256Hex(b)))
	}
	return linhas
}

// sha256Hex devolve o resumo SHA-256 do conteúdo em hexadecimal.
func sha256Hex(b []byte) string {
	soma := sha256.Sum256(b)
	return hex.EncodeToString(soma[:])
}

// criteriosDe aplica o critério binário do tipo do alvo e devolve as lacunas encontradas.
func criteriosDe(a Alvo, info os.FileInfo, teste TesteFunc) []Lacuna {
	switch a.Tipo {
	case AlvoVerbo:
		return criteriosVerbo(a, teste)
	case AlvoFuncao:
		return criteriosFuncao(a, info)
	case AlvoDocumento:
		return criteriosDocumento(a, info)
	case AlvoBinario:
		return criteriosBinario(info)
	case AlvoProcesso, AlvoPeca, AlvoRepoExterno:
		return criteriosConteudo(a, info)
	}
	return []Lacuna{{Descricao: fmt.Sprintf("tipo de alvo %q desconhecido", string(a.Tipo))}}
}

// criteriosVerbo é o critério binário do verbo: fonte real, teste unitário e critério mecânico.
func criteriosVerbo(a Alvo, teste TesteFunc) []Lacuna {
	var lacs []Lacuna
	if !temFonte(a.Caminho) {
		lacs = append(lacs, Lacuna{Descricao: "o verbo não tem arquivo de fonte não-teste"})
	}
	if !temTeste(a.Caminho) {
		lacs = append(lacs, Lacuna{Descricao: "o verbo não tem teste unitário com função Test"})
	}
	if teste == nil {
		lacs = append(lacs, Lacuna{Descricao: "o verbo não tem critério binário mecânico"})
	}
	return lacs
}

// temFonte diz se o diretório do verbo tem ao menos um arquivo de fonte não-teste.
func temFonte(dir string) bool {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entradas {
		if strings.HasSuffix(e.Name(), ".go") && !strings.HasSuffix(e.Name(), "_test.go") {
			return true
		}
	}
	return false
}

// temTeste diz se o diretório do verbo tem ao menos um teste unitário com função Test.
func temTeste(dir string) bool {
	entradas, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entradas {
		if !strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		if b, err := os.ReadFile(filepath.Join(dir, e.Name())); err == nil && strings.Contains(string(b), "func Test") {
			return true
		}
	}
	return false
}

// criteriosFuncao exige a função nomeada dentro do arquivo do alvo.
func criteriosFuncao(a Alvo, info os.FileInfo) []Lacuna {
	if info.IsDir() {
		return []Lacuna{{Descricao: "a função precisa de arquivo, não diretório"}}
	}
	b, err := os.ReadFile(a.Caminho)
	if err != nil {
		return []Lacuna{{Descricao: "arquivo da função ilegível"}}
	}
	if a.Nome == "" || !strings.Contains(string(b), "func "+a.Nome) {
		return []Lacuna{{Descricao: fmt.Sprintf("a função %q não aparece no arquivo", a.Nome)}}
	}
	return nil
}

// criteriosDocumento exige conteúdo real no documento.
func criteriosDocumento(a Alvo, info os.FileInfo) []Lacuna {
	if info.IsDir() {
		return []Lacuna{{Descricao: "o documento precisa de arquivo, não diretório"}}
	}
	b, err := os.ReadFile(a.Caminho)
	if err != nil || len(strings.Fields(string(b))) == 0 {
		return []Lacuna{{Descricao: "o documento não tem conteúdo"}}
	}
	return nil
}

// criteriosBinario exige o binário não-vazio.
func criteriosBinario(info os.FileInfo) []Lacuna {
	if info.IsDir() || info.Size() == 0 {
		return []Lacuna{{Descricao: "o binário não tem conteúdo"}}
	}
	return nil
}

// criteriosConteudo exige o item não-vazio: arquivo com tamanho ou diretório com arquivos.
func criteriosConteudo(a Alvo, info os.FileInfo) []Lacuna {
	if !info.IsDir() {
		if info.Size() == 0 {
			return []Lacuna{{Descricao: "o item não tem conteúdo"}}
		}
		return nil
	}
	entradas, err := os.ReadDir(a.Caminho)
	if err != nil || len(entradas) == 0 {
		return []Lacuna{{Descricao: "o diretório não tem arquivos"}}
	}
	return nil
}

// fecharRelatorio numera os tickets das lacunas, define o veredito e devolve o código de saída.
func fecharRelatorio(rel *Relatorio) int {
	for i := range rel.Lacunas {
		id := fmt.Sprintf("TICK-AUD-%02d", i+1)
		rel.Lacunas[i].Ticket = id
		rel.Tickets = append(rel.Tickets, Ticket{ID: id, Descricao: "Fechar a lacuna: " + rel.Lacunas[i].Descricao})
	}
	if len(rel.Lacunas) == 0 {
		rel.Veredito = VereditoAprovado
		return ExitAprovado
	}
	rel.Veredito = VereditoReprovado
	rel.Refatoracao = sugerirRefatoracao(*rel)
	return ExitReprovado
}

// sugerirRefatoracao aponta o primeiro ajuste que fecha a lacuna principal do alvo.
func sugerirRefatoracao(rel Relatorio) string {
	return fmt.Sprintf("Refatore %s (%s): %s. Reavalie com o critério binário depois do ajuste.",
		rel.Alvo.Nome, rel.Alvo.Tipo, rel.Lacunas[0].Descricao)
}

// Renderizar escreve as cinco seções fixas do relatório: Disseção, Lacunas, Tickets, Refatoração e Veredito.
func Renderizar(rel Relatorio) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Relatório de Auditoria — %s (%s)\n\n", rel.Alvo.Nome, rel.Alvo.Tipo)
	secaoDissecao(&b, rel)
	secaoLacunas(&b, rel)
	secaoTickets(&b, rel)
	secaoRefatoracao(&b, rel)
	secaoVeredito(&b, rel)
	return b.String()
}

// secaoDissecao lista a disseção do alvo.
func secaoDissecao(b *strings.Builder, rel Relatorio) {
	b.WriteString("## Disseção\n\n")
	for _, l := range rel.Dissecao {
		fmt.Fprintf(b, "- %s\n", l)
	}
	b.WriteString("\n")
}

// secaoLacunas lista as lacunas com o ticket de cada uma.
func secaoLacunas(b *strings.Builder, rel Relatorio) {
	b.WriteString("## Lacunas\n\n")
	if len(rel.Lacunas) == 0 {
		b.WriteString("Nenhuma lacuna.\n\n")
		return
	}
	for _, l := range rel.Lacunas {
		fmt.Fprintf(b, "- %s (`%s`)\n", l.Descricao, l.Ticket)
	}
	b.WriteString("\n")
}

// secaoTickets lista as unidades de planejamento que fecham as lacunas.
func secaoTickets(b *strings.Builder, rel Relatorio) {
	b.WriteString("## Tickets\n\n")
	if len(rel.Tickets) == 0 {
		b.WriteString("Nenhum ticket.\n\n")
		return
	}
	for _, t := range rel.Tickets {
		fmt.Fprintf(b, "- `%s`: %s\n", t.ID, t.Descricao)
	}
	b.WriteString("\n")
}

// secaoRefatoracao traz a sugestão de refatoração do alvo.
func secaoRefatoracao(b *strings.Builder, rel Relatorio) {
	b.WriteString("## Refatoração\n\n")
	if strings.TrimSpace(rel.Refatoracao) == "" {
		b.WriteString("Nenhuma refatoração sugerida.\n\n")
		return
	}
	fmt.Fprintf(b, "%s\n\n", rel.Refatoracao)
}

// secaoVeredito fecha o relatório com o veredito e as evidências mecânicas.
func secaoVeredito(b *strings.Builder, rel Relatorio) {
	fmt.Fprintf(b, "## Veredito\n\n**%s**\n\n### Evidências\n\n", rel.Veredito)
	if len(rel.Evidencias) == 0 {
		b.WriteString("Nenhuma evidência mecânica registrada.\n\n")
		return
	}
	for _, ev := range rel.Evidencias {
		fmt.Fprintf(b, "- Comando: `%s` — exit %d\n  - Saída: %s\n", ev.Comando, ev.Exit, emUmaLinha(ev.Saida, 200))
	}
	b.WriteString("\n")
}

// emUmaLinha deixa a saída da evidência numa linha só, cortada no limite de runas.
func emUmaLinha(s string, limite int) string {
	um := strings.Join(strings.Fields(s), " ")
	if r := []rune(um); len(r) > limite {
		return string(r[:limite]) + "..."
	}
	return um
}

// VerbosAuditados são os seis verbos do LegoLoop que a auditoria universal examina.
var VerbosAuditados = []string{"matriz", "listar", "refinar", "remover", "compor", "fechar"}

// AuditarVerbos audita os seis verbos sob raiz/internal e devolve os relatórios com os códigos de saída.
func AuditarVerbos(raiz string, teste TesteFunc) ([]Relatorio, []int, error) {
	rels := make([]Relatorio, 0, len(VerbosAuditados))
	codigos := make([]int, 0, len(VerbosAuditados))
	for _, v := range VerbosAuditados {
		a := Alvo{Tipo: AlvoVerbo, Nome: v, Caminho: filepath.Join(raiz, "internal", v)}
		rel, code, err := Auditar(a, teste)
		if err != nil {
			return nil, nil, err
		}
		rels = append(rels, rel)
		codigos = append(codigos, code)
	}
	return rels, codigos, nil
}

// AuditarVerbosEGravar audita os seis verbos e grava o relatório atômico de cada um
// em dir/<verbo>.md, com lacunas, tickets, veredito e evidência mecânica.
func AuditarVerbosEGravar(raiz string, teste TesteFunc, dir string) ([]Relatorio, []int, error) {
	rels, codigos, err := AuditarVerbos(raiz, teste)
	if err != nil {
		return nil, nil, err
	}
	for _, rel := range rels {
		if err := EscreverRelatorio(filepath.Join(dir, rel.Alvo.Nome+".md"), Renderizar(rel)); err != nil {
			return nil, nil, err
		}
	}
	return rels, codigos, nil
}

// TestePadraoEm roda o critério binário real do verbo: go test -count=1 no pacote dele sob a raiz.
func TestePadraoEm(raiz string) TesteFunc {
	return func(a Alvo) ResultadoTeste {
		comando := "go test -count=1 ./internal/" + a.Nome + "/..."
		cmd := exec.Command("go", "test", "-count=1", "./internal/"+a.Nome+"/...")
		cmd.Dir = raiz
		saida, err := cmd.CombinedOutput()
		return ResultadoTeste{Comando: comando, Saida: string(saida), Exit: codigoDeErro(err)}
	}
}

// codigoDeErro converte o erro do comando no código de saída do sistema operacional.
func codigoDeErro(err error) int {
	if err == nil {
		return ExitAprovado
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return ExitReprovado
}

// RelatorioUniversal consolida os relatórios num documento único com o veredito de cada alvo auditado.
func RelatorioUniversal(rels []Relatorio, agora time.Time) string {
	var b strings.Builder
	b.WriteString("# Relatório de Auditoria Universal dos 6 Verbos\n\n")
	fmt.Fprintf(&b, "- **Data:** %s\n- **Processo:** %s\n\n", agora.Format("2006-01-02"), SlugAuditoriaAtomica)
	b.WriteString("| Verbo | Veredito | Código | Lacunas |\n| :--- | :--- | :---: | :---: |\n")
	for _, rel := range rels {
		fmt.Fprintf(&b, "| %s | %s | %d | %d |\n", rel.Alvo.Nome, rel.Veredito, CodigoDe(rel), len(rel.Lacunas))
	}
	b.WriteString("\n")
	for _, rel := range rels {
		b.WriteString(Renderizar(rel))
		b.WriteString("\n---\n\n")
	}
	b.WriteString(apendiceProvas)
	return b.String()
}

// apendiceProvas registra as mutações provadas do Green-Red-Green e o que não foi provado.
const apendiceProvas = `## Mutações provadas

Cada mutação forçada fez o teste alvo falhar e foi restaurada em seguida (Green-Red-Green).

- M1: veredito invertido em fecharRelatorio — TestAlvoArbitrario e TestAuditoriaSeisVerbosCriterioBinario falharam — restaurada.
- M2: etapa e05_veredito removida em etapasAuditoria — TestScaffoldAuditoriaAtomica falhou — restaurada.
- M3: seção Refatoração omitida em Renderizar — TestAlvoArbitrario falhou — restaurada.
- M4: exit do critério binário ignorado em Auditar — TestAuditoriaSeisVerbosCriterioBinario falhou — restaurada.
- M5: erro de alvo inexistente ignorado em AuditarEGravar — TestAlvoInexistenteExitDois falhou — restaurada.
- M6: gravação dos relatórios atômicos suprimida em AuditarVerbosEGravar — TestAuditoriaSeisVerbosCriterioBinario falhou — restaurada.

## Não provado

- "-race" não rodou: ambiente sem CGO.
- Auditoria com LLM real não rodou: os testes usam o critério binário simulado.
- Distribuição universal aos harnesses não foi executada nesta sessão.
- A suíte do repositório tem flakes de concorrência pré-existentes (internal/runner, internal/studio): passam isolados e falham de forma intermitente sob a carga da suíte paralela.
`

// CodigoDe devolve o código de saída do relatório: 0 aprovado, 1 reprovado.
func CodigoDe(rel Relatorio) int {
	if rel.Veredito == VereditoAprovado {
		return ExitAprovado
	}
	return ExitReprovado
}

// ProcessoAuditoriaAtomica monta o Processo da auditoria universal nas cinco etapas do critério binário.
func ProcessoAuditoriaAtomica(now time.Time) (processo.Processo, error) {
	etapas, err := etapasAuditoria()
	if err != nil {
		return processo.Processo{}, err
	}
	return processo.Processo{
		Nome:     NomeAuditoriaAtomica,
		Versao:   "1.0.0",
		Natureza: "software",
		Objetivo: "Auditar qualquer alvo com critério binário e emitir veredito com evidências",
		CriadoEm: now.UTC(),
		Etapas:   etapas,
	}, nil
}

// etapasAuditoria monta as cinco etapas: disseção, lacunas, tickets, refatoração e veredito.
func etapasAuditoria() ([]processo.Etapa, error) {
	defs := []struct {
		id, titulo, tipo, produz string
		depende                  []string
	}{
		{"e01_dissecao", "Dissecar o alvo da auditoria", "dissecar", "out/01_dissecao.md", nil},
		{"e02_lacunas", "Listar lacunas do alvo", "verificar", "out/02_lacunas.md", []string{"e01_dissecao"}},
		{"e03_tickets", "Gerar tickets das lacunas", "planejar", "out/03_tickets.md", []string{"e02_lacunas"}},
		{"e04_refatoracao", "Sugerir refatoração do alvo", "refatorar", "out/04_refatoracao.md", []string{"e03_tickets"}},
		{"e05_veredito", "Emitir veredito binário com evidências", "verificar", "out/05_veredito.md", []string{"e04_refatoracao"}},
	}
	out := make([]processo.Etapa, 0, len(defs))
	for _, d := range defs {
		art := processo.Artefato{Caminho: d.produz, Tipo: "markdown", Promete: true}
		gate, err := processo.GateParaArtefato(art)
		if err != nil {
			return nil, err
		}
		out = append(out, processo.Etapa{
			ID: d.id, Titulo: d.titulo, Tipo: d.tipo, Peca: NomeAuditoriaAtomica,
			Consome: []string{}, Produz: art, Gate: gate, Depende: d.depende,
		})
	}
	return out, nil
}

// MaterializarAuditoriaAtomica grava o Processo auditoria-atomica com o nascimento 5/5 do scaffold:
// manifestos do Processo, skill e Quinteto ISO, tudo pelas projeções — nada escrito à mão.
func MaterializarAuditoriaAtomica(roots Roots, now time.Time) ([]string, error) {
	p, err := ProcessoAuditoriaAtomica(now)
	if err != nil {
		return nil, err
	}
	pr, err := processo.Projetar(p, now)
	if err != nil {
		return nil, err
	}
	files, err := arquivosAuditoriaAtomica(roots, p, pr, now)
	if err != nil {
		return nil, err
	}
	plano := MaterializePlan{Name: p.Nome, Roots: roots, Files: files, Force: true}
	if _, err := Materialize(context.Background(), plano, OSFS()); err != nil {
		return nil, err
	}
	return caminhosDe(files), nil
}

// caminhosDe devolve os caminhos gravados, na ordem do plano.
func caminhosDe(files []PlannedFile) []string {
	out := make([]string, len(files))
	for i, f := range files {
		out[i] = f.Path
	}
	return out
}

// arquivosAuditoriaAtomica monta o plano de gravação do Processo: nascimento 5/5, skill e Quinteto ISO.
func arquivosAuditoriaAtomica(roots Roots, p processo.Processo, pr *processo.Projecoes, now time.Time) ([]PlannedFile, error) {
	mem, err := memoriaDe(p)
	if err != nil {
		return nil, err
	}
	ac, err := anticorposDe(p)
	if err != nil {
		return nil, err
	}
	intent := intentAuditoriaAtomica(p.Nome, now)
	docs := func(sub, file string) string { return filepath.Join(roots.Docs, "docs", sub, file) }
	dir := filepath.Join(roots.Vault, caminhos.DirProcessos, SlugAuditoriaAtomica)
	files := []PlannedFile{
		{Kind: "skill", Path: filepath.Join(roots.Vault, "skills", p.Nome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Kind: "studio", Path: filepath.Join(dir, SlugAuditoriaAtomica+caminhos.SufixoStudio), Data: pr.Studio},
		{Kind: "dag", Path: filepath.Join(dir, caminhos.ArquivoDAG), Data: pr.DAG},
		{Kind: "intent", Path: docs("intents", intent.Name), Data: intent.Data},
		{Kind: "spec", Path: docs("specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Kind: "plan", Path: docs("plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Kind: "route", Path: docs("routes", pr.RouteNome), Data: pr.Route},
		{Kind: "dag_iso", Path: docs("dags", pr.DAGNome), Data: pr.DAG},
		{Kind: "triagem", Path: filepath.Join(dir, triagem.ArquivoTriagem), Data: pr.Triagem},
		{Kind: "memoria", Path: filepath.Join(dir, triagem.ArquivoMemoria), Data: mem},
		{Kind: "anticorpos", Path: filepath.Join(dir, "anticorpos.json"), Data: ac},
		{Kind: "manual", Path: docs("manuals", pr.ManualNome), Data: []byte(pr.Manual)},
		{Kind: "livro", Path: docs("books", filepath.Join(p.Nome, "livro.json")), Data: pr.Livro},
	}
	for _, f := range pr.Fichas {
		files = append(files, PlannedFile{Kind: "fichas_tecnicas", Path: docs("fichas", filepath.Join(p.Nome, f.EtapaID+".md")), Data: []byte(f.Conteudo)})
	}
	return files, nil
}

// intentAuditoriaAtomica renderiza o INTENT do Processo pelo template canônico da entrevista.
func intentAuditoriaAtomica(nome string, now time.Time) namedFile {
	date := now.Format("2006-01-02")
	slug := strings.TrimPrefix(nome, namePrefix)
	file := fmt.Sprintf("%s_INTENT_%s.md", date, upperSnake(slug))
	tema := "Auditoria universal agnóstica de alvo"
	doc := interview.RenderIntent(interview.IntentView{
		Theme: tema, Date: date, SessionID: "fabricar-" + slug, Confidence: 100, Natureza: "software",
		Mirror:  espelhoAuditoria(),
		SpecRel: interview.SpecRelFor(file, tema, date),
	})
	return namedFile{Name: file, Data: []byte(doc)}
}

// espelhoAuditoria traz o espelho socrático do Processo de auditoria universal.
func espelhoAuditoria() *interview.Mirror {
	return &interview.Mirror{
		Outcome: "Veredito binário com evidências para qualquer alvo auditado",
		User:    "Operador do LegoLoop que audita verbos, processos, peças, funções, documentos, binários e repos externos",
		WhyNow:  "A auditoria universal é o contrato do Épico 08 do Mega Plano",
		Success: "Todo alvo auditado recebe veredito APROVADO ou REPROVADO com evidência mecânica",
		Constraints: []string{
			"Critério binário com código de saída 0, 1 ou 2",
			"Zero artefato órfão para alvo inexistente",
		},
		OutOfScope: []string{"Correção automática do alvo reprovado"},
	}
}
