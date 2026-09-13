#!/usr/bin/env bash
# Copy EasyTier FFI sources onto the Clash.Meta submodule working tree.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
overlay="$root/core/easytier_overlay"
dest="$root/core/Clash.Meta"

if [[ ! -d "$overlay" ]]; then
  echo "missing $overlay" >&2
  exit 1
fi
if [[ ! -d "$dest" ]]; then
  echo "missing $dest (run git submodule update --init --recursive)" >&2
  exit 1
fi

find "$overlay" -type f ! -name SOURCE.txt -print0 | while IFS= read -r -d '' file; do
  rel="${file#"$overlay"/}"
  mkdir -p "$dest/$(dirname "$rel")"
  cp -a "$file" "$dest/$rel"
done

echo "Applied EasyTier overlay to $dest"
