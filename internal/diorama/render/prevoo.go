package render

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/diorama/ativos"
	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
	"github.com/heverton-dev/LegoLoop/internal/jail"
)

type ItemPreVoo struct {
	Nome    string `json:"nome"`
	Nivel   string `json:"nivel"`
	Detalhe string `json:"detalhe,omitempty"`
}

type RelatorioPreVoo struct {
	Bloqueado bool         `json:"bloqueado"`
	Itens     []ItemPreVoo `json:"itens"`
}

// DepsPreVoo isola as sondas externas para testes herméticos.
type DepsPreVoo struct {
	ProvedorRegistrado func(motor string) bool
	Ferramenta         func(nome string) bool
}

// ContratosRender reune as entradas imutaveis dos Ciclos 1 e 2 e a configuracao do Ciclo 3.
type ContratosRender struct {
	Projeto    roteiro.Projeto
	Marca      roteiro.Marca
	Script     roteiro.Script
	Storyboard ativos.Storyboard
	Sprites    map[string]ativos.DescritorSprite
	Config     ConfigRender
}

func pastasDeSaida() []string { return []string{"audio", "data", "renders"} }

func resolverSeguro(projDir, relativo string) (string, error) {
	return jail.DentroDaRaiz(projDir, relativo)
}

func lerBytes(projDir, nome string) ([]byte, error) {
	caminho, err := resolverSeguro(projDir, nome)
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	return bytes.TrimPrefix(b, []byte{0xEF, 0xBB, 0xBF}), nil
}

func lerJSON(projDir, nome string, dst any, estrito bool) error {
	b, err := lerBytes(projDir, nome)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	if estrito {
		dec.DisallowUnknownFields()
	}
	return dec.Decode(dst)
}

func bloqueio(nome string, err error) error {
	return fmt.Errorf("%w: %s: %v", ErrPreVooBloqueado, nome, err)
}

// CarregarContratos le e valida todas as entradas sem acionar adaptadores externos.
func CarregarContratos(projDir string) (ContratosRender, error) {
	var c ContratosRender
	if err := carregarCiclo1(projDir, &c); err != nil {
		return c, err
	}
	if err := lerJSON(projDir, "cenas.json", &c.Storyboard, true); err != nil {
		return c, bloqueio("cenas.json", err)
	}
	if err := carregarSprites(projDir, &c); err != nil {
		return c, err
	}
	if err := carregarConfig(projDir, &c); err != nil {
		return c, err
	}
	if err := validarCobertura(c.Script, c.Storyboard); err != nil {
		return c, err
	}
	return c, validarSpritesReferenciados(c.Storyboard, c.Sprites)
}

func carregarCiclo1(projDir string, c *ContratosRender) error {
	alvos := []struct {
		nome string
		dst  any
	}{{"projeto.json", &c.Projeto}, {"marca.json", &c.Marca}, {"script.json", &c.Script}}
	for _, a := range alvos {
		if err := lerJSON(projDir, a.nome, a.dst, false); err != nil {
			return bloqueio(a.nome, err)
		}
	}
	return nil
}

func carregarSprites(projDir string, c *ContratosRender) error {
	c.Sprites = map[string]ativos.DescritorSprite{}
	dir, err := resolverSeguro(projDir, "sprites")
	if err != nil {
		return bloqueio("sprites", err)
	}
	arquivos, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	for _, arq := range arquivos {
		var d ativos.DescritorSprite
		if err := lerJSON(projDir, filepath.Join("sprites", filepath.Base(arq)), &d, false); err != nil || strings.TrimSpace(d.ID) == "" {
			return bloqueio("sprites/"+filepath.Base(arq), errors.New("descritor ilegivel ou sem id"))
		}
		c.Sprites[d.ID] = d
	}
	return nil
}

func carregarConfig(projDir string, c *ContratosRender) error {
	caminho, err := resolverSeguro(projDir, "render.json")
	if err != nil {
		return bloqueio("render.json", err)
	}
	if _, err := os.Stat(caminho); os.IsNotExist(err) {
		c.Config = ConfigPadrao()
		return nil
	}
	var cfg ConfigRender
	if err := lerJSON(projDir, "render.json", &cfg, true); err != nil {
		return bloqueio("render.json", err)
	}
	if err := cfg.Validar(); err != nil {
		return bloqueio("render.json", err)
	}
	c.Config = cfg.Normalizar()
	return nil
}

func validarCobertura(s roteiro.Script, sb ativos.Storyboard) error {
	if len(s.Lines) == 0 || len(sb.Cenas) != len(s.Lines) {
		return bloqueio("cobertura", fmt.Errorf("cenas.json tem %d cenas para %d linhas de script.json", len(sb.Cenas), len(s.Lines)))
	}
	vistos := map[int]bool{}
	for _, cena := range sb.Cenas {
		if cena.IndiceLinha < 0 || cena.IndiceLinha >= len(s.Lines) || vistos[cena.IndiceLinha] {
			return bloqueio("indice_linha", fmt.Errorf("indice %d fora do intervalo ou duplicado", cena.IndiceLinha))
		}
		vistos[cena.IndiceLinha] = true
	}
	return nil
}

func validarSpritesReferenciados(sb ativos.Storyboard, sprites map[string]ativos.DescritorSprite) error {
	for _, cena := range sb.Cenas {
		for _, el := range cena.ElementosMoveis {
			if _, ok := sprites[el.SpriteID]; el.SpriteID != "" && !ok {
				return bloqueio("sprite_id", fmt.Errorf("%q (linha %d) nao existe em sprites/", el.SpriteID, cena.IndiceLinha))
			}
		}
	}
	return nil
}

func normalizarDepsPreVoo(d DepsPreVoo) DepsPreVoo {
	if d.ProvedorRegistrado == nil {
		d.ProvedorRegistrado = RegistroPadrao(nil).Existe
	}
	if d.Ferramenta == nil {
		d.Ferramenta = ferramentaLocal
	}
	return d
}

func ferramentaLocal(nome string) bool {
	if nome == "faster-whisper" {
		return exec.Command("python", "-c", "import faster_whisper").Run() == nil
	}
	_, err := exec.LookPath(nome)
	return err == nil
}

// ExecutarPreVoo valida contratos primeiro e so entao sonda provedor, ferramentas e escrita.
func ExecutarPreVoo(ctx context.Context, projDir string, deps DepsPreVoo) (RelatorioPreVoo, ContratosRender, error) {
	var rel RelatorioPreVoo
	deps = normalizarDepsPreVoo(deps)
	c, err := CarregarContratos(projDir)
	if err != nil {
		err = registrarFalta(&rel, "contratos", err)
		return rel, c, err
	}
	adicionar(&rel, ItemPreVoo{Nome: "contratos", Nivel: "OK"})
	passos := []func() error{
		func() error { return checarProvedor(&rel, c.Marca, deps) },
		func() error { return checarFerramentas(&rel, deps) },
		func() error { return checarEscrita(&rel, projDir) },
	}
	for _, passo := range passos {
		if err := ctx.Err(); err != nil {
			return rel, c, bloqueio("contexto", err)
		}
		if err := passo(); err != nil {
			return rel, c, err
		}
	}
	return rel, c, nil
}

func checarProvedor(rel *RelatorioPreVoo, m roteiro.Marca, deps DepsPreVoo) error {
	if !deps.ProvedorRegistrado(m.Voz.Motor) {
		err := fmt.Errorf("%w: %q", ErrProvedorDesconhecido, m.Voz.Motor)
		return registrarFalta(rel, "provedor_voz", err)
	}
	adicionar(rel, ItemPreVoo{Nome: "provedor_voz", Nivel: "OK", Detalhe: m.Voz.Motor})
	return nil
}

func checarFerramentas(rel *RelatorioPreVoo, deps DepsPreVoo) error {
	for _, nome := range []string{"ffmpeg", "faster-whisper"} {
		if !deps.Ferramenta(nome) {
			return registrarFalta(rel, nome, bloqueio(nome, errors.New("ferramenta ausente")))
		}
		adicionar(rel, ItemPreVoo{Nome: nome, Nivel: "OK"})
	}
	return nil
}

func checarEscrita(rel *RelatorioPreVoo, projDir string) error {
	for _, sub := range pastasDeSaida() {
		if err := sondarEscrita(projDir, sub); err != nil {
			return registrarFalta(rel, "pasta_"+sub, bloqueio(sub, err))
		}
		adicionar(rel, ItemPreVoo{Nome: "pasta_" + sub, Nivel: "OK"})
	}
	return nil
}

func sondarEscrita(projDir, sub string) error {
	dir, err := resolverSeguro(projDir, sub)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	probe := filepath.Join(dir, ".prevoo_probe")
	if err := os.WriteFile(probe, []byte("ok"), 0o644); err != nil {
		return err
	}
	return os.Remove(probe)
}

func adicionar(rel *RelatorioPreVoo, it ItemPreVoo) {
	rel.Itens = append(rel.Itens, it)
	if it.Nivel == "FALTA" {
		rel.Bloqueado = true
	}
}

func registrarFalta(rel *RelatorioPreVoo, nome string, err error) error {
	adicionar(rel, ItemPreVoo{Nome: nome, Nivel: "FALTA", Detalhe: err.Error()})
	return err
}
