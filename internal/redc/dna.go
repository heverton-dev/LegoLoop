package redc

import (
	"fmt"
	"strings"
)

// ExtractDNA sintetiza o DNA comportamental e visual a partir dos vídeos de topo.
func ExtractDNA(canalSlug string, topVideos []VideoEntry) (ChannelDNA, error) {
	if len(topVideos) == 0 {
		return ChannelDNA{}, ErrVideosInsuf
	}
	slug := strings.TrimSpace(canalSlug)
	if slug == "" {
		slug = "canal-analisado"
	}

	dna := ChannelDNA{
		CanalSlug:    slug,
		TomDeVoz:     inferTomDeVoz(topVideos),
		Arquetipo:    "Mestre / Especialista Prático",
		RitmoCortes:  calculateAveragePace(topVideos),
		DesignTokens: defaultDesignTokens(),
		MotionTokens: defaultMotionTokens(),
		TopGanchos:   collectHooks(topVideos),
	}
	return dna, nil
}

func inferTomDeVoz(videos []VideoEntry) string {
	if len(videos) > 0 && strings.Contains(strings.ToLower(videos[0].Gancho3s), "segredo") {
		return "Magnético, revelador e focado em alta conversão"
	}
	return "Direto ao ponto, autoritário e dinâmico"
}

func calculateAveragePace(videos []VideoEntry) string {
	var totalDur int
	for _, v := range videos {
		totalDur += v.DuracaoSeg
	}
	avg := 50
	if len(videos) > 0 && totalDur > 0 {
		avg = totalDur / len(videos)
	}
	if avg <= 45 {
		return "Cortes ultra-rápidos a cada 1.8s a 2.5s"
	}
	return fmt.Sprintf("Ritmo dinâmico com corte e zoom a cada 3.0s (duração média %ds)", avg)
}

func defaultDesignTokens() map[string]string {
	return map[string]string{
		"cor_primaria":      "#0F172A",
		"cor_destaque":      "#F59E0B",
		"cor_fundo":         "#050505",
		"fonte_titulos":     "Anton, Impact, sans-serif",
		"fonte_legendas":    "Montserrat ExtraBold, sans-serif",
		"tamanho_destaque":  "72px",
		"padding_container": "32px",
	}
}

func defaultMotionTokens() map[string]string {
	return map[string]string{
		"transicao":         "smooth-snap (curva bezier acelerada 200ms)",
		"efeito_cartoes":    "vidro_3d (blur 16px, borda 1px neon, inclinacao 12deg)",
		"zoom_enfase":       "punch_in 1.15x nos primeiros 0.8s",
		"tipografia_motion": "letras_cineticas_palavra_por_palavra",
	}
}

func collectHooks(videos []VideoEntry) []string {
	hooks := make([]string, 0, len(videos))
	for _, v := range videos {
		if h := strings.TrimSpace(v.Gancho3s); h != "" {
			hooks = append(hooks, h)
		}
	}
	return hooks
}
