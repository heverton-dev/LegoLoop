package legopress

import (
	"os/exec"
	"strings"
)

type DoctorReport struct {
	PandocAvailable bool   `json:"pandoc_available"`
	PandocVersion   string `json:"pandoc_version,omitempty"`
	TypstAvailable  bool   `json:"typst_available"`
	TypstVersion    string `json:"typst_version,omitempty"`
}

type Compiler struct {
	pandocBin string
	typstBin  string
	run       func(name string, args ...string) ([]byte, error)
}

func NewCompiler() *Compiler {
	pandoc, _ := exec.LookPath("pandoc")
	typst, _ := exec.LookPath("typst")
	return &Compiler{
		pandocBin: pandoc,
		typstBin:  typst,
	}
}

func NewCompilerWithRunner(pandocBin, typstBin string, run func(string, ...string) ([]byte, error)) *Compiler {
	return &Compiler{
		pandocBin: pandocBin,
		typstBin:  typstBin,
		run:       run,
	}
}

func (c *Compiler) BuildArgs(markdownPath string, m Manifest) []string {
	template := m.Template
	if template == "" {
		template = "ativos/livro.typst"
	}
	output := m.OutputPDF
	if output == "" {
		output = "livro.pdf"
	}
	return []string{
		"pandoc",
		markdownPath,
		"--template", template,
		"--pdf-engine", "typst",
		"-o", output,
	}
}

func (c *Compiler) Doctor() *DoctorReport {
	run := c.run
	if run == nil {
		run = func(name string, args ...string) ([]byte, error) {
			return exec.Command(name, args...).Output()
		}
	}
	report := &DoctorReport{}
	if c.pandocBin != "" {
		report.PandocAvailable = true
		if out, err := run(c.pandocBin, "--version"); err == nil {
			report.PandocVersion = extractFirstLine(string(out))
		}
	}
	if c.typstBin != "" {
		report.TypstAvailable = true
		if out, err := run(c.typstBin, "--version"); err == nil {
			report.TypstVersion = extractFirstLine(string(out))
		}
	}
	return report
}

func (c *Compiler) Compile(m Manifest) error {
	if err := m.Validate(); err != nil {
		return err
	}
	return nil
}

func extractFirstLine(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}
