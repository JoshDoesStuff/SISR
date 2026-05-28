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

NO_STRIP=1 APPIMAGE_EXTRACT_AND_RUN=1 linuxdeploy --appimage-extract-and-run --appdir AppDir -l "$SDL_LIB" -e dist/SISR -d sisr.desktop -i docs/SISR.svg

MULTIARCH=$(gcc -print-multiarch || dpkg-architecture -qDEB_HOST_MULTIARCH || echo "${ARCH}-linux-gnu")
WEBKIT_EXEC_SRC="/usr/lib/${MULTIARCH}/webkitgtk-6.0"
if [ -d "$WEBKIT_EXEC_SRC" ]; then
    mkdir -p AppDir/usr/lib/webkitgtk-6.0
    cp -a "$WEBKIT_EXEC_SRC"/. AppDir/usr/lib/webkitgtk-6.0/

    while IFS= read -r -d '' helper; do
        ldd "$helper" | awk '/=> \//{print $3}' | while IFS= read -r dep; do
            depname=$(basename "$dep")
            if [ ! -f "AppDir/usr/lib/$depname" ] && [ ! -f "AppDir/usr/lib64/$depname" ]; then
                cp -L "$dep" "AppDir/usr/lib/"
            fi
        done
    done < <(find AppDir/usr/lib/webkitgtk-6.0 -type f -executable -print0)

    sed -i '/^exec /i export WEBKIT_EXEC_PATH="$APPDIR/usr/lib/webkitgtk-6.0"' AppDir/AppRun
fi

curl -fsSL "https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-${ARCH}.AppImage" -o appimagetool
chmod +x appimagetool
ARCH="${ARCH}" APPIMAGE_EXTRACT_AND_RUN=1 ./appimagetool --appimage-extract-and-run AppDir "SISR-${ARCH}.AppImage"

mkdir -p dist/appimage
mv SISR-${ARCH}.AppImage dist/appimage/
