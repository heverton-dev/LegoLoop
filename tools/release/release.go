package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Alvo representa um par sistema operacional e arquitetura.
type Alvo struct {
	SO   string
	Arch string
}

// AlvosPadrao lista os seis alvos oficiais do release.
var AlvosPadrao = []Alvo{
	{SO: "windows", Arch: "amd64"},
	{SO: "windows", Arch: "arm64"},
	{SO: "linux", Arch: "amd64"},
	{SO: "linux", Arch: "arm64"},
	{SO: "darwin", Arch: "amd64"},
	{SO: "darwin", Arch: "arm64"},
}

// Config configura uma execucao de montagem de release.
type Config struct {
	Raiz       string
	OutDir     string
	Versao     string
	SO         string
	Arch       string
	Todos      bool
	Cosign     bool
	AssinarMSI bool
}

func main() {
	cfg, err := lerFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(2)
	}
	if err := executarCLI(cfg, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
}

func lerFlags(args []string) (Config, error) {
	fs := flag.NewFlagSet("release", flag.ContinueOnError)
	cfg := Config{}
	fs.StringVar(&cfg.SO, "so", "", "sistema operacional alvo (windows, linux, darwin)")
	fs.StringVar(&cfg.Arch, "arch", "", "arquitetura alvo (amd64, arm64)")
	fs.BoolVar(&cfg.Todos, "all", false, "compilar e empacotar todos os seis alvos")
	fs.StringVar(&cfg.OutDir, "out", "dist", "diretorio de saida dos artefatos")
	fs.StringVar(&cfg.Versao, "versao", "dev", "versao do release")
	fs.StringVar(&cfg.Raiz, "raiz", ".", "diretorio raiz do repositorio")
	fs.BoolVar(&cfg.Cosign, "cosign", false, "declarar assinatura cosign ativa no CI")
	fs.BoolVar(&cfg.AssinarMSI, "assinar-msi", false, "declarar certificado Authenticode para MSI")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	return cfg, validarConfig(cfg)
}

func validarConfig(cfg Config) error {
	if !cfg.Todos && (strings.TrimSpace(cfg.SO) == "" || strings.TrimSpace(cfg.Arch) == "") {
		return errors.New("informe --all ou o par --so e --arch")
	}
	if strings.TrimSpace(cfg.OutDir) == "" {
		return errors.New("diretorio --out vazio")
	}
	return nil
}

func executarCLI(cfg Config, saida io.Writer) error {
	if cfg.Todos {
		rel, err := Montar(cfg)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(saida, "release %s montado em %s (%d artefatos)\n", rel.Versao, cfg.OutDir, len(rel.Artefatos))
		return err
	}
	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return err
	}
	dest := filepath.Join(cfg.OutDir, NomeBinario(cfg.SO))
	hash, err := CompilarBinario(cfg.Raiz, cfg.SO, cfg.Arch, cfg.Versao, dest)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(saida, "%s  %s\n", hash, dest)
	return err
}

// NomeBinario retorna o nome do executavel para o SO informado.
func NomeBinario(so string) string {
	if so == "windows" {
		return "legoloop.exe"
	}
	return "legoloop"
}

// CompilarBinario compila ./cmd/legoloop de forma deterministica e retorna seu SHA-256.
func CompilarBinario(raiz, so, arch, _ string, destino string) (string, error) {
	if err := os.MkdirAll(filepath.Dir(destino), 0755); err != nil {
		return "", err
	}
	cmd := exec.Command(
		"go", "build",
		"-trimpath",
		"-buildvcs=false",
		"-ldflags=-s -w -buildid=",
		"-o", destino,
		"./cmd/legoloop",
	)
	cmd.Dir = raiz
	cmd.Env = ambienteBuild(os.Environ(), so, arch)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build %s/%s: %w: %s", so, arch, err, strings.TrimSpace(string(out)))
	}
	return HashArquivo(destino)
}

func ambienteBuild(base []string, so, arch string) []string {
	filtradas := make([]string, 0, len(base)+3)
	for _, kv := range base {
		if strings.HasPrefix(kv, "CGO_ENABLED=") || strings.HasPrefix(kv, "GOOS=") || strings.HasPrefix(kv, "GOARCH=") {
			continue
		}
		filtradas = append(filtradas, kv)
	}
	return append(filtradas, "CGO_ENABLED=0", "GOOS="+so, "GOARCH="+arch)
}

// HashArquivo calcula o SHA-256 hexadecimal de um arquivo em disco.
func HashArquivo(caminho string) (string, error) {
	f, err := os.Open(caminho)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Montar executa o Processo completo de release para todos os seis alvos.
func Montar(cfg Config) (RelatorioRelease, error) {
	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return RelatorioRelease{}, err
	}
	tempBins, err := os.MkdirTemp("", "legoloop-release-*")
	if err != nil {
		return RelatorioRelease{}, err
	}
	defer os.RemoveAll(tempBins)

	pacotes, arquivos, err := compilarEEmpacotarAlvos(cfg, tempBins)
	if err != nil {
		return RelatorioRelease{}, err
	}
	return finalizarArtefatos(cfg, pacotes, arquivos)
}

func compilarEEmpacotarAlvos(cfg Config, tempBins string) ([]ArtefatoAssinatura, []string, error) {
	var pacotes []ArtefatoAssinatura
	var arquivos []string
	for _, alvo := range AlvosPadrao {
		binPath := filepath.Join(tempBins, alvo.SO+"_"+alvo.Arch, NomeBinario(alvo.SO))
		if _, err := CompilarBinario(cfg.Raiz, alvo.SO, alvo.Arch, cfg.Versao, binPath); err != nil {
			return nil, nil, err
		}
		pkgPath, err := Empacotar(alvo.SO, alvo.Arch, cfg.Versao, binPath, cfg.OutDir)
		if err != nil {
			return nil, nil, err
		}
		hash, err := HashArquivo(pkgPath)
		if err != nil {
			return nil, nil, err
		}
		nome := filepath.Base(pkgPath)
		arquivos = append(arquivos, nome)
		pacotes = append(pacotes, ArtefatoAssinatura{Arquivo: nome, Tipo: "pacote", SHA256: hash})
	}
	return pacotes, arquivos, nil
}

func finalizarArtefatos(cfg Config, pacotes []ArtefatoAssinatura, arquivos []string) (RelatorioRelease, error) {
	modulo, goVer, err := LerModuloEGo(cfg.Raiz)
	if err != nil {
		return RelatorioRelease{}, err
	}
	sbomPath, err := GerarSBOM(cfg.Raiz, cfg.Versao, cfg.OutDir)
	if err != nil {
		return RelatorioRelease{}, err
	}
	sumsPath, err := GerarSHA256SUMS(cfg.OutDir, arquivos)
	if err != nil {
		return RelatorioRelease{}, err
	}
	rel := GerarRelatorioAssinatura(cfg, pacotes, filepath.Base(sumsPath), filepath.Base(sbomPath), modulo, goVer)
	if _, err := GravarRelatorio(cfg.OutDir, rel); err != nil {
		return RelatorioRelease{}, err
	}
	return rel, nil
}
