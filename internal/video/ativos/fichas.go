package ativos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
)

// ErrRenderizadorNaoRecorte indica estilo cujo renderizador nao usa fichas e rembg (AC-E09.3).
var ErrRenderizadorNaoRecorte = errors.New("diorama-ativos: renderizador do estilo nao e de recorte")

// RenderizadorDeRecorte indica renderizador "recorte" ou "recorte-<variante>", que segue fichas e rembg.
func RenderizadorDeRecorte(renderizador string) bool {
	r := strings.ToLower(strings.TrimSpace(renderizador))
	return r == "recorte" || strings.HasPrefix(r, "recorte-")
}

// EntradaEtapa3Recorte reune o estilo, as fichas 4 faixas e o elenco da Etapa 3 de recorte.
type EntradaEtapa3Recorte struct {
	Projeto    roteiro.Projeto
	Estilo     inventario.Estilo
	Fichas     []FichaAtivo
	Elenco     []FichaPersonagem
	Deps       DepsCenarios
	Recortador RecortadorSprite
}

// ExecutarEtapa3Recorte mantem o caminho legado: gera as fichas e depois recorta com rembg (AC-E09.3).
func ExecutarEtapa3Recorte(ctx context.Context, projDir string, in EntradaEtapa3Recorte) ([]DescritorSprite, error) {
	if !RenderizadorDeRecorte(in.Estilo.Renderizador) {
		return nil, fmt.Errorf("%w: estilo %q renderizador %q", ErrRenderizadorNaoRecorte, in.Estilo.ID, in.Estilo.Renderizador)
	}
	if _, err := ExecutarGeracaoFichas4Faixas(ctx, projDir, in.Projeto, in.Fichas, in.Deps); err != nil {
		return nil, err
	}
	return ExecutarEtapa3Elenco(ctx, projDir, in.Elenco, in.Recortador)
}

// CompilarPromptFicha4Faixas monta o prompt deterministico no padrao 4 faixas (Turnaround, Expressions, Poses, Build Notes).
func CompilarPromptFicha4Faixas(f FichaAtivo) string {
	return fmt.Sprintf(
		"Character/Asset Master Sheet for %s (%s). %s. "+
			"4 horizontal bands on clean neutral light background: "+
			"1) TURNAROUND: %s. "+
			"2) EXPRESSIONS / STATES: %s. "+
			"3) POSES / ACTIONS: %s. "+
			"4) BUILD NOTES: Material=%s; Details=%s; Accessories=%s; Rules=%s. "+
			"Consistent proportions, flat crisp lighting, no watermarks.",
		f.Nome, f.Categoria, f.DescricaoImutavel,
		strings.Join(f.Turnaround, ", "),
		strings.Join(f.Expressoes, ", "),
		strings.Join(f.Poses, ", "),
		f.BuildNotes.CorpoMaterial, f.BuildNotes.OlhosDetalhes,
		f.BuildNotes.AcessoriosEfeitos, f.BuildNotes.Regras,
	)
}

// ExtrairFichasDeEntrada consolida fichas de 4 faixas a partir de script.ativos e dispositivos de plano.md.
func ExtrairFichasDeEntrada(script roteiro.Script, planoMD string) []FichaAtivo {
	saida := []FichaAtivo{fichaPersonagemPadrao4Faixas(script.Title)}
	vistos := map[string]bool{"narrador": true}
	for _, a := range script.Ativos {
		id := strings.TrimSpace(a.ID)
		if id == "" || vistos[id] {
			continue
		}
		vistos[id] = true
		saida = append(saida, converterAtivoPrevisto(a))
	}
	for i, disp := range ExtrairDispositivosPlanoMD(planoMD) {
		id := fmt.Sprintf("obj_batida_%02d", i+1)
		if vistos[id] {
			continue
		}
		vistos[id] = true
		saida = append(saida, fichaObjetoDeDispositivo(id, disp))
	}
	return saida
}

// ExecutarGeracaoFichas4Faixas gera as folhas mestres (Flow ou PNG local) com idempotencia por hash.
func ExecutarGeracaoFichas4Faixas(ctx context.Context, projDir string, proj roteiro.Projeto, fichas []FichaAtivo, deps DepsCenarios) ([]FichaAtivo, error) {
	deps = normalizarDepsCenarios(deps)
	cache := carregarCacheFichas(projDir)
	saida := make([]FichaAtivo, 0, len(fichas))
	for _, f := range fichas {
		if err := f.Validar(); err != nil {
			return nil, err
		}
		gerada, err := gerarOuReusarFicha4Faixas(ctx, projDir, proj, f, cache, deps)
		if err != nil {
			return nil, err
		}
		saida = append(saida, gerada)
	}
	if err := gravarJSONSeguro(projDir, filepath.Join("flow", "fichas.json"), saida); err != nil {
		return nil, err
	}
	return saida, nil
}

func gerarOuReusarFicha4Faixas(ctx context.Context, projDir string, proj roteiro.Projeto, f FichaAtivo, cache map[string]FichaAtivo, deps DepsCenarios) (FichaAtivo, error) {
	f.PromptFlow = CompilarPromptFicha4Faixas(f)
	h := sha256.Sum256([]byte(f.ID + "|" + f.PromptFlow + "|" + proj.EstiloVisual + "|" + proj.FlowProjeto))
	f.HashEntrada = hex.EncodeToString(h[:])
	if ant, ok := cache[f.ID]; ok && ant.HashEntrada == f.HashEntrada && arquivoExisteSeguro(projDir, ant.FolhaMestrePNG) {
		return ant, nil
	}
	relPNG := filepath.ToSlash(filepath.Join("flow", "fichas", f.ID+"_sheet.png"))
	if proj.EstiloVisual == "codigo" {
		if err := gerarSpritePNGLocal(projDir, relPNG, f.Categoria); err != nil {
			return FichaAtivo{}, err
		}
	} else {
		if deps.Flow == nil {
			return FichaAtivo{}, fmt.Errorf("%w: cliente Flow ausente para ficha %s", ErrPreVooBloqueado, f.ID)
		}
		req := RequisicaoFlow{CenarioID: "ficha_" + f.ID, Prompt: f.PromptFlow, FlowProjeto: proj.FlowProjeto}
		dados, err := executarComBackoff(ctx, req, deps)
		if err != nil {
			return FichaAtivo{}, err
		}
		if err := gravarBytesSeguro(projDir, relPNG, dados); err != nil {
			return FichaAtivo{}, err
		}
	}
	f.FolhaMestrePNG = relPNG
	if err := gravarJSONSeguro(projDir, filepath.Join("flow", "fichas", f.ID+".json"), f); err != nil {
		return FichaAtivo{}, err
	}
	return f, nil
}

func carregarCacheFichas(projDir string) map[string]FichaAtivo {
	out := make(map[string]FichaAtivo)
	var lista []FichaAtivo
	if err := lerJSONSeguro(projDir, filepath.Join("flow", "fichas.json"), &lista); err != nil {
		return out
	}
	for _, f := range lista {
		out[f.ID] = f
	}
	return out
}

func fichaPersonagemPadrao4Faixas(titulo string) FichaAtivo {
	return FichaAtivo{
		ID:                "narrador",
		Categoria:         "personagem",
		Nome:              "Narrador",
		DescricaoImutavel: "Estilizado personagem miniatura 3D minimalista para " + titulo,
		Turnaround:        []string{"Front", "Three-quarter", "Side", "Back three-quarter", "Back"},
		Expressoes:        []string{"Neutral", "Wide", "Narrowed", "Slanted in", "Squeezed"},
		Poses:             []string{"Explaining", "Pointing", "Thinking", "Reacting"},
		BuildNotes: BuildNotes{
			CorpoMaterial:     "Matte terracotta and charcoal stylized miniature body",
			OlhosDetalhes:     "Two crisp white glowing slot eyes with expressive lids",
			AcessoriosEfeitos: "Minimalist lapel microphone and subtle rim light",
			Regras:            "Keep exact proportions, white eyes and palette in every shot",
		},
	}
}

func converterAtivoPrevisto(a roteiro.AtivoPrevisto) FichaAtivo {
	cat := a.Categoria
	if cat != "personagem" && cat != "objeto" && cat != "motion" {
		cat = "objeto"
	}
	regras := strings.Join(a.BuildNotes, "; ")
	if strings.TrimSpace(regras) == "" {
		regras = "Keep exact silhouette and palette in every scene"
	}
	return FichaAtivo{
		ID:                a.ID,
		Categoria:         cat,
		Nome:              a.Nome,
		DescricaoImutavel: a.DescricaoImutavel,
		Turnaround:        []string{"Front", "Three-quarter", "Side", "Back three-quarter", "Back"},
		Expressoes:        []string{"Idle", "Active", "Highlight", "Alert", "Resolved"},
		Poses:             []string{"Entry", "Main action", "Emphasis", "Exit"},
		BuildNotes: BuildNotes{
			CorpoMaterial:     a.DescricaoImutavel,
			OlhosDetalhes:     "Clean miniature diorama surface detail",
			AcessoriosEfeitos: "Consistent lighting and scale",
			Regras:            regras,
		},
	}
}

func fichaObjetoDeDispositivo(id, dispositivo string) FichaAtivo {
	return FichaAtivo{
		ID:                id,
		Categoria:         "objeto",
		Nome:              "Objeto " + id,
		DescricaoImutavel: dispositivo,
		Turnaround:        []string{"Front", "Three-quarter", "Side", "Back three-quarter", "Back"},
		Expressoes:        []string{"Neutral", "Highlighted", "Animated peak", "Warning", "Settled"},
		Poses:             []string{"Spawn", "Idle cycle", "Key motion", "Transition"},
		BuildNotes: BuildNotes{
			CorpoMaterial:     "Miniature isometric diorama prop: " + dispositivo,
			OlhosDetalhes:     "Crisp tactile edges and readable silhouette",
			AcessoriosEfeitos: "Consistent 3D miniature scale",
			Regras:            "No text artifacts; maintain exact color and proportions",
		},
	}
}
