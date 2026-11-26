package ativos

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/heverton-dev/LegoLoop/internal/video/inventario"
	"github.com/heverton-dev/LegoLoop/internal/video/roteiro"
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

// ValidarElenco exige estilo igual ao inventario e cast so com itens do estilo e partes suportadas.
func ValidarElenco(el Elenco, inv InventarioEstilo) error {
	if strings.TrimSpace(el.Estilo) == "" || el.Estilo != inv.Estilo {
		return fmt.Errorf("%w: elenco usa estilo %q, inventario %q", ErrEstiloInexistente, el.Estilo, inv.Estilo)
	}
	var erros []error
	for i, p := range el.Personagens {
		if strings.TrimSpace(p.ID) == "" {
			erros = append(erros, fmt.Errorf("%w: personagem %d sem id", ErrElencoInvalido, i))
			continue
		}
		erros = append(erros, validarCast(p, inv)...)
	}
	return errors.Join(erros...)
}

// validarCast devolve um erro por campo do cast, em ordem alfabetica de campo.
func validarCast(p PersonagemElenco, inv InventarioEstilo) []error {
	campos := make([]string, 0, len(p.Cast))
	for c := range p.Cast {
		campos = append(campos, c)
	}
	sort.Strings(campos)
	var erros []error
	for _, c := range campos {
		if err := validarCampoCast(p.ID, c, p.Cast[c], inv); err != nil {
			erros = append(erros, err)
		}
	}
	return erros
}

func validarCampoCast(id, campo, valor string, inv InventarioEstilo) error {
	if campoCor(campo) {
		return nil
	}
	tipo, forma, ok := tipoCampoCast(campo)
	if !ok {
		return fmt.Errorf("%w: personagem %q campo %q desconhecido", ErrParteNaoSuportada, id, campo)
	}
	if !inv.Tem(tipo, valor) {
		return fmt.Errorf("%w: personagem %q campo %q item %q inexistente no estilo %q", ErrEstiloInexistente, id, campo, valor, inv.Estilo)
	}
	if forma && !inv.Suporta(valor) {
		return fmt.Errorf("%w: personagem %q campo %q parte %q fora de partes_suportadas", ErrParteNaoSuportada, id, campo, valor)
	}
	return nil
}

// ValidarStoryboardV2 exige personagem do elenco e pose e humor do inventario; um erro por elemento.
func ValidarStoryboardV2(sb Storyboard, el Elenco, inv InventarioEstilo) error {
	ids := make(map[string]struct{}, len(el.Personagens))
	for _, p := range el.Personagens {
		ids[p.ID] = struct{}{}
	}
	var erros []error
	for _, cena := range sb.Cenas {
		for _, em := range cena.ElementosMoveis {
			if em.Tipo != "personagem" {
				continue
			}
			if err := validarElementoV2(cena.IndiceLinha, em, ids, inv); err != nil {
				erros = append(erros, err)
			}
		}
	}
	return errors.Join(erros...)
}

func validarElementoV2(linha int, em ElementoMovel, ids map[string]struct{}, inv InventarioEstilo) error {
	if _, ok := ids[em.PersonagemID]; !ok {
		return fmt.Errorf("%w: %w: linha %d elemento %q personagem_id %q", ErrStoryboardInvalido, ErrPersonagemForaDoElenco, linha, em.ID, em.PersonagemID)
	}
	if !inv.Tem("pose", em.Pose) {
		return fmt.Errorf("%w: %w: linha %d elemento %q pose %q", ErrStoryboardInvalido, ErrItemForaDoInventario, linha, em.ID, em.Pose)
	}
	if !inv.Tem("humor", em.Humor) {
		return fmt.Errorf("%w: %w: linha %d elemento %q humor %q", ErrStoryboardInvalido, ErrItemForaDoInventario, linha, em.ID, em.Humor)
	}
	return nil
}

// CarregarInventarioEstilo le estilo.json e os itens aprovados dos tipos usados por cast e cenas.
func CarregarInventarioEstilo(s inventario.Store, estilo string) (InventarioEstilo, error) {
	e, err := s.LerEstilo(estilo)
	if errors.Is(err, inventario.ErrEstiloInexistente) {
		return InventarioEstilo{}, fmt.Errorf("%w: %w", ErrEstiloInexistente, err)
	}
	if err != nil {
		return InventarioEstilo{}, err
	}
	inv := InventarioEstilo{Estilo: e.ID, PartesSuportadas: e.PartesSuportadas, Itens: map[string][]string{}}
	for _, tipo := range []string{"pele", "traje", "penteado", "adorno", "aderecos", "pose", "humor"} {
		itens, err := s.Buscar(estilo, tipo, nil)
		if err != nil {
			return InventarioEstilo{}, err
		}
		for _, it := range itens {
			inv.Itens[tipo] = append(inv.Itens[tipo], it.ID)
		}
	}
	return inv, nil
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
