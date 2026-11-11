package legopress

import (
	"errors"
	"fmt"
	"strings"
)

var ErrTopicosVazios = errors.New("legopress: lista de topicos vazia")

type RawTopic struct {
	Title       string `json:"title"`
	ContentType string `json:"content_type"`
}

type Partitioner struct{}

func NewPartitioner() *Partitioner {
	return &Partitioner{}
}

func (p *Partitioner) Partition(topics []RawTopic) ([]Part, error) {
	if len(topics) == 0 {
		return nil, ErrTopicosVazios
	}

	parts := make([]Part, 0, len(topics))
	for i, t := range topics {
		if strings.TrimSpace(t.Title) == "" {
			return nil, fmt.Errorf("legopress: topico %d com titulo vazio", i)
		}
		level := classifyLevel(t.ContentType)
		file := fmt.Sprintf("partes/%02d-%s.md", i, sanitizeSlug(t.Title))
		parts = append(parts, Part{
			ID:    fmt.Sprintf("P%02d", i+1),
			Title: t.Title,
			File:  file,
			Level: level,
		})
	}
	return parts, nil
}

func classifyLevel(contentType string) int {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "concept", "intro", "manifesto":
		return 1
	case "code", "architecture", "internals", "micro", "governance":
		return 3
	default:
		return 2
	}
}

func sanitizeSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
}
