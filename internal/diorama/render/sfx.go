package render

import (
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

var ErrSFXSemArquivo = errors.New("diorama-render: evento de SFX sem arquivo no catalogo")

// AtrasoObjetoS e o intervalo entre o inicio da fala e o destaque do objeto.
const AtrasoObjetoS = 0.25

// CatalogoSFX mapeia o tipo de evento ao arquivo de audio, relativo a pasta do projeto.
type CatalogoSFX map[string]string

// CatalogoPadrao aponta para os SFX que o pacote sintetiza quando o projeto nao traz os seus.
func CatalogoPadrao() CatalogoSFX {
	return CatalogoSFX{
		"troca_cena":       "audio/sfx/troca_cena.wav",
		"movimento_objeto": "audio/sfx/movimento_objeto.wav",
		"destaque":         "audio/sfx/destaque.wav",
	}
}

// CarregarCatalogo usa o catalogo de render.json quando houver; senao, o padrao.
func CarregarCatalogo(projDir string, cfg ConfigRender) (CatalogoSFX, error) {
	if cfg.SFXCatalogo == "" {
		return CatalogoPadrao(), nil
	}
	var cat CatalogoSFX
	if err := lerJSON(projDir, cfg.SFXCatalogo, &cat, true); err != nil {
		return nil, fmt.Errorf("%w: sfx_catalogo %q: %v", ErrSFXSemArquivo, cfg.SFXCatalogo, err)
	}
	for tipo := range cat {
		if !sfxValido(tipo) {
			return nil, invalido("catalogo de SFX com tipo %q desconhecido", tipo)
		}
	}
	return cat, nil
}

type eventoPrevisto struct {
	tipo  string
	tempo float64
}

func temCamada(b Batida, nome string) bool {
	for _, c := range b.Camadas {
		if c.Camada == nome {
			return true
		}
	}
	return false
}

// eventosDaBatida aplica as regras: troca de cenario, destaque por shake e objeto na 1a batida da fala.
func eventosDaBatida(b Batida, cenarioAnterior string, linhaAnterior int, primeira bool) []eventoPrevisto {
	var ev []eventoPrevisto
	if primeira || b.CenarioID != cenarioAnterior {
		ev = append(ev, eventoPrevisto{"troca_cena", b.InicioS})
	}
	if b.Camera == "shake" {
		ev = append(ev, eventoPrevisto{"destaque", b.InicioS})
	}
	if (primeira || b.IndiceLinha != linhaAnterior) && temCamada(b, "aderecos") {
		ev = append(ev, eventoPrevisto{"movimento_objeto", math.Min(b.InicioS+AtrasoObjetoS, b.FimS)})
	}
	sort.SliceStable(ev, func(i, j int) bool { return ev[i].tempo < ev[j].tempo })
	return ev
}

func resolverSFX(projDir string, cat CatalogoSFX, tipo string) (string, error) {
	rel := strings.TrimSpace(cat[tipo])
	if rel == "" {
		return "", fmt.Errorf("%w: tipo %q sem entrada no catalogo", ErrSFXSemArquivo, tipo)
	}
	caminho, err := resolverSeguro(projDir, rel)
	if err != nil {
		return "", fmt.Errorf("%w: %q: %v", ErrSFXSemArquivo, rel, err)
	}
	if fi, err := os.Stat(caminho); err != nil || fi.IsDir() {
		return "", fmt.Errorf("%w: arquivo %q ausente", ErrSFXSemArquivo, rel)
	}
	return rel, nil
}

func resolverEventos(projDir string, cat CatalogoSFX, previstos []eventoPrevisto) ([]EventoSFX, error) {
	var sfx []EventoSFX
	for _, p := range previstos {
		rel, err := resolverSFX(projDir, cat, p.tipo)
		if err != nil {
			return nil, err
		}
		sfx = append(sfx, EventoSFX{Tipo: p.tipo, TempoS: p.tempo, Arquivo: rel})
	}
	return sfx, nil
}

// GerarEventosSFX substitui os eventos de cada batida pelos derivados da propria timeline.
// Evento sem entrada ou sem arquivo no catalogo bloqueia. A entrada nao muda.
func GerarEventosSFX(projDir string, tl Timeline, cat CatalogoSFX) (Timeline, error) {
	saida := Timeline{SchemaVersion: tl.SchemaVersion, DuracaoS: tl.DuracaoS, Batidas: make([]Batida, len(tl.Batidas))}
	anterior, linhaAnterior := "", -1
	for i, b := range tl.Batidas {
		b.Camadas = append([]CamadaItem(nil), b.Camadas...)
		sfx, err := resolverEventos(projDir, cat, eventosDaBatida(b, anterior, linhaAnterior, i == 0))
		if err != nil {
			return Timeline{}, err
		}
		b.SFX = sfx
		saida.Batidas[i] = b
		anterior, linhaAnterior = b.CenarioID, b.IndiceLinha
	}
	return saida, nil
}
