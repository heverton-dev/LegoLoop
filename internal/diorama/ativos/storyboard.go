package ativos

import (
	"fmt"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/diorama/roteiro"
)

// ExecutarEtapa5Storyboard compila cenas.json cobrindo 100% das linhas com metafora visual e >= 4 elementos moveis.
func ExecutarEtapa5Storyboard(projDir string, proj roteiro.Projeto, script roteiro.Script, sprites []DescritorSprite, cenarios []Cenario) (Storyboard, error) {
	if len(script.Lines) == 0 || len(cenarios) == 0 {
		return Storyboard{}, fmt.Errorf("%w: script ou cenarios vazios", ErrStoryboardInvalido)
	}
	sb := Storyboard{
		SchemaVersion: 1,
		ProjetoSlug:   proj.Slug,
		TotalLinhas:   len(script.Lines),
		Cenas:         make([]Cena, 0, len(script.Lines)),
	}
	spriteID := primeiroSpriteID(sprites)
	for i, linha := range script.Lines {
		cen := cenarios[(i/2)%len(cenarios)]
		cena := Cena{
			IndiceLinha:     i,
			Texto:           linha.Text,
			CenarioID:       cen.ID,
			MetaforaVisual:  derivarMetaforaVisual(linha),
			ElementosMoveis: comporElementosMoveis(i, spriteID),
		}
		sb.Cenas = append(sb.Cenas, cena)
	}
	if err := ValidarStoryboard(sb, script, cenarios); err != nil {
		return Storyboard{}, err
	}
	if err := gravarJSONSeguro(projDir, "cenas.json", sb); err != nil {
		return Storyboard{}, err
	}
	return sb, nil
}

// ValidarStoryboard aplica o Gate da Etapa 5: 100% das linhas, cenario valido, metafora e >= 4 elementos moveis.
func ValidarStoryboard(sb Storyboard, script roteiro.Script, cenarios []Cenario) error {
	if len(sb.Cenas) != len(script.Lines) {
		return fmt.Errorf("%w: cobertura %d/%d linhas", ErrStoryboardInvalido, len(sb.Cenas), len(script.Lines))
	}
	cenIDs := indexarCenarios(cenarios)
	for i, cena := range sb.Cenas {
		if err := validarCenaIndividual(i, cena, cenIDs); err != nil {
			return err
		}
	}
	return nil
}

func validarCenaIndividual(indiceEsperado int, cena Cena, cenIDs map[string]struct{}) error {
	if cena.IndiceLinha != indiceEsperado {
		return fmt.Errorf("%w: linha %d fora de ordem (%d)", ErrStoryboardInvalido, indiceEsperado, cena.IndiceLinha)
	}
	if _, ok := cenIDs[cena.CenarioID]; !ok {
		return fmt.Errorf("%w: linha %d referencia cenario inexistente %q", ErrStoryboardInvalido, indiceEsperado, cena.CenarioID)
	}
	if strings.TrimSpace(cena.MetaforaVisual) == "" {
		return fmt.Errorf("%w: linha %d sem metafora_visual", ErrStoryboardInvalido, indiceEsperado)
	}
	if len(cena.ElementosMoveis) < MinElementosMoveis {
		return fmt.Errorf("%w: linha %d com %d elementos moveis (minimo %d)", ErrStoryboardInvalido, indiceEsperado, len(cena.ElementosMoveis), MinElementosMoveis)
	}
	return nil
}

func indexarCenarios(cenarios []Cenario) map[string]struct{} {
	m := make(map[string]struct{}, len(cenarios))
	for _, c := range cenarios {
		m[c.ID] = struct{}{}
	}
	return m
}

func primeiroSpriteID(sprites []DescritorSprite) string {
	if len(sprites) == 0 {
		return "narrador_pose_1"
	}
	return sprites[0].ID
}

// ExtrairDispositivosPlanoMD extrai os dispositivos visuais (cenarios/objetos) da tabela do plano.md.
func ExtrairDispositivosPlanoMD(planoMD string) []string {
	var dispositivos []string
	for _, linha := range strings.Split(planoMD, "\n") {
		t := strings.TrimSpace(linha)
		if !strings.HasPrefix(t, "|") {
			continue
		}
		cols := strings.Split(strings.Trim(t, "|"), "|")
		if len(cols) < 5 {
			continue
		}
		batida := strings.TrimSpace(cols[0])
		disp := strings.TrimSpace(cols[4])
		if batida == "" || strings.EqualFold(batida, "Batida") || strings.HasPrefix(batida, ":---") || strings.HasPrefix(batida, "---") || disp == "" {
			continue
		}
		dispositivos = append(dispositivos, disp)
	}
	return dispositivos
}

func derivarMetaforaVisual(linha roteiro.LinhaRoteiro) string {
	switch linha.Natureza {
	case "fato_consolidado":
		return "quadro_evidencia_diorama"
	case "evidencia_preliminar":
		return "lupa_hipotese_camadas"
	case "alegacao_terceiro":
		return "balao_citacao_palco"
	default:
		return "transicao_ponte_narrativa"
	}
}

func comporElementosMoveis(idx int, spriteID string) []ElementoMovel {
	return []ElementoMovel{
		{ID: fmt.Sprintf("l%d_personagem", idx), Tipo: "personagem", SpriteID: spriteID, Acao: "gesto_enfatiza"},
		{ID: fmt.Sprintf("l%d_objeto", idx), Tipo: "objeto", Acao: "destaque_foco"},
		{ID: fmt.Sprintf("l%d_camera", idx), Tipo: "camera", Acao: "travelling_suave"},
		{ID: fmt.Sprintf("l%d_efeito", idx), Tipo: "efeito", Acao: "particulas_profundidade"},
	}
}
