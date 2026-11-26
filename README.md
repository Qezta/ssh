<div id="top">
    <div align="center">
        <img alt='An abstract image of a donut-like object' title='Qezta' height=250 width=250 src='https://qezta.com/qezta.gif' /><sup>™</sup>
        <p>The Qezta torus logo, including its animated version, is a trademark of Qezta (OPC) Pvt. Ltd.</p>
        <p>© 2025–2026 Qezta (OPC) Pvt. Ltd. All rights reserved.</p>
        <h1 align='center'>ssh</h1>
        <strong>ssh divit@qezta.com: SSH terminal portfolio mirroring divit.qezta.com, built on wish + Bubble Tea</strong>
    </div>

</div>

---

<div align='center'>
    <a href="https://github.com/Qezta/ssh/stargazers">
        <img src="https://img.shields.io/github/stars/Qezta/ssh?&style=for-the-badge&logoColor=white&color=purple&logo=starship" alt="stars"/>
    </a>
    <a href="https://github.com/Qezta/ssh/">
        <img src="https://img.shields.io/github/repo-size/Qezta/ssh?&style=for-the-badge&logoColor=white&color=purple&logo=github" alt="size"/>
    </a>
    <a href="https://github.com/Qezta/ssh/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/Qezta/ssh?&style=for-the-badge&logoColor=white&color=purple&logo=unlicense" alt="license"/>
    </a>
    <a href="https://github.com/nixos/nixpkgs">
        <img src="https://img.shields.io/badge/Nixpkgs-unstable-blue.svg?style=for-the-badge&logo=NixOS&logoColor=white&color=purple" alt="nixpkgs"/>
    </a>
    <img src="https://img.shields.io/github/languages/top/Qezta/ssh?style=for-the-badge&color=purple" alt="repo-top-language"/>
</div>

---

## Contents

- [Overview](#overview)
- [How it works](#how-it-works)
- [Stack](#stack)
- [Commands](#commands)
- [Development](#development)
- [Nix](#nix)
- [Security](#security)
- [License](#license)

---

## Overview

**`ssh divit@qezta.com`** — a real SSH server that drops you straight into a
terminal portfolio, no real shell access involved.

This mirrors the fake browser terminal on
[divit.qezta.com](https://divit.qezta.com) (see
[`Terminal.svelte`](https://github.com/Qezta/divit/blob/main/src/lib/components/Terminal.svelte)
and
[`terminal-commands.ts`](https://github.com/Qezta/divit/blob/main/src/lib/data/terminal-commands.ts)
in [`Qezta/divit`](https://github.com/Qezta/divit)), but backed by a genuine
SSH connection instead of a JS widget. The server never grants a real shell,
only a sandboxed [Bubble Tea](https://github.com/charmbracelet/bubbletea)
TUI — but which persona you land in depends on the username you connect as:

```sh
ssh divit@qezta.com   # Divit's personal portfolio (div@ works too)
ssh hi@qezta.com      # Qezta org landing (hello@ works too)
```

There's a single SSH listener on `qezta.com`, so identity is chosen by
username rather than running a separate server per persona — see
`identityFor` in [`internal/tui/commands.go`](./internal/tui/commands.go).
`divit`/`div` get the personal identity and `hi`/`hello` get the org
identity. Any other username (`qezta`, `guest`, ...) gets a short guide
pointing at those two and is disconnected.

## How it works

```mermaid
flowchart LR
  client["ssh &lt;user&gt;@qezta.com"] --> wish["wish SSH server"]
  wish --> mw["middleware<br/>Bubble Tea · logging"]
  mw --> id{"identityFor(user)"}
  id -- "divit / div" --> me["personal portfolio TUI"]
  id -- "anything else" --> org["Qezta org TUI"]
```

No shell is ever started: every session is a Bubble Tea program.

## Stack

- [`wish`](https://github.com/charmbracelet/wish) — SSH server middleware
- [`bubbletea`](https://github.com/charmbracelet/bubbletea) — TUI framework
- [`lipgloss`](https://github.com/charmbracelet/lipgloss) — styling

## Commands

Shared: `clear`, `nix-build`, `cowsay <text>`.

`divit`/`div` (personal): `help`, `about`, `projects
[--ml|--infrastructure|--functional]`, `skills [--expert]`, `contact`,
`experience`, `resume`.

Anything else (org): `help`, `about`, `repos`, `contact`.

## Development

```sh
go run .                        # serves on :23234 by default
ssh -p 23234 divit@localhost    # personal identity
ssh -p 23234 hi@localhost       # org identity
```

| Variable | Default | Purpose |
| ------------------- | ---------------------- | --------------------------------------- |
| `SSH_HOST` | `0.0.0.0` | Bind address |
| `SSH_PORT` | `23234` | Bind port (proxy/forward 22 to this in prod) |
| `SSH_HOST_KEY_PATH` | `.ssh/qezta_ed25519` | Where the server's host key is generated/read |

```sh
go build ./...   # build
go vet ./...     # vet
gofmt -l .       # format check
```

## Nix

```sh
nix develop      # devshell with go, gopls, golangci-lint
nix build        # static server binary at result/bin/ssh
nix build .#image  # Linux only: OCI image tarball (dockerTools)
```

## Deployment

`packages.<linux>.image` runs on the home-lab k3s cluster (host T3, see
`hosts/nixos/T3/ssh-portfolio.nix` in `infra-nixCfg`). The image is preloaded
into k3s from the Nix store, so no registry is involved. Its tag is the
content hash, so every real change rolls out a new pod. It runs as UID 65532
and reads or generates its host key at `/data/ssh_host_ed25519_key`, which
lives on a persistent volume so clients' `known_hosts` entries stay valid.

## Security

This server intentionally accepts any username with any key or password —
see [SECURITY.md](./SECURITY.md) for what's in and out of scope.

## License

MIT — see [LICENSE](./LICENSE).

<div align="right">

[![][back-to-top]](#top)

</div>

[back-to-top]: https://img.shields.io/badge/-BACK_TO_TOP-151515?style=flat-square&color=purple
