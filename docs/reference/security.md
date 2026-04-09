---
title: Security
description: Security features in rmn — TLS enforcement, file permissions, API key protection.
---

# Security

- **TLS 1.2+** — all HTTPS connections enforce TLS 1.2 as the minimum version
- **Secure file permissions** — config files are created with `0600` and config directories with `0700`; rmn refuses to read files with insecure permissions
- **API key protection** — the API key header is automatically removed on cross-origin redirects to prevent credential leakage
- **HTTP warning** — a warning is displayed when connecting over plain HTTP
