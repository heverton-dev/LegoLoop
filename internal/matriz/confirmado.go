package matriz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"legoloop/internal/processo"
)

// confirmadoFile guarda os perfis que o usuário confirmou ou definiu, por hash do conteúdo. Fica fora da matriz
// de propósito: a matriz é regerada, e a palavra do usuário não pode se perder junto.
const confirmadoFile = "perfil.confirmado.json"

var (
	// ErrMatrizAusente: ainda não existe matriz no cofre.
	ErrMatrizAusente = errors.New("matriz: o cofre ainda não tem matriz")
	// ErrPecaNaoEncontrada: nenhum nome nem caminho da matriz casa com o pedido.
	ErrPecaNaoEncontrada = errors.New("matriz: peça não encontrada")
	// ErrPecaAmbigua: o nome casa com conteúdos diferentes; é preciso informar o caminho.
	ErrPecaAmbigua = errors.New("matriz: o nome casa com mais de uma peça")
	// ErrSemPerfil: a peça não tem perfil para confirmar; use definir.
	ErrSemPerfil = errors.New("matriz: a peça não tem perfil")
)

type confirmacao struct {
	Tipo   string      `json:"tipo"`
	Hash   string      `json:"hash"`
	Nome   string      `json:"nome"`
	Perfil PerfilDeUso `json:"perfil"`
}

type arquivoConfirmado struct {
	Versao      int           `json:"versao"`
	Confirmacao []confirmacao `json:"confirmacoes"`
}

// lerConfirmados devolve os perfis confirmados por (tipo, hash). Sem arquivo, não há nenhum.
func lerConfirmados(vault string) (map[string]confirmacao, error) {
	out := map[string]confirmacao{}
	b, err := os.ReadFile(filepath.Join(vault, confirmadoFile))
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	var a arquivoConfirmado
	if err := json.Unmarshal(b, &a); err != nil {
		return nil, fmt.Errorf("matriz: %s inválido: %w", confirmadoFile, err)
	}
	for _, c := range a.Confirmacao {
		out[chaveHash(c.Tipo, c.Hash)] = c
	}
	return out, nil
}

func arquivoDasConfirmacoes(vault string, m map[string]confirmacao) processo.File {
	a := arquivoConfirmado{Versao: 1}
	for _, k := range sortedKeys(m) {
		a.Confirmacao = append(a.Confirmacao, m[k])
	}
	return json2(filepath.Join(vault, confirmadoFile), "confirmado", a)
}

// Confirmar aceita o perfil atual da peça (e das suas cópias idênticas) como certo. A origem vira "confirmado".
func Confirmar(vault, alvo string) (Entrada, error) {
	reg, idx, err := localizar(vault, alvo)
	if err != nil {
		return Entrada{}, err
	}
	p := reg.Entradas[idx[0]].Perfil
	if p == nil {
		return Entrada{}, fmt.Errorf("%w: %s. Defina o perfil com legoloop matriz perfil definir", ErrSemPerfil, alvo)
	}
	return gravarConfirmacao(vault, reg, idx, *p)
}

// Definir grava o perfil que o usuário escolheu. Todo termo precisa existir no vocabulário, na categoria certa.
func Definir(vault, alvo string, p PerfilDeUso) (Entrada, error) {
	reg, idx, err := localizar(vault, alvo)
	if err != nil {
		return Entrada{}, err
	}
	voc, err := LoadVocab(vault)
	if err != nil {
		return Entrada{}, err
	}
	if err := validarPerfil(voc, p); err != nil {
		return Entrada{}, err
	}
	if atual := reg.Entradas[idx[0]].Perfil; atual != nil {
		p.Consome, p.Produz, p.QuandoNaoUsar = atual.Consome, atual.Produz, atual.QuandoNaoUsar
	}
	return gravarConfirmacao(vault, reg, idx, p)
}

func validarPerfil(voc Vocab, p PerfilDeUso) error {
	if len(p.Etapas)+len(p.Naturezas)+len(p.CasosDeUso) == 0 {
		return fmt.Errorf("%w: informe ao menos um termo", ErrVocabulario)
	}
	for _, g := range []struct {
		cat    Categoria
		termos []string
	}{{CatEtapa, p.Etapas}, {CatNatureza, p.Naturezas}, {CatCasoDeUso, p.CasosDeUso}} {
		for _, t := range g.termos {
			if !voc.Has(g.cat, t) {
				return fmt.Errorf("%w: %q não é %s do vocabulário. Veja os termos com legoloop matriz vocab listar", ErrVocabulario, t, g.cat)
			}
		}
	}
	return nil
}

// localizar lê a matriz e acha as entradas da peça pelo caminho ou pelo nome. Cópias idênticas (mesmo tipo e
// hash) contam como uma peça só; nomes iguais com conteúdos diferentes são ambíguos.
func localizar(vault, alvo string) (*Registro, []int, error) {
	reg, err := LerRegistro(filepath.Join(vault, "registry_matrix.json"))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, fmt.Errorf("%w. Rode legoloop matriz gerar", ErrMatrizAusente)
	}
	if err != nil {
		return nil, nil, err
	}
	var idx []int
	for i, e := range reg.Entradas {
		if e.Ref == alvo || e.Nome == alvo {
			idx = append(idx, i)
		}
	}
	if len(idx) == 0 {
		return nil, nil, fmt.Errorf("%w: %s", ErrPecaNaoEncontrada, alvo)
	}
	return reg, idx, ambiguidade(reg, idx, alvo)
}

func ambiguidade(reg *Registro, idx []int, alvo string) error {
	vistos := map[string]bool{}
	var refs []string
	for _, i := range idx {
		e := reg.Entradas[i]
		vistos[chaveHash(e.Tipo, e.Hash)] = true
		refs = append(refs, e.Ref)
	}
	if len(vistos) == 1 {
		return nil
	}
	return fmt.Errorf("%w: %s. Informe o caminho: %s", ErrPecaAmbigua, alvo, strings.Join(refs, ", "))
}

// gravarConfirmacao grava, numa só operação atômica, o arquivo de confirmações, a matriz geral e a matriz do tipo.
func gravarConfirmacao(vault string, reg *Registro, idx []int, p PerfilDeUso) (Entrada, error) {
	p.Origem = OrigemConfirmado
	alvo := reg.Entradas[idx[0]]
	chave := chaveHash(alvo.Tipo, alvo.Hash)
	conf, err := lerConfirmados(vault)
	if err != nil {
		return Entrada{}, err
	}
	conf[chave] = confirmacao{Tipo: alvo.Tipo, Hash: alvo.Hash, Nome: alvo.Nome, Perfil: p}
	aplicarNaLista(reg.Entradas, chave, p)
	tipo, err := arquivoDoTipo(vault, alvo.Tipo, chave, p)
	if err != nil {
		return Entrada{}, err
	}
	files := []processo.File{arquivoDasConfirmacoes(vault, conf), json2(filepath.Join(vault, "registry_matrix.json"), "registro", reg), tipo}
	st := processo.Store{Roots: processo.Roots{Vault: vault}}
	if err := st.Write(context.Background(), files, processo.WriteOpts{Force: true}); err != nil {
		return Entrada{}, err
	}
	return reg.Entradas[idx[0]], nil
}

func aplicarNaLista(es []Entrada, chave string, p PerfilDeUso) {
	for i := range es {
		if chaveHash(es[i].Tipo, es[i].Hash) == chave {
			c := p
			es[i].Perfil, es[i].Retentar = &c, false
		}
	}
}

func arquivoDoTipo(vault, tipo, chave string, p PerfilDeUso) (processo.File, error) {
	path := filepath.Join(vault, "matrix", tipo+".json")
	b, err := os.ReadFile(path)
	if err != nil {
		return processo.File{}, fmt.Errorf("matriz: a matriz por tipo %s não pôde ser lida: %w", tipo, err)
	}
	var a arquivoPorTipo
	if err := json.Unmarshal(b, &a); err != nil {
		return processo.File{}, fmt.Errorf("matriz: %s inválido: %w", path, err)
	}
	aplicarNaLista(a.Entradas, chave, p)
	return json2(path, "tipo", a), nil
}

// invalidadas lista as peças cuja confirmação não vale mais: o nome continua na matriz, mas o conteúdo mudou.
func invalidadas(conf map[string]confirmacao, entradas []Entrada) []string {
	hashes, nomes := map[string]bool{}, map[string]bool{}
	for _, e := range entradas {
		hashes[chaveHash(e.Tipo, e.Hash)] = true
		nomes[e.Tipo+"|"+e.Nome] = true
	}
	var out []string
	for k, c := range conf {
		msg := c.Nome + ": a confirmação do perfil perdeu a validade porque o conteúdo mudou. A peça foi catalogada de novo. Confirme de novo com legoloop matriz perfil confirmar"
		if !hashes[k] && nomes[c.Tipo+"|"+c.Nome] && !slices.Contains(out, msg) {
			out = append(out, msg)
		}
	}
	sort.Strings(out)
	return out
}
