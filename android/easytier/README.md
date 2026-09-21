# Android EasyTier outbound (FFI, no TUN)

FlClash on Android can join an EasyTier mesh as a `type: easytier` outbound
and keep the existing Tailscale outbound. The Core home directory is the
single place overwrite scripts should point `ffi-library` at.

## What is wired

The embedded mihomo tree is `core/Clash.Meta`, a git submodule of
`chen08209/Clash.Meta` (`FlClash`, pin `70f0570`). The FlClash wrapper
(`core/*.go`) needs that tree's `AllProxies` / `ProvidersSnapshot` patches.

The EasyTier-capable branch
https://github.com/g3237582/mihomo/tree/cursor/easytier-outbound-d1fd
([PR #1](https://github.com/g3237582/mihomo/pull/1), smoke-tested on
glq-mesh) is MetaCubeX Alpha and does not include those FlClash APIs. This
fork therefore keeps the Clash.Meta pin and applies
`core/easytier_overlay/` (FFI outbound from that branch) during Core build
and `go test`:

```bash
git submodule update --init --recursive
tool/apply_easytier_overlay.sh
```

The setup hook (`plugins/setup/setup_hooks`) applies the same overlay
before it compiles `libclash.so` / `FlClashCore`. Tailscale parsing is
unchanged. This environment cannot push a FlClash+EasyTier branch to
`g3237582/mihomo`, so the overlay lives in this repository.

Primary path: `ffi-library` + `instance-name` + exactly one of `config` /
`config-file`, `udp: true`. Mihomo forces `[flags] no_tun = true` and
`bind_device = false`. WireGuard Portal is fallback only.

## Install FlClash

1. Update the submodule, apply the overlay, then install Flutter and Go.

   ```bash
   git submodule update --init --recursive
   tool/apply_easytier_overlay.sh
   ```

2. Build the Android arm64 FFI library (needs Android NDK + Rust):

   ```bash
   export ANDROID_NDK=/path/to/ndk
   tool/bundle_easytier_ffi.sh
   ```

   Artifact: `android/easytier/jniLibs/arm64-v8a/libeasytier_ffi.so`

   AGP packages that directory through the `:core` extra `jniLibs` source
   set. Do not commit a real `.so` unless you built it yourself.

3. Build the APK (needs Android SDK + NDK; Flutter runs the Go Core hook):

   ```bash
   dart setup.dart android
   ```

   Or, after `flutter pub get`:

   ```bash
   flutter build apk --release --target-platform android-arm64
   ```

   This environment has no Android SDK/NDK, so an APK is not produced here.

## Where `ffi-library` points

| Stage | Path |
| --- | --- |
| Drop-in (source) | `android/easytier/jniLibs/arm64-v8a/libeasytier_ffi.so` |
| Inside the APK | `lib/arm64-v8a/libeasytier_ffi.so` |
| Extracted native dir | `{applicationInfo.nativeLibraryDir}/libeasytier_ffi.so` |
| After Core init | `{home-dir}/libeasytier_ffi.so` |

`home-dir` is `Context.getFilesDir()` (Flutter `appPath.homeDirPath`).
Overwrite scripts should use:

```yaml
ffi-library: ./libeasytier_ffi.so
```

On Android, mihomo's `IsSafePath` already allows any path. Staging into
home keeps the same relative path as desktop.

## Profile overwrite

Copy `examples/easytier/easytier-tailscale-overwrite.js` into a FlClash
overwrite script (`Profiles → overwrite → script`). Replace every
`CHANGE_ME_*` placeholder. YAML shape:
`examples/easytier/easytier-tailscale.yaml`.

Rules:

- `10.77.0.0/24` → `easytier`
- `100.64.0.0/10` → `tailscale`
- RFC1918 + loopback → `DIRECT`

## Verify on a mesh (glq-mesh style)

Do not treat "instance started" as connectivity.

1. Install the arm64 APK that contains `libeasytier_ffi.so`.
2. Confirm the staged library exists:

   ```text
   {filesDir}/libeasytier_ffi.so
   ```

   `adb shell run-as com.follow.clash ls files/libeasytier_ffi.so`
   (use `.dev` suffix for debug builds).

3. Fill placeholders with your mesh `network_name`, `network_secret`, and
   peer URI. Keep those values off git.

4. Start Core. In logs, look for
   `[EasyTier](...) FFI instance ... started (no_tun)`.

5. EasyTier overlay: from the phone, reach a listener on
   `10.77.0.0/24` (example: `http://10.77.0.2:8080/` through the mixed
   port). Success is a completed TCP response.

6. Tailscale: the node must reach **Online** in the tailnet admin
   (LastSeen updating), not stay Offline. Then reach a `100.64.0.0/10`
   address the same way (`ping`/`curl` `http://100.x.x.x` through the
   mixed port). Opening that URL must not flash-exit. If TailscaleIPs
   are still empty after the backend wait, Core should log
   `tailscale: no valid IPv4 address (backend not ready)` (or a recovered
   `tailscale panic: ...`) instead of aborting the process.

7. LAN: `10/8`, `172.16/12`, `192.168/16` must stay `DIRECT`.

## Residual risks

- FFI data-plane ABI v3 is IPv4-only.
- `DataPlaneSocketAddr` by-value submits (`tcp_connect_submit`,
  `udp_send_submit`) use the SysV stack path on amd64 and an AAPCS64
  hidden pointer (`*SocketAddr`) on arm64 and other arches. The
  linux/amd64 stub ABI is executed in this tree; Android arm64 now
  uses the pointer path and still needs a device/mesh check.
- `interface-name` / `routing-mark` / `dialer-proxy` do not apply to
  sockets created inside `libeasytier_ffi`.
- One native session per `ffi-library` + `instance-name`.
- Desktop/iOS FFI packaging is out of scope.
- `FORCE_ANET` prefers netlink `interfaceTable`, then falls back to
  `net.Interfaces()` when netlink errors or returns an empty list. HarmonyOS
  still needs a device check that tsnet reaches Running / Online.
