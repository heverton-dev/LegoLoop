package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var dataNormalizada = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// Empacotar cria o arquivo .zip (Windows) ou .tar.gz (Linux/Darwin) com timestamps fixos.
func Empacotar(so, arch, versao, binPath, outDir string) (string, error) {
	base := fmt.Sprintf("legoloop_%s_%s_%s", versao, so, arch)
	nomeInterno := NomeBinario(so)
	if so == "windows" {
		dest := filepath.Join(outDir, base+".zip")
		return dest, gerarZip(dest, nomeInterno, binPath)
	}
	dest := filepath.Join(outDir, base+".tar.gz")
	return dest, gerarTarGz(dest, nomeInterno, binPath)
}

func gerarZip(destPath, nomeInterno, srcPath string) error {
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	hdr := &zip.FileHeader{
		Name:     nomeInterno,
		Method:   zip.Deflate,
		Modified: dataNormalizada,
	}
	hdr.SetMode(0755)
	w, err := zw.CreateHeader(hdr)
	if err != nil {
		return err
	}
	if err := copiarArquivoParaWriter(srcPath, w); err != nil {
		return err
	}
	return zw.Close()
}

func gerarTarGz(destPath, nomeInterno, srcPath string) error {
	info, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	gz.ModTime = dataNormalizada
	gz.Name = ""
	tw := tar.NewWriter(gz)
	hdr := &tar.Header{
		Name:    nomeInterno,
		Mode:    0755,
		Size:    info.Size(),
		ModTime: dataNormalizada,
		Format:  tar.FormatPAX,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if err := copiarArquivoParaWriter(srcPath, tw); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}
	return gz.Close()
}

func copiarArquivoParaWriter(srcPath string, w io.Writer) error {
	in, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer in.Close()
	_, err = io.Copy(w, in)
	return err
}

// GerarSHA256SUMS grava o arquivo SHA256SUMS no formato "<sha256>  <arquivo>".
func GerarSHA256SUMS(outDir string, arquivos []string) (string, error) {
	ordenados := append([]string(nil), arquivos...)
	sort.Strings(ordenados)
	var b strings.Builder
	for _, nome := range ordenados {
		hash, err := HashArquivo(filepath.Join(outDir, nome))
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&b, "%s  %s\n", hash, nome)
	}
	sumsPath := filepath.Join(outDir, "SHA256SUMS")
	if err := os.WriteFile(sumsPath, []byte(b.String()), 0644); err != nil {
		return "", err
	}
	return sumsPath, nil
}

// VerificarSHA256SUMS confere cada linha do arquivo SHA256SUMS contra os arquivos em outDir.
func VerificarSHA256SUMS(outDir, sumsPath string) error {
	dados, err := os.ReadFile(sumsPath)
	if err != nil {
		return err
	}
	linhas := strings.Split(strings.TrimSpace(string(dados)), "\n")
	for _, linha := range linhas {
		if strings.TrimSpace(linha) == "" {
			continue
		}
		partes := strings.SplitN(linha, "  ", 2)
		if len(partes) != 2 {
			return fmt.Errorf("linha invalida em SHA256SUMS: %q", linha)
		}
		atual, err := HashArquivo(filepath.Join(outDir, strings.TrimSpace(partes[1])))
		if err != nil {
			return err
		}
		if atual != partes[0] {
			return fmt.Errorf("checksum divergente para %s: esperado %s, obtido %s", partes[1], partes[0], atual)
		}
	}
	return nil
}
