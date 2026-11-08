package redc

import (
	"sort"
	"strings"
)

// MethodCluster agrupa vídeos que ensinam o mesmo método ou metodologia.
type MethodCluster struct {
	Slug        string   `json:"slug"`
	Metodo      string   `json:"metodo"`
	ClusterTema string   `json:"cluster_tema"`
	VideoIDs    []string `json:"video_ids"`
	ScoreTotal  float64  `json:"score_total"`
}

// ClusterMethods consolida os métodos dos vídeos eliminando redundâncias.
func ClusterMethods(videos []VideoEntry) []MethodCluster {
	groups := make(map[string]*MethodCluster)
	for _, v := range videos {
		key := strings.ToLower(strings.TrimSpace(v.Metodo))
		if key == "" {
			key = "geral"
		}
		c, ok := groups[key]
		if !ok {
			c = &MethodCluster{
				Slug:        slugify(v.Metodo),
				Metodo:      v.Metodo,
				ClusterTema: v.ClusterTema,
			}
			groups[key] = c
		}
		c.VideoIDs = append(c.VideoIDs, v.ID)
		c.ScoreTotal += v.Score
	}
	return toSortedClusterList(groups)
}

func toSortedClusterList(groups map[string]*MethodCluster) []MethodCluster {
	list := make([]MethodCluster, 0, len(groups))
	for _, c := range groups {
		list = append(list, *c)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].ScoreTotal > list[j].ScoreTotal
	})
	return list
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
