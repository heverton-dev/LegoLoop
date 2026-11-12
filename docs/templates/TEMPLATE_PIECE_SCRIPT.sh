#!/bin/sh
# {{NOME_DO_SCRIPT}} - peça Script do LegoLoop (template V2.2.0).
# Uso: script.sh --target-file ARQUIVO [--project-root RAIZ] [--dry-run] [--json]
# Códigos de saída: 0 sucesso, 1 falha de execução, 2 uso ou entrada inválida.

set -eu

target=""
root="$(pwd)"
dry_run=0
json=0

usage_error() {
    printf 'ERRO DE USO: %s\n' "$1" >&2
    exit 2
}

exec_error() {
    printf 'FALHA DE EXECUÇÃO: %s\n' "$1" >&2
    exit 1
}

info() {
    if [ "$json" -eq 0 ]; then
        printf '%s\n' "$1"
    fi
}

emit_json() {
    if [ "$json" -eq 1 ]; then
        printf '{"status":"%s","target":"%s","changed":%s,"exit_code":0}\n' "$1" "$target" "$2"
    fi
}

while [ "$#" -gt 0 ]; do
    case "$1" in
        --target-file)
            [ "$#" -ge 2 ] || usage_error "--target-file exige um valor"
            target="$2"
            shift 2
            ;;
        --project-root)
            [ "$#" -ge 2 ] || usage_error "--project-root exige um valor"
            root="$2"
            shift 2
            ;;
        --dry-run)
            dry_run=1
            shift
            ;;
        --json)
            json=1
            shift
            ;;
        *)
            usage_error "argumento desconhecido: $1"
            ;;
    esac
done

[ -n "$target" ] || usage_error "--target-file é obrigatório"
[ -d "$root" ] || usage_error "raiz inexistente: $root"

root_full="$(cd "$root" && pwd -P)"
case "$target" in
    /*) candidate="$target" ;;
    *) candidate="$root_full/$target" ;;
esac
case "$candidate" in
    */../* | */..) usage_error "alvo fora da raiz do projeto: $target" ;;
esac
case "$candidate" in
    "$root_full"/*) ;;
    *) usage_error "alvo fora da raiz do projeto: $target" ;;
esac
[ -f "$candidate" ] || usage_error "arquivo alvo não encontrado: $candidate"
target="$candidate"

normalize() {
    sed -e 's/\r$//' -e 's/[[:space:]]*$//' "$1" | sed -e :a -e '/^\n*$/{$d;N;ba' -e '}'
}

tmp="$target.legoloop-tmp"
trap 'rm -f "$tmp"' EXIT

normalize "$target" >"$tmp" || exec_error "falha ao normalizar $target"

if cmp -s "$tmp" "$target"; then
    info "Nada a fazer: o arquivo já está normalizado ($target)."
    emit_json "noop" "false"
    exit 0
fi

if [ "$dry_run" -eq 1 ]; then
    info "[DRY-RUN] O arquivo seria normalizado: $target. Nenhuma alteração gravada."
    emit_json "dry_run" "true"
    exit 0
fi

[ -w "$target" ] || exec_error "arquivo somente leitura, nada foi gravado: $target"
mv -f "$tmp" "$target" || exec_error "falha ao substituir $target"

info "Sucesso: arquivo normalizado ($target)."
emit_json "applied" "true"
exit 0
