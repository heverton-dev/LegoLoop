// Package renomear renomeia um Processo de forma determinística em todos os
// pontos do repositório e do Cofre, com prévia (--dry-run), backup com
// MANIFESTO.tsv e registro em espinha/renomeados.json.
package renomear

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/caminhos"
	"github.com/heverton-dev/LegoLoop/internal/jail"
	"github.com/heverton-dev/LegoLoop/internal/processo"
)

var (
	// ErrNomeInvalido sinaliza nome antigo ou novo inválido ou iguais entre si.
	ErrNomeInvalido = errors.New("renomear: nome inválido")
	// ErrNaoEncontrado sinaliza que o Processo antigo não existe no projeto nem no Cofre.
	ErrNaoEncontrado = errors.New("renomear: processo não encontrado")
	// ErrColisao sinaliza que já existe um Processo com o novo nome.
	ErrColisao = errors.New("renomear: destino já existe")
)

// Opcoes configura uma operação de renomeação de Processo.
type Opcoes struct {
	Antigo    string
	Novo      string
	Raiz      string
	Vault     string
	Home      string
	Docs      string
	BackupDir string
	DryRun    bool
	Agora     func() time.Time
}

// Alteracao descreve uma mudança prevista ou executada.
type Alteracao struct {
	Tipo    string `json:"tipo"`
	Origem  string `json:"origem"`
	Destino string `json:"destino,omitempty"`
	Detalhe string `json:"detalhe,omitempty"`
}

// Relatorio resume o plano ou o resultado da renomeação.
type Relatorio struct {
	Antigo     string      `json:"antigo"`
	Novo       string      `json:"novo"`
	DryRun     bool        `json:"dry_run"`
	BackupDir  string      `json:"backup_dir,omitempty"`
	Manifesto  string      `json:"manifesto,omitempty"`
	Registro   string      `json:"registro"`
	Alteracoes []Alteracao `json:"alteracoes"`
}

// RegistroRenomeado representa uma entrada em espinha/renomeados.json.
type RegistroRenomeado struct {
	Antigo string `json:"antigo"`
	Novo   string `json:"novo"`
	Data   string `json:"data"`
}

// Executar valida, planeja e (quando DryRun é falso) aplica a renomeação.
func Executar(op Opcoes) (Relatorio, error) {
	op = normalizarOpcoes(op)
	if err := validarOpcoes(op); err != nil {
		return Relatorio{}, err
	}
	rel, err := planejar(op)
	if err != nil || op.DryRun {
		return rel, err
	}
	if err := aplicar(op, &rel); err != nil {
		return Relatorio{}, err
	}
	return rel, nil
}

func normalizarOpcoes(op Opcoes) Opcoes {
	if strings.TrimSpace(op.Raiz) == "" {
		op.Raiz = "."
	}
	if strings.TrimSpace(op.Docs) == "" {
		op.Docs = filepath.Join(op.Raiz, "docs")
	}
	if op.Agora == nil {
		op.Agora = time.Now
	}
	return op
}

func validarOpcoes(op Opcoes) error {
	if err := jail.ValidarNomeProcesso(op.Antigo); err != nil {
		return fmt.Errorf("%w: %w", ErrNomeInvalido, err)
	}
	if err := jail.ValidarNomeProcesso(op.Novo); err != nil {
		return fmt.Errorf("%w: %w", ErrNomeInvalido, err)
	}
	if op.Antigo == op.Novo {
		return fmt.Errorf("%w: antigo e novo são iguais (%q)", ErrNomeInvalido, op.Antigo)
	}
	return conferirExistenciaEColisao(op)
}

func conferirExistenciaEColisao(op Opcoes) error {
	origens := caminhosExistentesProcesso(op, op.Antigo)
	if len(origens) == 0 {
		return fmt.Errorf("%w: %q", ErrNaoEncontrado, op.Antigo)
	}
	if col := caminhosExistentesProcesso(op, op.Novo); len(col) > 0 {
		return fmt.Errorf("%w: %q já existe em %s", ErrColisao, op.Novo, col[0])
	}
	return nil
}

func caminhosExistentesProcesso(op Opcoes, slug string) []string {
	var encontrados []string
	candidatos := []string{filepath.Join(op.Raiz, caminhos.DirProcessos, slug)}
	if op.Vault != "" {
		candidatos = append(candidatos,
			filepath.Join(op.Vault, caminhos.DirProcessos, slug),
			filepath.Join(op.Vault, "skills", slug),
			processo.CaminhoNoCofre(op.Vault, slug),
		)
	}
	for _, c := range candidatos {
		if _, err := os.Lstat(c); err == nil {
			encontrados = append(encontrados, c)
		}
	}
	return encontrados
}

func planejar(op Opcoes) (Relatorio, error) {
	regPath := caminhoRegistro(op.Raiz)
	rel := Relatorio{
		Antigo:     op.Antigo,
		Novo:       op.Novo,
		DryRun:     op.DryRun,
		Registro:   regPath,
		Alteracoes: []Alteracao{},
	}
	rel.Alteracoes = append(rel.Alteracoes, planejarPastaProjeto(op)...)
	rel.Alteracoes = append(rel.Alteracoes, planejarCofreEVinculos(op)...)
	docsAlt, err := planejarDocs(op)
	if err != nil {
		return Relatorio{}, err
	}
	rel.Alteracoes = append(rel.Alteracoes, docsAlt...)
	rel.Alteracoes = append(rel.Alteracoes, Alteracao{
		Tipo:    "registro",
		Origem:  regPath,
		Destino: regPath,
		Detalhe: fmt.Sprintf("%s -> %s", op.Antigo, op.Novo),
	})
	return rel, nil
}

func aplicar(op Opcoes, rel *Relatorio) error {
	bdir := resolverBackupDir(op)
	bk, err := novoBackup(bdir)
	if err != nil {
		return err
	}
	if err := aplicarPastaProjeto(op, bk); err != nil {
		return err
	}
	if err := aplicarCofreEVinculos(op, bk); err != nil {
		return err
	}
	if err := aplicarDocs(op, rel.Alteracoes, bk); err != nil {
		return err
	}
	if err := registrarRenomeado(op, bk); err != nil {
		return err
	}
	manifesto, err := bk.gravarManifesto()
	if err != nil {
		return err
	}
	rel.BackupDir = bdir
	rel.Manifesto = manifesto
	return nil
}

func caminhoRegistro(raiz string) string {
	return filepath.Join(raiz, "espinha", "renomeados.json")
}

// BuscarRenomeado consulta espinha/renomeados.json sob raiz procurando slug como antigo.
func BuscarRenomeado(raiz, slug string) (RegistroRenomeado, bool) {
	regs, err := lerRenomeados(caminhoRegistro(raiz))
	if err != nil {
		return RegistroRenomeado{}, false
	}
	for _, r := range regs {
		if r.Antigo == slug {
			return r, true
		}
	}
	return RegistroRenomeado{}, false
}

func lerRenomeados(caminho string) ([]RegistroRenomeado, error) {
	dados, err := os.ReadFile(caminho)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var regs []RegistroRenomeado
	if err := json.Unmarshal(dados, &regs); err != nil {
		return nil, err
	}
	return regs, nil
}

func registrarRenomeado(op Opcoes, bk *backupCtx) error {
	caminho := caminhoRegistro(op.Raiz)
	if _, err := os.Stat(caminho); err == nil {
		if err := bk.salvar(caminho, "atualizar_registro"); err != nil {
			return err
		}
	}
	regs, err := lerRenomeados(caminho)
	if err != nil {
		return err
	}
	novoReg := RegistroRenomeado{
		Antigo: op.Antigo,
		Novo:   op.Novo,
		Data:   op.Agora().Format("2006-01-02"),
	}
	regs = upsertRegistro(regs, novoReg)
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		return err
	}
	dados, err := json.MarshalIndent(regs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(caminho, append(dados, '\n'), 0o644)
}

func upsertRegistro(regs []RegistroRenomeado, item RegistroRenomeado) []RegistroRenomeado {
	for i := range regs {
		if regs[i].Antigo == item.Antigo {
			regs[i] = item
			return regs
		}
	}
	return append(regs, item)
}
