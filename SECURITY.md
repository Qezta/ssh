# Security Policy

## Supported Versions

Only the latest release of this project receives security updates.

| Version | Supported |
| ------- | --------- |
| latest | ✓ |
| older | ✗ |

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Report vulnerabilities privately by emailing **security@qezta.com** with:

- A description of the vulnerability and its potential impact
- Steps to reproduce or proof-of-concept code
- Any suggested mitigations (optional)

You can expect an acknowledgement within **48 hours** and a status update within **7 days**. We will coordinate a fix and disclosure timeline with you before going public.

## Scope

This policy covers the source code in this repository. Vulnerabilities in third-party dependencies should be reported upstream to the respective maintainers.

This server intentionally accepts SSH connections from any username with any
key or password — it never grants a real shell, only a sandboxed TUI. There
is no meaningful auth boundary to bypass by design; reports about "anyone can
log in" are expected behavior, not a vulnerability. Reports about escaping
the TUI into a real shell, resource exhaustion, or crashing the host process
are in scope.
