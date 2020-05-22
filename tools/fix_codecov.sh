#!/usr/bin/env bash

set -euo pipefail

rootDir="bazel-out"

echo "mode: set"

while IFS= read -r -d '' i
do
    # >&2 echo "${i} -> ${i//.dat/.txt}"
    >&2 echo "parsing: ${i}"
    # cp "${i}" "${i//.dat/.txt}"
    tail -n +2 "${i}"
done <   <(find "${rootDir}"/ -name coverage.dat -print0)
