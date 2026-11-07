package legopress

import (
	"errors"
	"strings"
)

var (
	ErrManifestoInvalido = errors.New("legopress: manifesto invalido")
	ErrGeometriaFalhou   = errors.New("legopress: falha na auditoria geometrica")
	ErrModoDesconhecido  = errors.New("legopress: modo operacional desconhecido")
)

type Mode string

const (
	ModeCreate Mode = "create"
	ModeEdit   Mode = "edit"
	ModeSync   Mode = "sync"
	ModeAdapt  Mode = "adapt"
	ModeAudit  Mode = "audit"
)

func (m Mode) Validate() error {
	switch m {
	case ModeCreate, ModeEdit, ModeSync, ModeAdapt, ModeAudit:
		return nil
	default:
		return ErrModoDesconhecido
	}
}

type Part struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	File  string `json:"file"`
	Level int    `json:"level"`
	Words int    `json:"words,omitempty"`
}

type Manifest struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle,omitempty"`
	Author    string `json:"author,omitempty"`
	Format    string `json:"format,omitempty"`
	Template  string `json:"template,omitempty"`
	Parts     []Part `json:"parts"`
	OutputPDF string `json:"output_pdf,omitempty"`
}

func (m Manifest) Validate() error {
	if strings.TrimSpace(m.Title) == "" || len(m.Parts) == 0 {
		return ErrManifestoInvalido
	}
	for _, p := range m.Parts {
		if strings.TrimSpace(p.File) == "" {
			return ErrManifestoInvalido
		}
	}
	return nil
}

type GeometryReport struct {
	TotalPages      int   `json:"total_pages"`
	OrphanPages     []int `json:"orphan_pages,omitempty"`
	BrokenHeaders   []int `json:"broken_headers,omitempty"`
	MultiLineTables []int `json:"multiline_tables,omitempty"`
	Valid           bool  `json:"valid"`
}
