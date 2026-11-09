package redc

import (
	"encoding/json"
	"sort"
	"strings"
)

// CatalogVideos analisa e ordena os vídeos pelo score de tração e viralidade.
func CatalogVideos(raw []byte) ([]VideoEntry, error) {
	var list []VideoEntry
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, ErrDadosCorrompidos
	}
	if len(list) == 0 {
		return nil, ErrVideosInsuf
	}
	for i := range list {
		sanitizeEntry(&list[i])
		list[i].Score = calculateTractionScore(list[i])
	}
	sortVideosByScore(list)
	return list, nil
}

func sanitizeEntry(v *VideoEntry) {
	v.ID = strings.TrimSpace(v.ID)
	v.Titulo = strings.TrimSpace(v.Titulo)
	v.Gancho3s = strings.TrimSpace(v.Gancho3s)
	v.Metodo = strings.TrimSpace(v.Metodo)
	v.ClusterTema = strings.TrimSpace(v.ClusterTema)
}

func calculateTractionScore(v VideoEntry) float64 {
	engajamento := float64(v.Likes*2 + v.Comentarios*5)
	return float64(v.Views)*0.7 + engajamento*0.3
}

func sortVideosByScore(list []VideoEntry) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Score > list[j].Score
	})
}
