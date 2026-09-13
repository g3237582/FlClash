<div>

[**简体中文**](README_zh_CN.md)

</div>

## FlClash

[![Downloads](https://img.shields.io/github/downloads/chen08209/FlClash/total?style=flat-square&logo=github)](https://github.com/chen08209/FlClash/releases/)[![Last Version](https://img.shields.io/github/release/chen08209/FlClash/all.svg?style=flat-square)](https://github.com/chen08209/FlClash/releases/)[![License](https://img.shields.io/github/license/chen08209/FlClash?style=flat-square)](LICENSE)

[![Channel](https://img.shields.io/badge/Telegram-Channel-blue?style=flat-square&logo=telegram)](https://t.me/FlClash)

A multi-platform proxy client based on ClashMeta, simple and easy to use, open-source and ad-free.

on Desktop:
<p style="text-align: center;">
    <img alt="desktop" src="snapshots/desktop.gif">
</p>

on Mobile:
<p style="text-align: center;">
    <img alt="mobile" src="snapshots/mobile.gif">
</p>

## Features

✈️ Multi-platform: Android, Windows, macOS and Linux

💻 Adaptive multiple screen sizes, Multiple color themes available

💡 Based on Material You Design, [Surfboard](https://github.com/getsurfboard/surfboard)-like UI

☁️ Supports data sync via WebDAV

✨ Support subscription link, Dark mode

## Use

### Linux

⚠️ Make sure to install the following dependencies before using them

   ```bash
    sudo apt-get install libayatana-appindicator3-dev
   ```

### Android

Support the following actions

   ```bash
    com.follow.clash.action.START
    
    com.follow.clash.action.STOP
    
    com.follow.clash.action.TOGGLE
   ```

### Android EasyTier (FFI) + Tailscale

This fork can start a native EasyTier outbound (`type: easytier`) next to
the existing Tailscale outbound. Full steps: [android/easytier/README.md](android/easytier/README.md).

1. Install an arm64 APK that includes `libeasytier_ffi.so` (see Build below,
   then `tool/bundle_easytier_ffi.sh` before `dart setup.dart android`).
2. Point `ffi-library` at `./libeasytier_ffi.so`. Core copies the packaged
   library into the app files directory (`home-dir`) on init.
3. Use the overwrite script
   [examples/easytier/easytier-tailscale-overwrite.js](examples/easytier/easytier-tailscale-overwrite.js)
   (placeholders only). YAML shape:
   [examples/easytier/easytier-tailscale.yaml](examples/easytier/easytier-tailscale.yaml).
4. Verify `10.77.0.0/24` via EasyTier, `100.64.0.0/10` via Tailscale, and
   RFC1918 LAN as `DIRECT`. WireGuard Portal is fallback only.

The Clash.Meta submodule stays on the FlClash-patched pin. EasyTier FFI
from [g3237582/mihomo `cursor/easytier-outbound-d1fd`](https://github.com/g3237582/mihomo/tree/cursor/easytier-outbound-d1fd)
([PR #1](https://github.com/g3237582/mihomo/pull/1)) is applied from
`core/easytier_overlay/` (`tool/apply_easytier_overlay.sh`, also run by
the Core build hook). A raw swap to that Alpha branch would drop FlClash
`AllProxies` APIs.


## Download

<a href="https://chen08209.github.io/FlClash-fdroid-repo/repo?fingerprint=789D6D32668712EF7672F9E58DEEB15FBD6DCEEC5AE7A4371EA72F2AAE8A12FD"><img alt="Get it on F-Droid" src="snapshots/get-it-on-fdroid.svg" width="200px"/></a> <a href="https://github.com/chen08209/FlClash/releases"><img alt="Get it on GitHub" src="snapshots/get-it-on-github.svg" width="200px"/></a>

### Homebrew

```bash
brew tap chen08209/tap
brew install --cask flclash
```

## Build

1. Update submodules and apply the EasyTier FFI overlay
   ```bash
   git submodule update --init --recursive
   tool/apply_easytier_overlay.sh
   ```

2. Install `Flutter` and `Golang` environment

3. Build Application

    - android

        1. Install `Android SDK`, `Android NDK`

        2. Set `ANDROID_NDK` environment variable

        3. Build `libeasytier_ffi.so` (arm64) into
           `android/easytier/jniLibs/arm64-v8a/`, then run the build script

           ```bash
           tool/bundle_easytier_ffi.sh
           dart setup.dart android
           ```

           The setup hook also applies `core/easytier_overlay/`. See
           [android/easytier/README.md](android/easytier/README.md).

    - windows

        1. Requires a Windows client

        2. Install `GCC`, `Inno Setup`

        3. Run build script

           ```bash
           dart setup.dart windows
           ```

    - linux

        1. Requires a Linux client

        2. Dependencies are auto-installed by setup script, or manually:
           ```bash
           sudo apt-get install -y libayatana-appindicator3-dev
           ```

        3. Run build script

           ```bash
           dart setup.dart linux
           ```

    - macOS

        1. Requires a macOS client

        2. Run build script

           ```bash
           dart setup.dart macos
           ```

## Star

The easiest way to support developers is to click on the star (⭐) at the top of the page.

<p style="text-align: center;">
    <a href="https://api.star-history.com/svg?repos=chen08209/FlClash&Date">
        <img alt="start" width=50% src="https://api.star-history.com/svg?repos=chen08209/FlClash&Date"/>
    </a>
</p>
