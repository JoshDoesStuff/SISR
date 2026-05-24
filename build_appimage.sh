#!/usr/bin/env bash
set -euo pipefail

ARCH=$(uname -m)

shopt -s nullglob
sdl_libs=(dist/libSDL3.so*)
if [ ${#sdl_libs[@]} -eq 0 ]; then
	echo "Error: No SDL3 shared library found in dist/."
	ls -la dist
	exit 1
fi
SDL_LIB="${sdl_libs[0]}"

rm -rf AppDir
mkdir -p AppDir/usr/lib64
export LD_LIBRARY_PATH="dist${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
export LINUXDEPLOY_EXCLUDED_LIBRARIES=""

NO_STRIP=1 APPIMAGE_EXTRACT_AND_RUN=1 linuxdeploy --appimage-extract-and-run --appdir AppDir -l "$SDL_LIB" -e dist/SISR -d sisr.desktop -i docs/SISR.svg --output appimage
mkdir -p dist/appimage
mv SISR-${ARCH}.AppImage dist/appimage/
