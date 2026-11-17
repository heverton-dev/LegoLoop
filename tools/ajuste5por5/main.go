// Comando ajuste5por5 resolve as pendências finais do nascimento 5/5 em processos/:
// regenera triagens antigas que não validam contra triagem.v1, completa os manifestos
// que faltam com nome taxonômico seguro e grava o documento do Processo com nós só de
// identificação quando o DAG herda termos proibidos. Nunca toca em pipeline.dag.json.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"legoloop/internal/dag"
	"legoloop/internal/fabricar"
	"legoloop/internal/gate"
	"legoloop/internal/studiodoc"
	"legoloop/internal/taxonomy"
	"legoloop/internal/triagem"
)

// casos lista os Processos pendentes e o nome taxonômico seguro de cada um.
var casos = []struct{ dir, nome string }{
	{"processos/fronteira-ai-driven", "fronteira-ai-driven"},
	{"processos/legoloop-design-ui", "legoloop-design-ui"},
	{"processos/mega-pipeline", "mega-processo"},
	{"processos/legoloop-swarm", "legoloop-swarm"},
}

func main() {
	for _, c := range casos {
		d, ids, err := lerDAG(c.dir)
		if err != nil {
			fmt.Printf("%s: sem DAG: %v\n", c.dir, err)
			continue
		}
		triagemNova(c.dir, c.nome, ids)
		manifestos(c.dir, c.nome)
		documento(d, c.dir, c.nome)
	}
}

// lerDAG parseia o DAG do Processo e devolve o grafo e os IDs dos nós.
func lerDAG(dir string) (*dag.DAG, []string, error) {
	b, err := os.ReadFile(filepath.Join(dir, "pipeline.dag.json"))
	if err != nil {
		return nil, nil, err
	}
	d, err := dag.Parse(b)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]string, len(d.Nodes))
	for i, n := range d.Nodes {
		ids[i] = n.ID
	}
	return d, ids, nil
}

// triagemNova regenera a triagem no formato triagem.v1, preservando o original.
func triagemNova(dir, nome string, ids []string) {
	s := triagem.GerarTriagemPadrao(nome, ids)
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Printf("%s: triagem: %v\n", nome, err)
		return
	}
	if original, err := os.ReadFile(filepath.Join(dir, triagem.ArquivoTriagem)); err == nil {
		_ = os.WriteFile(filepath.Join(dir, triagem.ArquivoTriagem+".legoloop-bak"), original, 0o644)
	}
	gravar(dir, triagem.ArquivoTriagem, b, "triagem.v1.json", nome)
}

// manifestos garante anticorpos.json e memória de execução do nascimento.
func manifestos(dir, nome string) {
	alvo := filepath.Join(dir, "anticorpos.json")
	if _, err := os.Stat(alvo); err != nil {
		if b, err := fabricar.AnticorposPadrao(nome); err == nil {
			gravar(dir, "anticorpos.json", b, "anticorpos.v1.json", nome)
		} else {
			fmt.Printf("%s: anticorpos: %v\n", nome, err)
		}
	}
	mem := filepath.Join(dir, triagem.ArquivoMemoria)
	if _, err := os.Stat(mem); err == nil {
		return
	}
	b, err := json.MarshalIndent(triagem.MemoriaExecucao{Processo: nome}, "", "  ")
	if err != nil {
		fmt.Printf("%s: memória: %v\n", nome, err)
		return
	}
	gravar(dir, triagem.ArquivoMemoria, b, "memoria_execucao.v1.json", nome)
}

// documento grava o documento do Processo com nós só de identificação quando
// nenhum dialeto existe, para não herdar termos proibidos dos rótulos do DAG.
func documento(d *dag.DAG, dir, nome string) {
	if temDialeto(dir) {
		return
	}
	doc := studiodoc.Processo{ID: d.ID, Name: d.Name, Description: d.Description,
		Version: d.Version, Date: d.Date, Nodes: idNos(d)}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Printf("%s: documento: %v\n", nome, err)
		return
	}
	if violacoes, _ := taxonomy.ValidateText(string(b)); len(violacoes) > 0 {
		fmt.Printf("%s: documento ainda herda %d termo(s); nada gravado\n", nome, len(violacoes))
		return
	}
	gravar(dir, nome+".loop.json", b, "doc_processo.v1.json", nome)
}

// idNos reduz os nós do DAG ao campo de identificação, sem rótulos herdados.
func idNos(d *dag.DAG) []dag.Node {
	nos := make([]dag.Node, len(d.Nodes))
	for i, n := range d.Nodes {
		nos[i] = dag.Node{ID: n.ID}
	}
	return nos
}

// gravar valida o manifesto contra o schema oficial e o linter taxonômico antes do disco.
func gravar(dir, arquivo string, b []byte, esquema, nome string) {
	schema, err := os.ReadFile(filepath.Join("schemas", esquema))
	if err != nil {
		fmt.Printf("%s: schema %s: %v\n", nome, esquema, err)
		return
	}
	if !gate.ValidateJSONSchema(schema, b) {
		fmt.Printf("%s: %s não valida contra %s; nada gravado\n", nome, arquivo, esquema)
		return
	}
	if violacoes, _ := taxonomy.ValidateText(string(b)); len(violacoes) > 0 {
		fmt.Printf("%s: %s herda termo(s) proibido(s); nada gravado\n", nome, arquivo)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, arquivo), b, 0o644); err != nil {
		fmt.Printf("%s: gravar %s: %v\n", nome, arquivo, err)
		return
	}
	fmt.Printf("%s: gravado %s\n", nome, arquivo)
}

// temDialeto procura o documento do Processo em qualquer dialeto gravado.
func temDialeto(dir string) bool {
	achou := false
	_ = filepath.WalkDir(dir, func(p string, e os.DirEntry, err error) error {
		if err == nil && !e.IsDir() {
			n := e.Name()
			if len(n) > 10 && (n[len(n)-10:] == ".loop.json" || n[len(n)-12:] == ".studio.json") {
				achou = true
			}
		}
		return nil
	})
	return achou
}
