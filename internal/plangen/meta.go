package plangen

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// planMeta reúne o que PLAN e ROUTE compartilham: é a fonte única dos dois, o que os mantém espelhados.
type planMeta struct {
	Date                            string
	Theme                           string
	Level                           SupervisionLevel
	SpecRel, PlanRel, RouteRel      string // relativos à raiz do projeto, com "/"
	IntentRef                       string // idem; vazio quando a SPEC não informa
	SpecLink, IntentLink, RouteLink string // relativos a docs/plans (links do PLAN)
	SpecFile, RouteName             string
	Epics                           map[int]string
	SpecCriteria                    []string
}

var (
	specThemeRe  = regexp.MustCompile(`(?m)^# 📋 SPEC:\s*(.+?)\s*$`)
	intentRefRe  = regexp.MustCompile(`Referência da Intenção:\*\*\s*(?:\[([^\]]+)\]|(\S+))`)
	epicNameRe   = regexp.MustCompile(`(?m)^#{2,4}\s*(?:Épico\s+|E)?(\d+)\s*[—–-]\s*(.+?)\s*$`)
	specACRe     = regexp.MustCompile(`(?i)\*\s*\[\s*\]\s*(AC-[A-Za-z0-9_.-]+)`)
	firstNumbers = regexp.MustCompile(`\d+`)
)

type pathSet struct {
	baseDir, specPath, planPath, routePath string
}

// rel devolve p relativo à raiz do projeto, com "/". Fora da raiz, devolve só o nome do arquivo.
func (p pathSet) rel(target string) string {
	abs, err := filepath.Abs(target)
	if err != nil {
		return filepath.ToSlash(filepath.Base(target))
	}
	base, _ := filepath.Abs(p.baseDir)
	r, err := filepath.Rel(base, abs)
	if err != nil || strings.HasPrefix(r, "..") {
		return filepath.ToSlash(filepath.Base(target))
	}
	return filepath.ToSlash(r)
}

// link devolve o caminho relativo de docs/plans até o alvo (para links Markdown).
func (p pathSet) link(projectRel string) string {
	from := filepath.Join(p.baseDir, "docs", "plans")
	r, err := filepath.Rel(from, filepath.Join(p.baseDir, filepath.FromSlash(projectRel)))
	if err != nil {
		return projectRel
	}
	return filepath.ToSlash(r)
}

func newMeta(ps pathSet, specText, date string, level SupervisionLevel) planMeta {
	m := planMeta{Date: date, Level: level, SpecFile: filepath.Base(ps.specPath), Epics: map[int]string{}}
	m.SpecRel, m.PlanRel, m.RouteRel = ps.rel(ps.specPath), ps.rel(ps.planPath), ps.rel(ps.routePath)
	m.RouteName = strings.TrimSuffix(filepath.Base(ps.routePath), ".json")
	m.Theme = strings.TrimSuffix(m.SpecFile, ".md")
	if t := specThemeRe.FindStringSubmatch(specText); len(t) > 1 {
		m.Theme = t[1]
	}
	m.IntentRef = intentRefOf(ps, specText)
	m.SpecLink, m.RouteLink = ps.link(m.SpecRel), ps.link(m.RouteRel)
	if m.IntentRef != "" {
		m.IntentLink = ps.link(m.IntentRef)
	}
	for _, e := range epicNameRe.FindAllStringSubmatch(specText, -1) {
		if n, err := strconv.Atoi(e[1]); err == nil {
			m.Epics[n] = e[2]
		}
	}
	for _, c := range specACRe.FindAllStringSubmatch(specText, -1) {
		m.SpecCriteria = append(m.SpecCriteria, c[1])
	}
	return m
}

func intentRefOf(ps pathSet, specText string) string {
	m := intentRefRe.FindStringSubmatch(specText)
	if m == nil {
		return ""
	}
	v := strings.TrimSpace(m[1] + m[2])
	v = strings.ReplaceAll(v, `\\`, "/")
	v = strings.ReplaceAll(v, `\`, "/")
	if filepath.IsAbs(filepath.FromSlash(v)) {
		return ps.rel(v)
	}
	return v
}

// epicName devolve o nome do épico segundo a SPEC, ou "Épico E01" quando a SPEC não o nomeia.
func (m planMeta) epicName(epicID string) string {
	if n := epicNumber(epicID); n > 0 {
		if name, ok := m.Epics[n]; ok {
			return name
		}
	}
	return "Épico " + epicID
}

func epicNumber(epicID string) int {
	d := firstNumbers.FindString(epicID)
	n, _ := strconv.Atoi(d)
	return n
}
