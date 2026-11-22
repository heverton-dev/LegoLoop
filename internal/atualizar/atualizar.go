package atualizar

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
)

const (
	// CanalEstavel identifica o canal estavel de atualizacao.
	CanalEstavel = "estavel"
	// CanalBeta identifica o canal beta de atualizacao.
	CanalBeta = "beta"
	// BaseURLPadrao aponta para a origem oficial de releases.
	BaseURLPadrao = "https://github.com/heverton-dev/LegoLoop/releases/latest/download"
)

var (
	// ErrCanalInvalido sinaliza canal fora de estavel e beta.
	ErrCanalInvalido = errors.New("atualizar: canal invalido")
	// ErrShaDivergente sinaliza divergencia de SHA-256 no pacote baixado.
	ErrShaDivergente = errors.New("atualizar: sha256 divergente")
	// ErrShaNaoEncontrado sinaliza ausencia do artefato em SHA256SUMS.
	ErrShaNaoEncontrado = errors.New("atualizar: entrada ausente em SHA256SUMS")
	// ErrDerivaDetectada sinaliza deriva da Forja apos a atualizacao.
	ErrDerivaDetectada = errors.New("atualizar: deriva detectada apos atualizacao")
	// ErrSemBinario sinaliza caminho de binario vazio.
	ErrSemBinario = errors.New("atualizar: caminho do binario vazio")
	// ErrSemBaseURL sinaliza URL base vazia.
	ErrSemBaseURL = errors.New("atualizar: url base vazia")
	// ErrSemHome sinaliza diretorio HOME vazio.
	ErrSemHome = errors.New("atualizar: home vazio")
)

// Opcao configura uma execucao de atualizacao do binario.
type Opcao struct {
	Canal    string       `json:"canal"`
	BaseURL  string       `json:"base_url"`
	Binario  string       `json:"binario"`
	Home     string       `json:"home"`
	Artefato string       `json:"artefato,omitempty"`
	Cliente  *http.Client `json:"-"`
}

// Resultado resume a atualizacao executada e a revalidacao da Forja.
type Resultado struct {
	Canal          string `json:"canal"`
	Artefato       string `json:"artefato"`
	Binario        string `json:"binario"`
	BinarioPrev    string `json:"binario_prev"`
	SHA256Anterior string `json:"sha256_anterior,omitempty"`
	SHA256Novo     string `json:"sha256_novo"`
	Revalidado     bool   `json:"revalidado"`
}

// NomeArtefato devolve o nome canonico do pacote para o SO e arquitetura.
func NomeArtefato(goos, goarch string) string {
	nome := fmt.Sprintf("legoloop_%s_%s", goos, goarch)
	if goos == "windows" {
		return nome + ".exe"
	}
	return nome
}

// Validar confere os parametros obrigatorios antes de qualquer rede ou disco.
func (o Opcao) Validar() error {
	if o.Canal != CanalEstavel && o.Canal != CanalBeta {
		return fmt.Errorf("%w: %q", ErrCanalInvalido, o.Canal)
	}
	if strings.TrimSpace(o.BaseURL) == "" {
		return ErrSemBaseURL
	}
	if strings.TrimSpace(o.Binario) == "" {
		return ErrSemBinario
	}
	if strings.TrimSpace(o.Home) == "" {
		return ErrSemHome
	}
	return nil
}

// Executar baixa SHA256SUMS e o pacote, confere o SHA-256, guarda .prev e revalida a Forja.
func Executar(ctx context.Context, op Opcao) (Resultado, error) {
	if err := op.Validar(); err != nil {
		return Resultado{}, err
	}
	artefato := strings.TrimSpace(op.Artefato)
	if artefato == "" {
		artefato = NomeArtefato(runtime.GOOS, runtime.GOARCH)
	}
	cliente := op.Cliente
	if cliente == nil {
		cliente = http.DefaultClient
	}
	base := strings.TrimRight(op.BaseURL, "/") + "/" + op.Canal
	shaEsperado, err := obterShaEsperado(ctx, cliente, base+"/SHA256SUMS", artefato)
	if err != nil {
		return Resultado{}, err
	}
	dadosNovo, err := baixarBytes(ctx, cliente, base+"/"+artefato)
	if err != nil {
		return Resultado{}, err
	}
	shaNovo := calcularSHA256(dadosNovo)
	if !strings.EqualFold(shaNovo, shaEsperado) {
		return Resultado{}, fmt.Errorf("%w: esperado=%s obtido=%s", ErrShaDivergente, shaEsperado, shaNovo)
	}
	return substituirERevalidar(op, dadosNovo, artefato, shaNovo)
}

func obterShaEsperado(ctx context.Context, cliente *http.Client, url, artefato string) (string, error) {
	dados, err := baixarBytes(ctx, cliente, url)
	if err != nil {
		return "", fmt.Errorf("atualizar: baixar SHA256SUMS: %w", err)
	}
	return extrairSHA256(string(dados), artefato)
}

func extrairSHA256(conteudo, artefato string) (string, error) {
	sc := bufio.NewScanner(strings.NewReader(conteudo))
	for sc.Scan() {
		campos := strings.Fields(sc.Text())
		if len(campos) < 2 {
			continue
		}
		nome := strings.TrimPrefix(campos[1], "*")
		if nome == artefato {
			return strings.ToLower(campos[0]), nil
		}
	}
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("atualizar: ler SHA256SUMS: %w", err)
	}
	return "", fmt.Errorf("%w: %s", ErrShaNaoEncontrado, artefato)
}

func baixarBytes(ctx context.Context, cliente *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("atualizar: montar requisicao %s: %w", url, err)
	}
	resp, err := cliente.Do(req)
	if err != nil {
		return nil, fmt.Errorf("atualizar: requisitar %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fecharCorpo(resp, fmt.Errorf("atualizar: status HTTP %d em %s", resp.StatusCode, url))
	}
	dados, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fecharCorpo(resp, fmt.Errorf("atualizar: ler resposta %s: %w", url, err))
	}
	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("atualizar: fechar resposta %s: %w", url, err)
	}
	return dados, nil
}

func fecharCorpo(resp *http.Response, causa error) error {
	if err := resp.Body.Close(); err != nil {
		return fmt.Errorf("%w (fechar corpo: %v)", causa, err)
	}
	return causa
}

func calcularSHA256(dados []byte) string {
	soma := sha256.Sum256(dados)
	return hex.EncodeToString(soma[:])
}
