package homologacao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const HostCanonico = "roteador.signa-docs.com"

const SufixoBackup = ".lego" + "lo" + "op-bak"

type Endereco struct {
	Arquivo string
	Linha   int
	Host    string
}

type Sonda struct {
	Arquivo          string
	Host             string
	Resolve          bool
	Versao           string
	VersaoDisponivel string
	Antigo           bool
	Combos           []string
	Erro             string
}

type Mudanca struct {
	Arquivo   string
	Enderecos []Endereco
	Backup    string
	Aplicada  bool
	Conferida bool
}

type achado struct {
	linha, ini, fim int
	host            string
}

var (
	reCabecalho = regexp.MustCompile(`^\s*"?providers?"?\s*:`)
	reEndereco  = regexp.MustCompile(`(?i)"?baseurl"?\s*:\s*['"]?([^'",\s]+)`)
)

func achar(linha string, n, off int) []achado {
	t := strings.TrimSpace(linha)
	if strings.HasPrefix(t, "//") || strings.HasPrefix(t, "#") {
		return nil
	}
	var out []achado
	for _, m := range reEndereco.FindAllStringSubmatchIndex(linha, -1) {
		ini, fim := m[2], m[3]
		if i := strings.Index(linha[ini:fim], "://"); i >= 0 {
			ini += i + 3
		}
		if i := strings.Index(linha[ini:fim], "/"); i >= 0 {
			fim = ini + i
		}
		out = append(out, achado{n, off + ini, off + fim, linha[ini:fim]})
	}
	return out
}

func fimDeSecao(linha string) bool {
	if linha == "" {
		return false
	}
	return !strings.ContainsRune(" \t#\r\n", rune(linha[0]))
}

func varrer(texto string) []achado {
	var out []achado
	off, dentro, raiz := 0, false, false
	for i, linha := range strings.SplitAfter(texto, "\n") {
		switch {
		case reCabecalho.MatchString(linha):
			dentro, raiz = true, strings.HasPrefix(linha, "provider")
		case dentro && raiz && fimDeSecao(linha):
			dentro = false
		}
		if dentro {
			out = append(out, achar(linha, i+1, off)...)
		}
		off += len(linha)
	}
	return out
}

func LerEnderecos(arquivo string) ([]Endereco, error) {
	dados, err := os.ReadFile(arquivo)
	if err != nil {
		return nil, err
	}
	out := []Endereco{}
	for _, a := range varrer(string(dados)) {
		out = append(out, Endereco{Arquivo: arquivo, Linha: a.linha, Host: a.host})
	}
	return out, nil
}

func obter(ctx context.Context, cli *http.Client, url, chave string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	if chave != "" {
		req.Header.Set("Authorization", "Bearer "+chave)
	}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	corpo, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	return corpo, resp.StatusCode, err
}

func limpar(msg, chave string) string {
	if chave == "" {
		return msg
	}
	return strings.ReplaceAll(msg, chave, "***")
}

func lerVersao(corpo []byte) string {
	var v struct {
		Version string `json:"version"`
	}
	_ = json.Unmarshal(corpo, &v)
	return v.Version
}

func lerCombos(corpo []byte) []string {
	var v struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(corpo, &v)
	out := []string{}
	for _, m := range v.Data {
		out = append(out, m.ID)
	}
	return out
}

func SondarHost(ctx context.Context, cli *http.Client, esquema, host, chave string) Sonda {
	s := Sonda{Host: host, Combos: []string{}}
	if cli == nil {
		cli = http.DefaultClient
	}
	base := esquema + "://" + host
	corpo, status, err := obter(ctx, cli, base+"/api/version", chave)
	if err != nil {
		s.Erro = limpar(err.Error(), chave)
		return s
	}
	s.Resolve = true
	if status == http.StatusOK {
		s.Versao = lerVersao(corpo)
	} else {
		s.Erro = fmt.Sprintf("status %d em /api/version", status)
	}
	if corpo, status, err = obter(ctx, cli, base+"/v1/models", chave); err == nil && status == http.StatusOK {
		s.Combos = lerCombos(corpo)
	}
	return s
}

func ComparaVersao(a, b string) int {
	pa, pb := strings.Split(strings.TrimPrefix(a, "v"), "."), strings.Split(strings.TrimPrefix(b, "v"), ".")
	for i := 0; i < max(len(pa), len(pb)); i++ {
		x, y := parte(pa, i), parte(pb, i)
		if x != y {
			if x < y {
				return -1
			}
			return 1
		}
	}
	return 0
}

func parte(partes []string, i int) int {
	if i >= len(partes) {
		return 0
	}
	n, _ := strconv.Atoi(partes[i])
	return n
}

func (s *Sonda) Comparar(disponivel string) {
	s.VersaoDisponivel = disponivel
	s.Antigo = s.Versao != "" && disponivel != "" && ComparaVersao(s.Versao, disponivel) < 0
}

func filtrar(todos []achado, host string) []achado {
	var out []achado
	for _, a := range todos {
		if a.host == host {
			out = append(out, a)
		}
	}
	return out
}

func trocar(texto string, alvos []achado, novo string) string {
	for i := len(alvos) - 1; i >= 0; i-- {
		texto = texto[:alvos[i].ini] + novo + texto[alvos[i].fim:]
	}
	return texto
}

func conferir(arquivo, antigo, novo string, esperados int) error {
	ends, err := LerEnderecos(arquivo)
	if err != nil {
		return err
	}
	achados := 0
	for _, e := range ends {
		if e.Host == antigo {
			return fmt.Errorf("conferência falhou em %s: ainda há %s", arquivo, antigo)
		}
		if e.Host == novo {
			achados++
		}
	}
	if achados < esperados {
		return fmt.Errorf("conferência falhou em %s: %s não encontrado", arquivo, novo)
	}
	return nil
}

func Corrigir(arquivo, antigo, novo string, confirmar bool) (Mudanca, error) {
	m := Mudanca{Arquivo: arquivo}
	dados, err := os.ReadFile(arquivo)
	if err != nil || antigo == novo {
		return m, err
	}
	texto := string(dados)
	alvos := filtrar(varrer(texto), antigo)
	for _, a := range alvos {
		m.Enderecos = append(m.Enderecos, Endereco{Arquivo: arquivo, Linha: a.linha, Host: a.host})
	}
	if len(alvos) == 0 || !confirmar {
		return m, nil
	}
	info, err := os.Stat(arquivo)
	if err != nil {
		return m, err
	}
	m.Backup = arquivo + SufixoBackup
	// Um backup que já existe guarda o original de uma correção anterior (outro host do mesmo arquivo): não o sobrescreva.
	if _, err := os.Stat(m.Backup); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(m.Backup, dados, info.Mode().Perm()); err != nil {
			return m, err
		}
	}
	if err := os.WriteFile(arquivo, []byte(trocar(texto, alvos, novo)), info.Mode().Perm()); err != nil {
		return m, err
	}
	m.Aplicada = true
	if err := conferir(arquivo, antigo, novo, len(alvos)); err != nil {
		return m, err
	}
	m.Conferida = true
	return m, nil
}
