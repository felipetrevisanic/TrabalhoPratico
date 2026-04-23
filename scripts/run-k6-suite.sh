#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RESULT_DIR="${ROOT_DIR}/result/k6"
PROFILE_INPUT="${1:-single}"
TARGET_INPUT="${2:-all}"

case "${PROFILE_INPUT}" in
  unico) PROFILE_INPUT="single" ;;
  pequeno) PROFILE_INPUT="small" ;;
  medio) PROFILE_INPUT="medium" ;;
  grande) PROFILE_INPUT="large" ;;
esac

mkdir -p "${RESULT_DIR}"

TARGETS=(
  "csharp-rest"
  "csharp-graphql"
  "csharp-grpc"
  "go-rest"
  "go-graphql"
  "go-grpc"
  "rust-rest"
  "rust-graphql"
  "rust-grpc"
)

if [[ "${TARGET_INPUT}" != "all" ]]; then
  TARGETS=("${TARGET_INPUT}")
fi

for target in "${TARGETS[@]}"; do
  timestamp="$(date +%Y%m%d-%H%M%S)"
  summary_file="${RESULT_DIR}/${timestamp}-${PROFILE_INPUT}-${target}-summary.json"
  html_file="${RESULT_DIR}/${timestamp}-${PROFILE_INPUT}-${target}-report.html"

  echo "Executando perfil ${PROFILE_INPUT} para ${target}"
  k6 run \
    -e TARGET_ID="${target}" \
    -e DATA_PROFILE="${PROFILE_INPUT}" \
    -e HTML_REPORT_FILE="${html_file}" \
    --summary-export "${summary_file}" \
    "${ROOT_DIR}/k6/api-load.js"
done

echo "Suite k6 concluida."
