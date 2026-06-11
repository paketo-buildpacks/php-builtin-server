#!/usr/bin/env bash

set -eu
set -o pipefail

readonly PROGDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly BUILDPACKDIR="$(cd "${PROGDIR}/.." && pwd)"

function main() {
  rm -rf "${BUILDPACKDIR}/linux"

  build::target "linux/amd64"
  build::target "linux/arm64"
}

function build::target() {
  local target="${1}"
  local platform arch
  platform="$(echo "${target}" | cut -d '/' -f1)"
  arch="$(echo "${target}" | cut -d '/' -f2)"

  mkdir -p "${BUILDPACKDIR}/${platform}/${arch}/bin"

  GOOS="${platform}" \
  GOARCH="${arch}" \
  CGO_ENABLED=0 \
    go build \
      -ldflags="-s -w" \
      -o "${BUILDPACKDIR}/${platform}/${arch}/bin/run" \
      "${BUILDPACKDIR}/run"

  ln -sf "run" "${BUILDPACKDIR}/${platform}/${arch}/bin/detect"
  ln -sf "run" "${BUILDPACKDIR}/${platform}/${arch}/bin/build"
}

main