# Thelo

Thelo is the private network client of Qavats Ranch, LLC: **`thelod`**, the node
agent, and **`thelo`**, its command line. Devices running it form a **thelonet**, a private WireGuard®
network run by the [Thelo Controller](https://github.com/ashleemikeyoung/headscale).

This repository is a fork of [Tailscale](https://github.com/tailscale/tailscale). It is not
affiliated with or endorsed by Tailscale Inc.; "Tailscale" is their trademark. Tailscale's own
README is kept as [README.upstream.md](README.upstream.md).

## Branches

| Branch | What it is |
|---|---|
| `thelo` (default) | a Tailscale stable release with our commits on top (currently v1.104.5) |
| `main` | a plain mirror of upstream `main`, used to pick up security fixes |

## What is different from Tailscale

- **The names people read:** the CLI is `thelo` and the daemon `thelod`; a tailnet is a thelonet;
  Funnel is **Halo** (`thelo halo`); `thelo version` names Thelo and Qavats Ranch; `thelo licenses`
  prints Tailscale's licence in full. The brand layer is `cmd/tailscale/cli/brand.go`.
- **Its own daemon identity:** socket, state directory, systemd unit and macOS launchd label
  (`com.qavats.thelod`) are Thelo's, so it never clashes with a Tailscale install.
- **Unchanged, on purpose:** the protocol, the state format, the `tailscale0` interface, Go package
  paths and `--json` output. Nodes keep their identity when they move from tailscaled, and the
  stock Tailscale apps can still join a Thelo Controller.

We deploy `thelod` with `--no-logs-no-support`, so nothing is uploaded to Tailscale's servers.

## Building

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ./build_dist.sh --strip -o thelod ./cmd/tailscaled
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ./build_dist.sh --strip -o thelo ./cmd/tailscale
```

macOS works the same with `GOOS=darwin GOARCH=arm64`.

## Following upstream

Our commits are kept small so they rebase onto each new stable release:
`git rebase --onto vX.Y.Z <previous tag> thelo`. They touch `cmd/tailscale/cli`,
`cmd/tailscaled/install_darwin.go`, `paths`, `safesocket`, `license_test.go` and the licence
and README files.

## Ownership and licence

Thelo is the work and property of Qavats Ranch, LLC. Copyright © 2026 Qavats Ranch, LLC.

It is free software under the **GNU Affero General Public License, version 3 or later**
([LICENSE](LICENSE)). Code that comes from Tailscale keeps its BSD 3-Clause licence
([LICENSE-BSD-3-Clause](LICENSE-BSD-3-Clause)) and Tailscale's patent grant ([PATENTS](PATENTS));
[NOTICE](NOTICE) explains which covers what. The names and marks Thelo, Qavats and ElRoi
belong to Qavats Ranch, LLC and are not covered by the licence.

WireGuard is a registered trademark of Jason A. Donenfeld.
