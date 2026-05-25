# 🏗️ Building

## Prerequisites

- Go toolchain  
  [go.dev/dl](https://go.dev/dl/)
- Just command runner  
  [just.systems](https://just.systems/man/en/)
- CMake and Ninja  
  SDL3 is built from source via CMake w/ Ninja
- Node.js  
  Required to build the UI (`UI/`) and the CEF payloads (`cefpayloads/`) that
  interact with Steam's CEF remote debugging interface

!!! info "SDL3"
    SDL3 is compiled from source in `deps/SDL` via CMake + Ninja (`just build-sdl`).

### 🪟 Windows

- MSVC toolchain
- Ensure `cmake` and `ninja` are in PATH

### 🐧 Linux

#### 🏹 Arch Linux

```bash
sudo pacman -S ninja cmake pkg-config gtk3 webkit2gtk-4.1 xdotool libxss
```

#### 🟠 Ubuntu

```bash
sudo apt-get install ninja-build cmake pkg-config libgtk-3-dev \
  libsoup-3.0-dev libjavascriptcoregtk-4.1-dev libwebkit2gtk-4.1-dev \
  libxdo-dev libxss-dev
```

## Build

Use the `justfile`.

**1. Checkout**:

```bash
git clone git@github.com:Alia5/SISR.git
cd SISR
git submodule update --init --recursive
```

**2. Build**:

```bash
just build release
```

!!! info "VIIPER"
    CI builds bundle VIIPER (windows only)  
    Please install/run a VIIPER server before running SISR.    
    [VIIPER Docs](https://alia5.github.io/VIIPER/)  
