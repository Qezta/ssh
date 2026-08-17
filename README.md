<div align="center">

# ssh.qezta.com

**`ssh divit@qezta.com`** — a real SSH server that drops you straight into a
terminal portfolio, no real shell access involved.

</div>

---

This mirrors the fake browser terminal on
[divit.qezta.com](https://divit.qezta.com) (see
[`Terminal.svelte`](https://github.com/Qezta/divit/blob/main/src/lib/components/Terminal.svelte)
and
[`terminal-commands.ts`](https://github.com/Qezta/divit/blob/main/src/lib/data/terminal-commands.ts)
in [`Qezta/divit`](https://github.com/Qezta/divit)), but backed by a genuine
SSH connection instead of a JS widget. Connect as any username — the server
never grants a real shell, only a sandboxed [Bubble
Tea](https://github.com/charmbracelet/bubbletea) TUI.

```sh
ssh divit@qezta.com
```

## Stack

- [`wish`](https://github.com/charmbracelet/wish) — SSH server middleware
- [`bubbletea`](https://github.com/charmbracelet/bubbletea) — TUI framework
- [`lipgloss`](https://github.com/charmbracelet/lipgloss) — styling

## Commands

`help`, `about`, `projects [--ml|--infrastructure|--functional]`, `skills
[--expert]`, `contact`, `experience`, `resume`, `clear`, `nix-build`,
`cowsay <text>`.

## Development

```sh
go run .                        # serves on :23234 by default
ssh -p 23234 divit@localhost    # connect from another terminal
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
```

## Security

This server intentionally accepts any username with any key or password —
see [SECURITY.md](./SECURITY.md) for what's in and out of scope.

## License

MIT — see [LICENSE](./LICENSE).
