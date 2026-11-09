package fabricar

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"legoloop/internal/interview"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

func init() {
	ste.Default.Register("fabricar.lacuna_sem_peca", "A etapa %s não tem peça apta na matriz. Nada foi gravado. Dê um perfil a uma peça com legoloop matriz perfil definir, ou mude o objetivo.")
	ste.Default.Register("fabricar.lacunas_sem_peca", "As etapas %s não têm peça apta na matriz. Nada foi gravado. Dê um perfil a uma peça com legoloop matriz perfil definir, ou mude o objetivo.")
}

// ErrLacuna: uma etapa pedida ficou sem peça apta. A lacuna é decisão do usuário: o motor não inventa peça.
var ErrLacuna = processo.ErrLacuna

// namedFile é um documento já renderizado, com o nome que terá em docs/.
type namedFile struct {
	Name string
	Data []byte
}

// projetar é a Camada 4. O Processo nasce das necessidades e da Fusão (ProcessoDe), e SPEC, PLAN, ROUTE,
// DAG, loop e skill são projeções dele, conferidas entre si por processo.Projetar. Etapa sem peça para aqui:
// o erro é ErrLacuna e o Processo volta junto, para o dry-run montar o relatório da lacuna.
func projetar(name, objetivo string, c consulta, now time.Time) (processo.Processo, *processo.Projecoes, error) {
	p, err := ProcessoDe(name, objetivo, c.needs, c.m, now)
	if err != nil {
		return processo.Processo{}, nil, err
	}
	if l := lacunasDe(p); len(l) > 0 {
		return p, nil, fmt.Errorf("%w: %s", ErrLacuna, msgLacuna(l))
	}
	pr, err := processo.Projetar(p, now)
	if err != nil {
		return processo.Processo{}, nil, err
	}
	return p, pr, nil
}

// msgLacuna é a orientação STE-PT para as etapas sem peça: a mesma no erro e no aviso do cartão.
func msgLacuna(etapas []string) string {
	chave := "fabricar.lacuna_sem_peca"
	if len(etapas) > 1 {
		chave = "fabricar.lacunas_sem_peca"
	}
	return ste.Msg(chave, strings.Join(etapas, ", "))
}

// relatorioDaLacuna é o dry-run com lacuna: nada seria gravado, mas o Report e o cartão saem com a lacuna
// destacada, para o usuário ver a etapa e o que fazer. O chamador devolve o resultado junto com ErrLacuna.
func relatorioDaLacuna(d Deps, name string, clar Clarification, c consulta, p processo.Processo, now time.Time) *Result {
	plano := Plano{Lacunas: lacunasDe(p)}
	for _, h := range d.Harnesses {
		plano.Detectados = append(plano.Detectados, h.Nome)
	}
	warns := append(append(clar.Warnings(), c.avisos...), c.m.Warnings...)
	return &Result{CanonicalName: name, Slug: strings.TrimPrefix(name, namePrefix), CreatedAt: now,
		CuratedPieces: c.m.Used, Warnings: warns, Report: buildReport(d, clar, c, p, nil, plano)}
}

// lacunasDe devolve o nome das etapas sem peça e sem processo composto.
func lacunasDe(p processo.Processo) []string {
	var out []string
	for _, e := range p.Etapas {
		if e.Peca == "" && e.Compoe == "" {
			nome := e.Tipo
			if nome == "" {
				nome = e.Titulo
			}
			out = append(out, nome)
		}
	}
	return out
}

// objetivoDoProcesso é o resultado esperado do espelho; sem ele, o objetivo do pedido.
func objetivoDoProcesso(req Request, clar Clarification) string {
	if o := strings.TrimSpace(clar.Mirror.Outcome); o != "" {
		return o
	}
	return objetivoDe(req)
}

// intentDe renderiza o INTENT da entrevista com a natureza do processo, para que ele e o processo concordem.
func intentDe(name string, clar Clarification, natureza string, now time.Time) namedFile {
	date := now.Format("2006-01-02")
	slug := strings.TrimPrefix(name, namePrefix)
	file := fmt.Sprintf("%s_INTENT_%s.md", date, upperSnake(slug))
	doc := interview.RenderIntent(interview.IntentView{
		Theme: clar.Theme, Date: date, SessionID: "fabricar-" + slug, Confidence: clar.Confidence,
		Natureza: natureza, Mirror: &clar.Mirror, SpecRel: interview.SpecRelFor(file, clar.Theme, date),
	})
	return namedFile{Name: file, Data: []byte(doc)}
}

func upperSnake(slug string) string { return strings.ToUpper(strings.ReplaceAll(slug, "-", "_")) }

// plannedFiles monta a lista final: no cofre, a skill, o loop, o DAG e o processo.json; em docs/, o INTENT, as
// projeções SPEC, PLAN e ROUTE e o TEMPLATE. Tudo é gravado junto, numa única operação atômica.
func plannedFiles(roots Roots, p processo.Processo, pr *processo.Projecoes, intent namedFile) ([]PlannedFile, error) {
	pj, err := p.Marshal()
	if err != nil {
		return nil, err
	}
	docs := func(sub, file string) string { return filepath.Join(roots.Docs, "docs", sub, file) }
	vault := func(parts ...string) string { return filepath.Join(append([]string{roots.Vault}, parts...)...) }
	return []PlannedFile{
		{Kind: "skill", Path: vault("skills", p.Nome, "SKILL.md"), Data: []byte(pr.Skill)},
		{Kind: "loop", Path: vault("loops", p.Nome+".loop.json"), Data: pr.Loop},
		{Kind: "dag", Path: vault("loops", p.Nome+".dag.json"), Data: pr.DAG},
		{Kind: "processo", Path: processo.CaminhoNoCofre(roots.Vault, p.Nome), Data: pj},
		{Kind: "intent", Path: docs("intents", intent.Name), Data: intent.Data},
		{Kind: "spec", Path: docs("specs", pr.SpecNome), Data: []byte(pr.Spec)},
		{Kind: "plan", Path: docs("plans", pr.PlanNome), Data: []byte(pr.Plan)},
		{Kind: "route", Path: docs("routes", pr.RouteNome), Data: pr.Route},
		{Kind: "dag_iso", Path: docs("dags", pr.DAGNome), Data: pr.DAG},
		{Kind: "template", Path: docs("templates", templateFileName(p.Nome)), Data: []byte(EmitTemplate(p))},
	}, nil
}
