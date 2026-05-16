<img src="SISR.svg" align="right" width="128"/>
<br />

# SISR ✂️

**S**team **I**nput **S**ystem **R**edirector

SISR (pronounced "scissor") redirects Steam Input configurations to the system level (localhost or network).  

SISR takes controllers it receives from Steam (via Steam Input)
and forwards them as emulated but more compatible controllers (indistinguishable from real hardware) to the OS.

It can be used to circumvent issues with games and applications that
do not support Steam Input or otherwise pose challenges, like (but not limited to):

- Games with aggressive anti-cheat systems
- Emulators
- Windows Store games/apps
- Games with broken Steam Input support

All while still providing you the full feature set of SteamInput!

SISR can also be used to "tunnel"/forward Steam Input configurations over the network to other machines, including Keyboard/Mouse.  
This makes it possible to use devices like a Steam Deck as a dedicated controller without the need to stream the entire game.

The emulated controllers (and Keyboard/Mouse) are indistinguishable from real hardware and show up at system level.  
SISR achieves this by utilizing [VIIPER](https://github.com/Alia5/VIIPER) (requires **USBIP**).  
Unlike its predecessor [GlosSI](https://github.com/Alia5/GlosSI), it does not use the unmaintained [ViGEm](https://github.com/ViGEm/ViGEmBus) driver.


!!! warning
    **Highly experimental work in progress.** Everything is subject to change and may or may not work.  
    Expect bugs, crashes, and missing features.

!!! danger
    You **do** get the full functionality of SteamInput, _without having to launch your games from Steam_  
    For this Steam must still be running in the background, though.   
    **Please read the [introduction post](https://alia5.github.io/SISR/main/getting-started/introduction/) before you get started.**  

## ✨🛣️ Features / Roadmap

- Full SteamInput featureset while emulating compatible controllers (indistinguishable from real hardware)
  - Xbox360 _or_
  - DualShock 4
- **Non Steam Mode*
-  Xbox 360 controller emulation
- Multi-platform support (Windows, Linux)
   Multiple operation modes
    - Standalone background service (To be improved)
    - Steam overlay window mode
-  PS4 controller emulation
- Networked operation across computers
  - Use devices like a SteamDeck as dedicated controller without streaming the whole game/display
- ~~🚧 Xbox One controller emulation~~
- ~~🚧 Generic controller emulation~~
- 🚧 Gyro Passthrough
- 🚧 Bundling multiple devices into a single controller
- 🚧 Automatic HidHide integration

## 🚀 Getting started

- [Installation](getting-started/installation.md)

## ⚙️ Configuration

- [Configuration](config/config.md)
- [CLI Reference](config/cli.md)

## 🆘 Help

- [Guides](guides/overview.md)
- [Troubleshooting](misc/troubleshooting.md)
- [FAQ](misc/faq.md)

## 🛠️ Development

- [Building](dev/building.md)

## 🔗 Links

- [📥 Downloads](downloads/index.md)
- [GitHub Repository](https://github.com/Alia5/SISR)
- [SISR Releases](https://github.com/Alia5/SISR/releases)
- [VIIPER Docs](https://alia5.github.io/VIIPER/)
- [USBIP-Win2 (Windows USBIP)](https://github.com/vadimgrn/usbip-win2)
- [Changelog](changelog/)
