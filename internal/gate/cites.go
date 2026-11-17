package gate

import (
	"bytes"
	"errors"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// maxCitacoes limita o trabalho por artefato: acima disso o gate reprova em vez de varrer sem fim.
const maxCitacoes = 500

// minIdent descarta identificadores curtos demais para provar algo ("ok", "id", "go").
const minIdent = 3

// bt e o acento grave, que nao cabe dentro de string bruta.
const bt = "`"

var (
	citePathRe = regexp.MustCompile(`(?:[A-Za-z0-9_.-]+/)*[A-Za-z0-9_.-]+\.(?:go|md|json|sh|py|js|ts|tsx|css|html|ya?ml|toml|mod|sql)`)
	// citeColonRe: o numero vem logo apos o caminho, como "arquivo.go:12" ou "arquivo.go:12-15".
	citeColonRe = regexp.MustCompile(`^[` + bt + `')\]]*:(\d+)(?:\s*[-–]\s*(\d+))?`)
	// citeLooseRe: "(linha 12)", "linhas 12-15" ou "line 12" ate 80 caracteres depois do caminho, antes do proximo caminho.
	citeLooseRe = regexp.MustCompile(`(?i)^.{0,80}?(?:linhas?|lines?)\s*(\d+)(?:\s*[-–]\s*(\d+))?`)
	// citeSpanRe acha os trechos entre crases de uma linha; citeIdentRe aceita so nomes de codigo.
	citeSpanRe  = regexp.MustCompile(bt + `([^` + bt + `\n]+)` + bt)
	citeIdentRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)*(?:\(\))?$`)
)

var errMuitasCitacoes = errors.New("gate: citacoes demais no artefato")

type citation struct {
	file       string
	from, to   int
	artifactLn int
}

// extractCitations acha, linha a linha do texto, cada "caminho + numero de linha". O numero pertence ao caminho
// mais proximo a sua esquerda na mesma linha. Caminho sem numero nao e citacao.
func extractCitations(text string) []citation {
	var out []citation
	for ln, line := range strings.Split(text, "\n") {
		idx := citePathRe.FindAllStringIndex(line, -1)
		for i, loc := range idx {
			rest := line[loc[1]:]
			if i+1 < len(idx) {
				rest = line[loc[1]:idx[i+1][0]]
			}
			m := citeColonRe.FindStringSubmatch(rest)
			if m == nil {
				m = citeLooseRe.FindStringSubmatch(rest)
			}
			if m == nil {
				continue
			}
			from, _ := strconv.Atoi(m[1])
			to := from
			if m[2] != "" {
				to, _ = strconv.Atoi(m[2])
			}
			out = append(out, citation{file: line[loc[0]:loc[1]], from: from, to: to, artifactLn: ln + 1})
		}
	}
	return out
}

// CiteVerdict e o desfecho de uma citacao "arquivo:linha" do artefato.
type CiteVerdict struct {
	File         string   `json:"file"`
	From         int      `json:"from"`
	To           int      `json:"to"`
	ArtifactLine int      `json:"artifact_line"`
	InRange      bool     `json:"in_range"`      // o arquivo existe e o intervalo cabe nele
	MissingIdent []string `json:"missing_ident"` // identificadores da mesma linha ausentes dos arquivos citados nela
}

// OK diz se a citacao passa: intervalo valido e nenhum identificador ausente.
func (v CiteVerdict) OK() bool { return v.InRange && len(v.MissingIdent) == 0 }

type citeFile struct {
	data  []byte
	lines int
	ok    bool
}

func loadCiteFile(root RootFS, name string, cache map[string]citeFile) citeFile {
	if f, hit := cache[name]; hit {
		return f
	}
	f := citeFile{}
	if data, err := readGateFile(Env{Root: root}, name); err == nil {
		f = citeFile{data: data, ok: true, lines: bytes.Count(data, []byte("\n"))}
		if len(data) > 0 && data[len(data)-1] != '\n' {
			f.lines++
		}
	}
	cache[name] = f
	return f
}

// lineIdents devolve os identificadores entre crases de uma linha do artefato: nomes de codigo, sem caminhos,
// sem numeros de linha e com ao menos minIdent caracteres. Em "pkg.Func" vale o ultimo trecho; "Func()" vale "Func".
func lineIdents(line string) []string {
	var out []string
	seen := map[string]bool{}
	for _, m := range citeSpanRe.FindAllStringSubmatch(line, -1) {
		tok := strings.TrimSpace(m[1])
		if full := citePathRe.FindString(tok); full == tok || !citeIdentRe.MatchString(tok) {
			continue
		}
		tok = strings.TrimSuffix(tok, "()")
		if i := strings.LastIndex(tok, "."); i >= 0 {
			tok = tok[i+1:]
		}
		if len(tok) >= minIdent && !seen[tok] {
			seen[tok] = true
			out = append(out, tok)
		}
	}
	return out
}

// CiteVerdicts avalia cada citacao "arquivo:linha" do artefato contra a pasta base, dentro de root. Alem de
// existir, a citacao exige que cada identificador entre crases na mesma linha do artefato apareca (como palavra)
// em ao menos um dos arquivos citados nessa linha. Prova existencia e coerencia de nomes, nao o sentido do texto.
func CiteVerdicts(root RootFS, artifact, base string) ([]CiteVerdict, error) {
	data, err := readGateFile(Env{Root: root}, artifact)
	if err != nil {
		return nil, err
	}
	cites := extractCitations(string(data))
	if len(cites) > maxCitacoes {
		return nil, errMuitasCitacoes
	}
	base = path.Clean(filepath.ToSlash(base))
	cache := map[string]citeFile{}
	artLines := strings.Split(string(data), "\n")
	verdicts, filesByLine, firstByLine := montarVerdicts(root, base, cites, cache)
	marcarIdentsFaltantes(verdicts, artLines, filesByLine, firstByLine)
	return verdicts, nil
}

// montarVerdicts carrega cada arquivo citado e monta o veredito inicial de cada citacao, agrupando os arquivos
// vivos por linha do artefato e registrando a primeira citacao de cada linha.
func montarVerdicts(root RootFS, base string, cites []citation, cache map[string]citeFile) ([]CiteVerdict, map[int][]citeFile, map[int]int) {
	verdicts := make([]CiteVerdict, len(cites))
	filesByLine := map[int][]citeFile{}
	firstByLine := map[int]int{}
	for i, c := range cites {
		target := path.Join(base, c.file)
		f := citeFile{}
		if filepath.IsLocal(filepath.FromSlash(target)) {
			f = loadCiteFile(root, target, cache)
		}
		verdicts[i] = CiteVerdict{
			File: c.file, From: c.from, To: c.to, ArtifactLine: c.artifactLn,
			InRange: f.ok && c.from >= 1 && c.to >= c.from && f.lines >= c.to,
		}
		if f.ok {
			filesByLine[c.artifactLn] = append(filesByLine[c.artifactLn], f)
		}
		if _, seen := firstByLine[c.artifactLn]; !seen {
			firstByLine[c.artifactLn] = i
		}
	}
	return verdicts, filesByLine, firstByLine
}

// marcarIdentsFaltantes anota em MissingIdent da primeira citacao da linha cada identificador que nao aparece
// em nenhum dos arquivos citados nessa linha.
func marcarIdentsFaltantes(verdicts []CiteVerdict, artLines []string, filesByLine map[int][]citeFile, firstByLine map[int]int) {
	for ln, first := range firstByLine {
		for _, id := range lineIdents(artLines[ln-1]) {
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(id) + `\b`)
			found := false
			for _, f := range filesByLine[ln] {
				if re.Match(f.data) {
					found = true
					break
				}
			}
			if !found {
				verdicts[first].MissingIdent = append(verdicts[first].MissingIdent, id)
			}
		}
	}
}

// citesExist e o builtin cites_exist(artefato, base): verdadeiro so se o artefato traz ao menos uma citacao
// "arquivo:linha" e todas passam em CiteVerdicts (existem e os identificadores da linha aparecem no arquivo citado).
// Citacao fora de base ou com ".." conta como inexistente.
func citesExist(env Env, args []value, _ *regexp.Regexp) (value, error) {
	verdicts, err := CiteVerdicts(env.Root, expandGateVars(env, args[0].s), expandGateVars(env, args[1].s))
	if errors.Is(err, errMuitasCitacoes) {
		return value{kind: kBool}, nil
	}
	if err != nil {
		return value{}, err
	}
	if len(verdicts) == 0 {
		return value{kind: kBool}, nil
	}
	for _, v := range verdicts {
		if !v.OK() {
			return value{kind: kBool}, nil
		}
	}
	return value{kind: kBool, b: true}, nil
}
