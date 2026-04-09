---
title: Configuration
description: Configure rmn authentication, config file location, and per-command overrides.
---

# Configuration

## Authentication

```bash
# Set up with flags
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# Or interactively
rmn auth login

# Verify your configuration
rmn auth status
```

## Config file

Configuration is stored in `~/.config/rmn/config.json` (or `$XDG_CONFIG_HOME/rmn/config.json` if set):

```json
{
  "redmine_url": "https://your-redmine.example.com",
  "api_key": "your-api-key-here"
}
```

The config file is created with `0600` permissions. rmn refuses to read config files with insecure permissions.

## Per-command overrides

Global flags `--redmine-url` and `--api-key` override the stored configuration for a single command:

```bash
rmn issue list --redmine-url https://other-redmine.example.com --api-key OTHER_KEY
```

::: warning
When connecting over plain HTTP, rmn displays a warning that your API key is sent in plaintext. Always prefer HTTPS.
:::
