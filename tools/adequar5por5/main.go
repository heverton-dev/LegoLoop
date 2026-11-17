// Comando adequar5por5 completa o nascimento 5/5 dos Processos já gravados em processos/:
// complementa só o que falta (anticorpos, triagem, memória de execução e documento do
// Processo), nunca sobrescreve conteúdo existente e nunca toca em pipeline.dag.json.
// Cada manifesto novo passa pelo schema oficial e pelo linter taxonômico antes de gravar.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"legoloop/internal/dag"
	"legoloop/internal/fabricar"
	"legoloop/internal/gate"
	"legoloop/internal/studiodoc"
	"legoloop/internal/taxonomy"
	"legoloop/internal/triagem"
)

var agora = time.Now().UTC()

func main() {
	raiz := "processos"
	entradas, err := os.ReadDir(raiz)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ler %s: %v\n", raiz, err)
		os.Exit(1)
	}
	for _, e := range entradas {
		if !e.IsDir() {
			continue
		}
		processar(filepath.Join(raiz, e.Name()), e.Name())
	}
}

// processar completa um Processo e imprime cada manifesto gravado ou pendente.
func processar(dir, nome string) {
	d, ids, err := lerDAG(dir)
	if err != nil {
		fmt.Printf("%s: sem DAG utilizável: %v\n", nome, err)
		return
	}
	anticorpos(dir, nome)
	triagemDoProcesso(dir, nome, ids)
	memoria(dir, nome)
	documento(dir, nome, d)
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

// anticorpos grava o anticorpos.json do nascimento se faltar.
func anticorpos(dir, nome string) {
	alvo := filepath.Join(dir, "anticorpos.json")
	if existe(alvo) {
		return
	}
	b, err := fabricar.AnticorposPadrao(nome)
	if err != nil {
		fmt.Printf("%s: anticorpos: %v\n", nome, err)
		return
	}
	gravar(dir, "anticorpos.json", b, "anticorpos.v1.json", nome)
}

// triagemDoProcesso grava a triagem do nascimento ou completa as etapas que faltam,
// preservando o original como .legoloop-bak antes de qualquer reescrita.
func triagemDoProcesso(dir, nome string, ids []string) {
	alvo := filepath.Join(dir, triagem.ArquivoTriagem)
	s := triagem.GerarTriagemPadrao(nome, ids)
	original, err := os.ReadFile(alvo)
	if err == nil {
		var atual triagem.TriagemSpec
		if json.Unmarshal(original, &atual) == nil {
			antes := len(atual.Etapas)
			s = triagem.Completar(atual, nome, ids)
			if len(s.Etapas) == antes {
				return
			}
			preservar(dir, triagem.ArquivoTriagem, original)
		}
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Printf("%s: triagem: %v\n", nome, err)
		return
	}
	gravar(dir, triagem.ArquivoTriagem, b, "triagem.v1.json", nome)
}

// memoria grava a memória de execução inicial se faltar.
func memoria(dir, nome string) {
	alvo := filepath.Join(dir, triagem.ArquivoMemoria)
	if existe(alvo) {
		return
	}
	b, err := json.MarshalIndent(triagem.MemoriaExecucao{Processo: nome, AtualizadoEm: agora}, "", "  ")
	if err != nil {
		fmt.Printf("%s: memória: %v\n", nome, err)
		return
	}
	gravar(dir, triagem.ArquivoMemoria, b, "memoria_execucao.v1.json", nome)
}

// documento grava o documento do Processo se nenhum dialeto existir. O documento sai do
// DAG real e só é gravado se o conteúdo passar no linter taxonômico.
func documento(dir, nome string, d *dag.DAG) {
	if temDialeto(dir) {
		return
	}
	doc := studiodoc.Processo{ID: d.ID, Name: d.Name, Description: d.Description,
		Version: d.Version, Date: d.Date, Trigger: d.Trigger, Nodes: d.Nodes}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Printf("%s: documento: %v\n", nome, err)
		return
	}
	if violacoes, _ := taxonomy.ValidateText(string(b)); len(violacoes) > 0 {
		fmt.Printf("%s: PENDENTE documento do Processo: conteúdo herda %d termo(s) proibido(s) do DAG\n", nome, len(violacoes))
		return
	}
	gravar(dir, nome+".loop.json", b, "doc_processo.v1.json", nome)
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

// preservar guarda o original como .legoloop-bak antes de uma reescrita de complemento.
func preservar(dir, arquivo string, original []byte) {
	_ = os.WriteFile(filepath.Join(dir, arquivo+".legoloop-bak"), original, 0o644)
}

func existe(caminho string) bool {
	_, err := os.Stat(caminho)
	return err == nil
}

// temDialeto procura o documento do Processo em qualquer dialeto gravado.
func temDialeto(dir string) bool {
	achou := false
	_ = filepath.WalkDir(dir, func(p string, e os.DirEntry, err error) error {
		if err == nil && !e.IsDir() && (filepath.Ext(p) == ".json") {
			n := filepath.Base(p)
			if len(n) > 10 && (n[len(n)-10:] == ".loop.json" || n[len(n)-12:] == ".studio.json") {
				achou = true
			}
		}
		return nil
	})
	return achou
}
