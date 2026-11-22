#!/bin/sh
set -eu

DRY_RUN=0
BASE_URL=""
INSTALL_DIR=""
VERSAO="latest"
HOME_DIR="${HOME:-}"

while [ "$#" -gt 0 ]; do
	case "$1" in
		--dry-run)
			DRY_RUN=1
			shift
			;;
		--base-url)
			BASE_URL="${2:-}"
			shift 2
			;;
		--base-url=*)
			BASE_URL="${1#--base-url=}"
			shift
			;;
		--install-dir)
			INSTALL_DIR="${2:-}"
			shift 2
			;;
		--install-dir=*)
			INSTALL_DIR="${1#--install-dir=}"
			shift
			;;
		--versao)
			VERSAO="${2:-}"
			shift 2
			;;
		--versao=*)
			VERSAO="${1#--versao=}"
			shift
			;;
		--home)
			HOME_DIR="${2:-}"
			shift 2
			;;
		--home=*)
			HOME_DIR="${1#--home=}"
			shift
			;;
		*)
			printf 'erro: argumento desconhecido: %s\n' "$1" >&2
			exit 2
			;;
	esac
done

detectar_so() {
	case "$(uname -s)" in
		Linux*) printf 'linux' ;;
		Darwin*) printf 'darwin' ;;
		MINGW*|MSYS*|CYGWIN*) printf 'windows' ;;
		*)
			printf 'erro: sistema operacional nao suportado: %s\n' "$(uname -s)" >&2
			exit 1
			;;
	esac
}

detectar_arch() {
	case "$(uname -m)" in
		x86_64|amd64) printf 'amd64' ;;
		aarch64|arm64) printf 'arm64' ;;
		*)
			printf 'erro: arquitetura nao suportada: %s\n' "$(uname -m)" >&2
			exit 1
			;;
	esac
}

validar_url() {
	case "$1" in
		https://*|http://127.0.0.1*|http://localhost*|http://\[::1\]*)
			return 0
			;;
		*)
			printf 'erro: download exige HTTPS (ou loopback para testes): %s\n' "$1" >&2
			exit 1
			;;
	esac
}

baixar_arquivo() {
	url="$1"
	dest="$2"
	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$url" -o "$dest"
	elif command -v wget >/dev/null 2>&1; then
		wget -qO "$dest" "$url"
	else
		printf 'erro: curl ou wget nao encontrado\n' >&2
		exit 1
	fi
}

calcular_sha256() {
	arq="$1"
	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$arq" | awk '{print tolower($1)}'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$arq" | awk '{print tolower($1)}'
	elif command -v openssl >/dev/null 2>&1; then
		openssl dgst -sha256 "$arq" | awk '{print tolower($NF)}'
	else
		printf 'erro: utilitario SHA-256 nao encontrado\n' >&2
		exit 1
	fi
}

SO="$(detectar_so)"
ARCH="$(detectar_arch)"

if [ -z "$BASE_URL" ]; then
	if [ "$VERSAO" = "latest" ]; then
		BASE_URL="https://github.com/heverton-dev/LegoLoop/releases/latest/download"
	else
		BASE_URL="https://github.com/heverton-dev/LegoLoop/releases/download/${VERSAO}"
	fi
fi
BASE_URL="${BASE_URL%/}"
validar_url "$BASE_URL"

if [ -z "$INSTALL_DIR" ]; then
	INSTALL_DIR="${HOME_DIR}/.local/bin"
fi

if [ "$SO" = "windows" ]; then
	EXT="zip"
	BIN_NOME="legoloop.exe"
else
	EXT="tar.gz"
	BIN_NOME="legoloop"
fi

ARQUIVO="legoloop_${VERSAO}_${SO}_${ARCH}.${EXT}"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

SUMS_PATH="${TMP_DIR}/SHA256SUMS"
PKG_PATH="${TMP_DIR}/${ARQUIVO}"

baixar_arquivo "${BASE_URL}/SHA256SUMS" "$SUMS_PATH"
baixar_arquivo "${BASE_URL}/${ARQUIVO}" "$PKG_PATH"

ESPERADO="$(awk -v alvo="$ARQUIVO" '$2 == alvo {print tolower($1)}' "$SUMS_PATH" | head -n 1)"
if [ -z "$ESPERADO" ]; then
	printf 'erro: checksum SHA-256 ausente para %s em SHA256SUMS\n' "$ARQUIVO" >&2
	exit 1
fi

ATUAL="$(calcular_sha256 "$PKG_PATH")"
if [ "$ATUAL" != "$ESPERADO" ]; then
	printf 'erro: SHA-256 divergente para %s (esperado %s, obtido %s)\n' "$ARQUIVO" "$ESPERADO" "$ATUAL" >&2
	exit 1
fi

if [ "$DRY_RUN" -eq 1 ]; then
	printf 'DRY-RUN: SHA-256 conferido (%s); instalaria %s em %s e executaria legoloop instalar\n' "$ATUAL" "$BIN_NOME" "$INSTALL_DIR"
	exit 0
fi

mkdir -p "$INSTALL_DIR"
if [ "$EXT" = "zip" ]; then
	unzip -oq "$PKG_PATH" -d "$INSTALL_DIR"
else
	tar -xzf "$PKG_PATH" -C "$INSTALL_DIR"
fi

BIN_DEST="${INSTALL_DIR}/${BIN_NOME}"
chmod +x "$BIN_DEST" 2>/dev/null || true

if [ -n "$HOME_DIR" ]; then
	"$BIN_DEST" instalar --apply --home "$HOME_DIR"
else
	"$BIN_DEST" instalar --apply
fi
