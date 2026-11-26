package kraften

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Preflight extrai o pacote com segurança em dest e valida os 3 artefatos
// universais: Dockerfile (com HEALTHCHECK), docker-stack.yml na rede
// network_conexao e .env.example com ADMIN_USER/ADMIN_PASS.
func Preflight(zipPath, dest string) error {
	raiz, err := extrairZipSeguro(zipPath, dest)
	if err != nil {
		return err
	}
	return validarArtefatos(raiz)
}

// extrairZipSeguro valida todos os nomes antes de gravar qualquer byte.
func extrairZipSeguro(zipPath, dest string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", ErrInvalidDeployInput
	}
	defer r.Close()
	for _, f := range r.File {
		if _, err := nomeSeguro(f.Name); err != nil {
			return "", err
		}
	}
	for _, f := range r.File {
		if err := extrairEntrada(f, dest); err != nil {
			return "", err
		}
	}
	return encontrarRaizApp(dest)
}

// nomeSeguro normaliza separadores e recusa caminhos que escapem da raiz.
func nomeSeguro(nome string) (string, error) {
	n := strings.ReplaceAll(nome, "\\", "/")
	if n == "" || strings.HasPrefix(n, "/") || filepath.IsAbs(n) {
		return "", ErrInvalidDeployInput
	}
	limpo := filepath.ToSlash(filepath.Clean(n))
	if limpo == ".." || strings.HasPrefix(limpo, "../") || strings.Contains(limpo, ":") {
		return "", ErrInvalidDeployInput
	}
	return limpo, nil
}

// extrairEntrada grava uma entrada do zip estritamente dentro de dest.
func extrairEntrada(f *zip.File, dest string) error {
	limpo, _ := nomeSeguro(f.Name)
	alvo := filepath.Join(dest, filepath.FromSlash(limpo))
	if f.FileInfo().IsDir() {
		return os.MkdirAll(alvo, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(alvo), 0o755); err != nil {
		return err
	}
	src, err := f.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(alvo)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

// encontrarRaizApp localiza a pasta que contém o Dockerfile do pacote.
func encontrarRaizApp(dest string) (string, error) {
	var raiz string
	err := filepath.WalkDir(dest, func(caminho string, d os.DirEntry, err error) error {
		if err != nil || raiz != "" || d.IsDir() || d.Name() != "Dockerfile" {
			return err
		}
		raiz = filepath.Dir(caminho)
		return nil
	})
	if err != nil || raiz == "" {
		return "", ErrInvalidDeployInput
	}
	return raiz, nil
}
