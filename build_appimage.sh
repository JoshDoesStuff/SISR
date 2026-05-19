#!/usr/bin/env bash
set -e

ARCH=$(uname -m)
RUST_TARGET="${ARCH}-unknown-linux-gnu"

rm -rf AppDir
mkdir -p AppDir/usr/lib64
export LD_LIBRARY_PATH="target/release:$LD_LIBRARY_PATH"
export LINUXDEPLOY_EXCLUDED_LIBRARIES=""

NO_STRIP=1 APPIMAGE_EXTRACT_AND_RUN=1 linuxdeploy --appimage-extract-and-run --appdir AppDir -l target/release/libSDL3.so.0 -e target/release/SISR -d sisr.desktop -i docs/SISR.svg --output appimage
mkdir -p dist/appimage
mv SISR-${ARCH}.AppImage dist/appimage/
