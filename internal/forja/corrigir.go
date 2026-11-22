package forja

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/jail"
)

// ItemCorrigido registra uma skill eager movida para o Cofre com backup.
type ItemCorrigido struct {
	Harness string `json:"harness"`
	Skill   string `json:"skill"`
	Origem  string `json:"origem"`
	Destino string `json:"destino"`
	Backup  string `json:"backup"`
}

// ResultadoCorrecao resume as derivas corrigidas e seus backups.
type ResultadoCorrecao struct {
	Data       string          `json:"data"`
	BackupRaiz string          `json:"backup_raiz"`
	Cofre      string          `json:"cofre"`
	Corrigidos []ItemCorrigido `json:"corrigidos,omitempty"`
}

// CorrigirDeriva faz backup das skills eager encontradas, copia cada uma para
// o cofreDir e remove a pasta eager do harness.
func CorrigirDeriva(home, cofreDir, data string) (ResultadoCorrecao, error) {
	if strings.TrimSpace(home) == "" {
		return ResultadoCorrecao{}, ErrSemHome
	}
	if strings.TrimSpace(cofreDir) == "" {
		cofreDir = filepath.Join(home, ".aidd", "canonical", "skills")
	}
	if strings.TrimSpace(data) == "" {
		data = time.Now().UTC().Format("2006-01-02")
	}
	backupRaiz := filepath.Join(home, ".aidd", "backups", "deriva_"+data)
	res := ResultadoCorrecao{Data: data, BackupRaiz: backupRaiz, Cofre: cofreDir}
	derivas, err := DetectarDeriva(home)
	if err != nil || len(derivas) == 0 {
		return res, err
	}
	grupos, ordem := agruparDerivasPorHarness(derivas)
	for _, h := range ordem {
		itens, err := corrigirHarnessDeriva(home, cofreDir, backupRaiz, h, grupos[h])
		if err != nil {
			return res, err
		}
		res.Corrigidos = append(res.Corrigidos, itens...)
	}
	return res, nil
}

func agruparDerivasPorHarness(derivas []DerivaItem) (map[string][]DerivaItem, []string) {
	grupos := make(map[string][]DerivaItem)
	var ordem []string
	for _, d := range derivas {
		if _, ok := grupos[d.Harness]; !ok {
			ordem = append(ordem, d.Harness)
		}
		grupos[d.Harness] = append(grupos[d.Harness], d)
	}
	return grupos, ordem
}

func corrigirHarnessDeriva(home, cofreDir, backupRaiz, harness string, itens []DerivaItem) ([]ItemCorrigido, error) {
	dirBackup := filepath.Join(backupRaiz, harness)
	if err := os.MkdirAll(dirBackup, 0o755); err != nil {
		return nil, fmt.Errorf("forja: criar backup %s: %w", dirBackup, err)
	}
	var corrigidos []ItemCorrigido
	var manifesto strings.Builder
	manifesto.WriteString("acao\torigem\tbackup\tdestino\n")
	for _, d := range itens {
		item, err := moverSkillDeriva(home, cofreDir, dirBackup, d)
		if err != nil {
			return nil, err
		}
		corrigidos = append(corrigidos, item)
		manifesto.WriteString(fmt.Sprintf("movido\t%s\t%s\t%s\n",
			filepath.ToSlash(item.Origem), filepath.ToSlash(item.Backup), filepath.ToSlash(item.Destino)))
	}
	if err := os.WriteFile(filepath.Join(dirBackup, "MANIFESTO.tsv"), []byte(manifesto.String()), 0o644); err != nil {
		return nil, fmt.Errorf("forja: gravar MANIFESTO.tsv: %w", err)
	}
	return corrigidos, nil
}

func moverSkillDeriva(home, cofreDir, dirBackup string, d DerivaItem) (ItemCorrigido, error) {
	origemSegura, err := jail.DentroDaRaiz(home, d.Pasta)
	if err != nil {
		return ItemCorrigido{}, err
	}
	pastaBackup := filepath.Join(dirBackup, d.Skill)
	if err := copiarDiretorio(origemSegura, pastaBackup); err != nil {
		return ItemCorrigido{}, fmt.Errorf("forja: backup da skill %s: %w", d.Skill, err)
	}
	destinoCofre := filepath.Join(cofreDir, d.Skill)
	if err := copiarDiretorio(origemSegura, destinoCofre); err != nil {
		return ItemCorrigido{}, fmt.Errorf("forja: copiar skill %s ao cofre: %w", d.Skill, err)
	}
	if err := os.RemoveAll(origemSegura); err != nil {
		return ItemCorrigido{}, fmt.Errorf("forja: remover skill eager %s: %w", origemSegura, err)
	}
	return ItemCorrigido{
		Harness: d.Harness,
		Skill:   d.Skill,
		Origem:  origemSegura,
		Destino: destinoCofre,
		Backup:  pastaBackup,
	}, nil
}

func copiarDiretorio(origem, destino string) error {
	return filepath.WalkDir(origem, func(caminho string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(origem, caminho)
		if err != nil {
			return err
		}
		alvo := filepath.Join(destino, rel)
		if d.IsDir() {
			return os.MkdirAll(alvo, 0o755)
		}
		return copiarArquivoSimples(caminho, alvo)
	})
}

func copiarArquivoSimples(origem, destino string) error {
	in, err := os.Open(origem)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(destino), 0o755); err != nil {
		return err
	}
	out, err := os.Create(destino)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		if cerr := out.Close(); cerr != nil {
			return cerr
		}
		return err
	}
	return out.Close()
}
