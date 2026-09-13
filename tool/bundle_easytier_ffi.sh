#!/usr/bin/env bash
# Build Android arm64 libeasytier_ffi.so and place it in the FlClash jniLibs drop-in.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
dest="$root/android/easytier/jniLibs/arm64-v8a/libeasytier_ffi.so"
work="${EASYTIER_SRC:-}"
api="${ANDROID_API_LEVEL:-21}"

if [[ -z "${ANDROID_NDK:-}" && -z "${ANDROID_NDK_HOME:-}" ]]; then
  echo "ANDROID_NDK or ANDROID_NDK_HOME is required to cross-compile easytier-ffi." >&2
  echo "Install the Android NDK, then re-run this script." >&2
  exit 1
fi

if [[ -z "$work" ]]; then
  work="$(mktemp -d "${TMPDIR:-/tmp}/easytier-ffi.XXXXXX")"
  trap 'rm -rf "$work"' EXIT
  git clone --depth 1 https://github.com/EasyTier/EasyTier.git "$work"
fi

if ! command -v cargo >/dev/null; then
  echo "cargo is required." >&2
  exit 1
fi

if ! rustup target list --installed | grep -q '^aarch64-linux-android$'; then
  rustup target add aarch64-linux-android
fi

(
  cd "$work"
  cargo ndk -t arm64-v8a -p "$api" build --release -p easytier-ffi --features 'c-abi,ffi-dataplane' \
    || cargo build --release -p easytier-ffi --features 'c-abi,ffi-dataplane' --target aarch64-linux-android
)

so=""
for candidate in \
  "$work/target/aarch64-linux-android/release/libeasytier_ffi.so" \
  "$work/target/release/libeasytier_ffi.so"; do
  if [[ -f "$candidate" ]]; then
    so="$candidate"
    break
  fi
done

if [[ -z "$so" ]]; then
  echo "libeasytier_ffi.so was not produced. Check cargo/NDK output above." >&2
  exit 1
fi

mkdir -p "$(dirname "$dest")"
cp -f "$so" "$dest"
echo "Wrote $dest"
