package ste

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// ErrInvalido sinaliza valor inválido para modo, pausa ou harness.
var ErrInvalido = errors.New("ste: entrada inválida no alternador")

// Modo define o comportamento e estilo de resposta e governança.
type Modo string

const (
	ModoOff     Modo = "off"
	ModoStept   Modo = "stept"
	ModoCaveman Modo = "caveman"
	ModoHibrido Modo = "hibrido"
	ModoPuro    Modo = "puro"
)

// Pausa define a política de interrupção antes de ações com efeito.
type Pausa string

const (
	PausaNunca      Pausa = "nunca"
	PausaDestrutivo Pausa = "destrutivo"
	PausaSempre     Pausa = "sempre"
)

// HarnessID identifica os harnesses suportados.
type HarnessID string

const (
	HarnessGemini   HarnessID = "gemini"
	HarnessClaude   HarnessID = "claude"
	HarnessOMP      HarnessID = "omp"
	HarnessOpenCode HarnessID = "opencode"
	HarnessMimo     HarnessID = "mimo"
	HarnessCursor   HarnessID = "cursor"
	HarnessWindsurf HarnessID = "windsurf"
	HarnessVSCode   HarnessID = "vscode"
)

// ConfigPerfil define a configuração de governança de um harness ou padrão.
type ConfigPerfil struct {
	Modo     Modo   `json:"modo"`
	Pausa    Pausa  `json:"pausa"`
	Thinking string `json:"thinking,omitempty"`
}

// Config representa a estrutura persistida de governança do STE-PT.
type Config struct {
	Padrao    ConfigPerfil               `json:"padrao"`
	Harnesses map[HarnessID]ConfigPerfil `json:"harnesses,omitempty"`
}

// HarnessStatus reporta o estado em tempo real de um harness.
type HarnessStatus struct {
	ID        HarnessID    `json:"id"`
	Nome      string       `json:"nome"`
	Instalado bool         `json:"instalado"`
	Caminho   string       `json:"caminho"`
	Ativo     bool         `json:"ativo"`
	Perfil    ConfigPerfil `json:"perfil"`
	Origem    string       `json:"origem"`
}

// Validate valida se modo e pausa são conhecidos e aceitáveis.
func (c ConfigPerfil) Validate() error {
	switch c.Modo {
	case ModoOff, ModoStept, ModoCaveman, ModoHibrido, ModoPuro:
	default:
		return ErrInvalido
	}
	switch c.Pausa {
	case PausaNunca, PausaDestrutivo, PausaSempre:
	default:
		return ErrInvalido
	}
	return nil
}

// PerfilPadraoRetorno devolve a tupla canônica de governança.
func PerfilPadraoRetorno() ConfigPerfil {
	return ConfigPerfil{
		Modo:     ModoHibrido,
		Pausa:    PausaDestrutivo,
		Thinking: "caveman",
	}
}

// NovoConfig inicializa uma estrutura vazia com perfil padrão canônico.
func NovoConfig() *Config {
	return &Config{
		Padrao:    PerfilPadraoRetorno(),
		Harnesses: make(map[HarnessID]ConfigPerfil),
	}
}

// ValidarHarness verifica se o identificador de harness é aceito.
func ValidarHarness(h HarnessID) bool {
	switch h {
	case HarnessGemini, HarnessClaude, HarnessOMP, HarnessOpenCode, HarnessMimo, HarnessCursor, HarnessWindsurf, HarnessVSCode:
		return true
	default:
		return false
	}
}

// ResolverPerfil resolve herança entre perfil customizado do harness e padrão.
func (c *Config) ResolverPerfil(h HarnessID) (ConfigPerfil, string) {
	if c == nil {
		return PerfilPadraoRetorno(), "padrao"
	}
	if p, ok := c.Harnesses[h]; ok {
		return p, "override"
	}
	return c.Padrao, "padrao"
}

// DefinirPadrao atualiza o perfil global padrão.
func (c *Config) DefinirPadrao(p ConfigPerfil) error {
	if err := p.Validate(); err != nil {
		return err
	}
	c.Padrao = p
	return nil
}

// DefinirHarness atualiza ou cria o override de um harness específico.
func (c *Config) DefinirHarness(h HarnessID, p ConfigPerfil) error {
	if !ValidarHarness(h) {
		return ErrInvalido
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if c.Harnesses == nil {
		c.Harnesses = make(map[HarnessID]ConfigPerfil)
	}
	c.Harnesses[h] = p
	return nil
}

// RemoverHarness remove o override de um harness, voltando ao padrão.
func (c *Config) RemoverHarness(h HarnessID) {
	if c.Harnesses != nil {
		delete(c.Harnesses, h)
	}
}

// CaminhoConfigPadrao retorna o caminho de persistência ~/.stept/config.json.
func CaminhoConfigPadrao() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".stept", "config.json")
}

// CarregarConfig lê a configuração do disco ou cria padrão se ausente.
func CarregarConfig(caminho string) (*Config, error) {
	dados, err := os.ReadFile(caminho)
	if err != nil {
		if os.IsNotExist(err) {
			return NovoConfig(), nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(dados, &cfg); err != nil {
		return nil, err
	}
	if cfg.Harnesses == nil {
		cfg.Harnesses = make(map[HarnessID]ConfigPerfil)
	}
	return &cfg, nil
}

// SalvarConfig persiste a configuração no disco atomicamente com 0600.
func SalvarConfig(caminho string, cfg *Config) error {
	if cfg == nil {
		return ErrInvalido
	}
	if err := cfg.Padrao.Validate(); err != nil {
		return err
	}
	for _, p := range cfg.Harnesses {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	dados, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(caminho), 0755); err != nil {
		return err
	}
	return os.WriteFile(caminho, dados, 0600)
}
