//go:build linux

package orchestrix

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func memoriaLivreSO() (uint64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		campos := strings.Fields(sc.Text())
		if len(campos) >= 2 && campos[0] == "MemAvailable:" {
			kb, err := strconv.ParseUint(campos[1], 10, 64)
			return kb << 10, err
		}
	}
	return 0, ErrMemoriaIndisponivel
}
