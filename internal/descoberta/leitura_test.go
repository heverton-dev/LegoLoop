package descoberta_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/descoberta"
	"legoloop/internal/matriz"
)

type provedorTeste struct {
	so       string
	unidades []descoberta.Unidade
}

func (p provedorTeste) SO() string                           { return p.so }
func (p provedorTeste) ListarUnidades() []descoberta.Unidade { return p.unidades }

type cenarioLeitura struct {
	origem    string
	home      string
	fixa      string
	removivel string
	prov      descoberta.Provedor
}

func montarCenarioLeitura(t *testing.T) cenarioLeitura {
	t.Helper()
	origem := t.TempDir()
	home := filepath.Join(origem, "home")
	fixa := filepath.Join(origem, "unidade-fixa")
	removivel := filepath.Join(origem, "unidade-removivel")
	gravarArquivoTeste(t, filepath.Join(home, ".keep"), "home\n")
	gravarProjetoFixo(t, filepath.Join(fixa, "proj-alpha"))
	gravarProjetoRemovivel(t, filepath.Join(removivel, "proj-pen"))
	prov := provedorTeste{
		so: "linux",
		unidades: []descoberta.Unidade{
			{Caminho: fixa, Tipo: descoberta.TipoFixa},
			{Caminho: removivel, Tipo: descoberta.TipoRemovivel},
		},
	}
	return cenarioLeitura{origem: origem, home: home, fixa: fixa, removivel: removivel, prov: prov}
}

func gravarProjetoFixo(t *testing.T, proj string) {
	t.Helper()
	gravarArquivoTeste(t, filepath.Join(proj, "go.mod"), "module alpha\n\ngo 1.24\n")
	skill := "---\ncasos_de_uso: [tdd]\netapas: [construir]\nnaturezas: [software]\n---\n# Revisor\n"
	gravarArquivoTeste(t, filepath.Join(proj, "skills", "revisor", "SKILL.md"), skill)
	gravarArquivoTeste(t, filepath.Join(proj, "scripts", "checar.sh"), "#!/bin/sh\necho ok\n")
	dag := `{"id":"processo-alpha","name":"Processo Alpha","nodes":[]}`
	gravarArquivoTeste(t, filepath.Join(proj, "processos", "processo-alpha", "dag.json"), dag)
}

func gravarProjetoRemovivel(t *testing.T, proj string) {
	t.Helper()
	gravarArquivoTeste(t, filepath.Join(proj, "go.mod"), "module pen\n\ngo 1.24\n")
	skill := "---\ncasos_de_uso: [tdd]\netapas: [construir]\nnaturezas: [software]\n---\n# Removivel\n"
	gravarArquivoTeste(t, filepath.Join(proj, "skills", "pen-skill", "SKILL.md"), skill)
}

func gravarArquivoTeste(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatalf("criar diretorio %s: %v", caminho, err)
	}
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatalf("gravar %s: %v", caminho, err)
	}
}

func hashArvore(t *testing.T, raiz string) string {
	t.Helper()
	var caminhos []string
	err := filepath.WalkDir(raiz, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(raiz, p)
		if err != nil {
			return err
		}
		if rel != "." {
			caminhos = append(caminhos, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("percorrer arvore %s: %v", raiz, err)
	}
	sort.Strings(caminhos)
	return resumirCaminhos(t, raiz, caminhos)
}

func resumirCaminhos(t *testing.T, raiz string, caminhos []string) string {
	t.Helper()
	h := sha256.New()
	for _, rel := range caminhos {
		full := filepath.Join(raiz, filepath.FromSlash(rel))
		info, err := os.Stat(full)
		if err != nil {
			t.Fatalf("stat %s: %v", full, err)
		}
		h.Write([]byte(rel + "\x00"))
		if info.IsDir() {
			h.Write([]byte("D\x00"))
			continue
		}
		dados, err := os.ReadFile(full)
		if err != nil {
			t.Fatalf("ler %s: %v", full, err)
		}
		h.Write([]byte("F\x00"))
		h.Write(dados)
		h.Write([]byte("\x00"))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func executarDescobertaEMatriz(t *testing.T, cen cenarioLeitura, cofre string) ([]string, []descoberta.Recusa, *matriz.Summary) {
	t.Helper()
	raizes, recusadas := descoberta.Raizes(cen.home, nil, cen.prov, false)
	scanner := assets.NewScanner(assets.Config{
		Home:             cen.home,
		AutoDiscover:     true,
		ProvedorUnidades: cen.prov,
		TrustRemovable:   false,
		LookPath:         func(string) (string, error) { return "", os.ErrNotExist },
		Getenv:           func(string) string { return "" },
	})
	deps := matriz.Deps{
		Scan: func(ctx context.Context) ([]assets.Item, error) {
			return scanner.Scan(ctx).Items, nil
		},
		Catalogar: matriz.Catalogador{Vocab: matriz.DefaultVocab(), Brain: nil},
		Home:      cen.home,
		Now:       func() time.Time { return time.Date(2026, 11, 21, 12, 0, 0, 0, time.UTC) },
	}
	sum, err := matriz.Generate(context.Background(), matriz.Options{Vault: cofre}, deps)
	if err != nil {
		t.Fatalf("gerar matriz sem cerebro: %v", err)
	}
	return raizes, recusadas, sum
}

func validarRelatorioDescoberta(t *testing.T, cen cenarioLeitura, raizes []string, recusadas []descoberta.Recusa) {
	t.Helper()
	if !slices.Contains(raizes, cen.fixa) || !slices.Contains(raizes, cen.home) {
		t.Fatalf("raizes = %v, esperado conter fixa %q e home %q", raizes, cen.fixa, cen.home)
	}
	if slices.Contains(raizes, cen.removivel) {
		t.Fatalf("unidade removivel %q nao pode estar em raizes: %v", cen.removivel, raizes)
	}
	if len(recusadas) != 1 || recusadas[0].Caminho != cen.removivel || strings.TrimSpace(recusadas[0].Motivo) == "" {
		t.Fatalf("recusadas = %+v, esperado [%s com motivo]", recusadas, cen.removivel)
	}
}

func validarMatrizesIguais(t *testing.T, cofreA, cofreB string, sumA, sumB *matriz.Summary) {
	t.Helper()
	if sumA.Total < 2 || sumA.Total != sumB.Total {
		t.Fatalf("totais divergentes ou insuficientes: A=%d B=%d", sumA.Total, sumB.Total)
	}
	if !slices.Equal(sumA.Arquivos, sumB.Arquivos) {
		t.Fatalf("arquivos gerados diferem: A=%v B=%v", sumA.Arquivos, sumB.Arquivos)
	}
	for _, rel := range sumA.Arquivos {
		a, errA := os.ReadFile(filepath.Join(cofreA, filepath.FromSlash(rel)))
		b, errB := os.ReadFile(filepath.Join(cofreB, filepath.FromSlash(rel)))
		if errA != nil || errB != nil {
			t.Fatalf("ler arquivo %s: errA=%v errB=%v", rel, errA, errB)
		}
		if !bytes.Equal(a, b) {
			t.Fatalf("arquivo %s difere entre execucoes", rel)
		}
		if strings.Contains(string(a), "pen-skill") {
			t.Fatalf("peca de unidade removivel vazou para %s", rel)
		}
	}
}

func TestDescobertaSomenteLeitura(t *testing.T) {
	cen := montarCenarioLeitura(t)
	hashAntes := hashArvore(t, cen.origem)

	cofre1, cofre2 := t.TempDir(), t.TempDir()
	raizes1, recusadas1, sum1 := executarDescobertaEMatriz(t, cen, cofre1)
	raizes2, recusadas2, sum2 := executarDescobertaEMatriz(t, cen, cofre2)

	hashDepois := hashArvore(t, cen.origem)
	if hashAntes != hashDepois {
		t.Fatalf("arvore de origem foi alterada: antes=%s depois=%s", hashAntes, hashDepois)
	}
	validarRelatorioDescoberta(t, cen, raizes1, recusadas1)
	validarRelatorioDescoberta(t, cen, raizes2, recusadas2)
	validarMatrizesIguais(t, cofre1, cofre2, sum1, sum2)
}
