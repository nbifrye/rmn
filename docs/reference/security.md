---
title: Security
description: Security features in rmn — TLS enforcement, file permissions, API key protection.
---

# Security

rmn implements several security measures to protect your Redmine credentials and data in transit.

## TLS 1.2+

All HTTPS connections enforce TLS 1.2 as the minimum protocol version. This is configured on the HTTP transport used by every API request, preventing downgrade attacks to older, vulnerable TLS versions. Connections to servers that do not support TLS 1.2 or later will fail.

## Secure file permissions

Configuration files are created with `0600` (owner read/write only) and configuration directories with `0700` (owner read/write/execute only). When loading a config file, rmn checks the file permissions and **refuses to read files where group or other bits are set** (i.e., permissions more open than `0600`). If rmn detects insecure permissions, it returns an error message with the current permissions and a suggested `chmod 600` command.

## API key protection on redirects

The `X-Redmine-API-Key` header is automatically removed when the HTTP client follows a redirect to a different host. This prevents credential leakage if a Redmine instance redirects to a third-party server. Same-host redirects preserve the API key header. The client also enforces a maximum of 10 redirects to prevent redirect loops.

## HTTP warning

When the configured Redmine URL uses plain `http://` instead of `https://`, rmn prints a warning to stderr:

> Warning: using insecure HTTP connection. API key will be sent in plaintext.

Only `http` and `https` URL schemes are supported; any other scheme is rejected with an error.

## Response body size limit

API responses are read with a 10 MB size limit to prevent excessive memory consumption from unexpectedly large responses. Responses exceeding this limit are silently truncated.
