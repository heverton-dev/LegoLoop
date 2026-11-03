//go:build !windows

package fabricar

import "strings"

// isUntrustedRoot trata como não confiável o que é montado fora da raiz: mídia removível e volumes externos.
func isUntrustedRoot(root string) bool {
	if root == "/" {
		return false
	}
	for _, prefix := range []string{"/media/", "/mnt/", "/run/media/", "/Volumes/"} {
		if strings.HasPrefix(root, prefix) {
			return true
		}
	}
	return false
}
