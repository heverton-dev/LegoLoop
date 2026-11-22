//go:build darwin

package orchestrix

import (
	"os/exec"
	"strconv"
	"strings"
)

func memoriaLivreSO() (uint64, error) {
	out, err := exec.Command("vm_stat").Output()
	if err == nil {
		if bytes := parseVMStatLivre(string(out)); bytes > 0 {
			return bytes, nil
		}
	}
	return 0, ErrMemoriaIndisponivel
}

func parseVMStatLivre(saida string) uint64 {
	var pageSize uint64 = 4096
	var paginas uint64
	for _, linha := range strings.Split(saida, "\n") {
		if idx := strings.Index(linha, "page size of "); idx >= 0 {
			rest := linha[idx+len("page size of "):]
			if sp := strings.IndexByte(rest, ' '); sp > 0 {
				if ps, err := strconv.ParseUint(rest[:sp], 10, 64); err == nil && ps > 0 {
					pageSize = ps
				}
			}
			continue
		}
		if strings.HasPrefix(linha, "Pages free:") ||
			strings.HasPrefix(linha, "Pages inactive:") ||
			strings.HasPrefix(linha, "Pages speculative:") {
			partes := strings.Split(linha, ":")
			if len(partes) == 2 {
				num := strings.TrimRight(strings.TrimSpace(partes[1]), ".")
				if n, err := strconv.ParseUint(num, 10, 64); err == nil {
					paginas += n
				}
			}
		}
	}
	return paginas * pageSize
}
