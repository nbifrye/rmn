---
layout: home
title: rmn — Redmine CLI Tool
titleTemplate: Command-Line Client for Redmine

hero:
  name: rmn
  text: Redmine CLI Tool
  tagline: A fast, open-source command-line client for Redmine written in Go. Manage issues from your terminal or let AI agents handle it via the built-in MCP server.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/nbifrye/rmn

features:
  - title: Full Issue Lifecycle
    details: List, view, create, update, close, and delete Redmine issues from the command line.
  - title: MCP Server for AI Agents
    details: Expose Redmine operations to AI assistants like Claude Code via the Model Context Protocol.
  - title: Multiple Output Formats
    details: Human-readable table (default) and machine-readable JSON for scripting and automation.
  - title: Flexible Filtering
    details: Filter by project, status, assignee, and tracker with sorting and pagination support.
  - title: GitLab CLI-Inspired Aliases
    details: Familiar shorthand commands (ls, show, get, new, rm) for faster workflows.
  - title: Six Installation Methods
    details: Homebrew, mise, Nix, Go install, pre-built binaries, and build from source.
  - title: Shell Completion
    details: Auto-complete for Bash, Zsh, Fish, and PowerShell.
  - title: XDG-Compliant Config
    details: Respects $XDG_CONFIG_HOME for config file placement.
  - title: Security Hardened
    details: TLS 1.2+ enforcement, secure file permissions (0600), and API key protection on redirects.
  - title: 100% Test Coverage
    details: Enforced in CI with every pull request.
  - title: Cross-Platform
    details: Pre-built binaries for Linux, macOS, and Windows on both amd64 and arm64.
---

::: warning
This project is experimental. The API and CLI interface are not yet stable and may introduce breaking changes without notice.
:::

## Quick Start

```bash
# 1. Install
brew tap nbifrye/rmn https://github.com/nbifrye/rmn.git
brew install nbifrye/rmn/rmn

# 2. Authenticate
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# 3. List your issues
rmn issue list -a me

# 4. View issue details
rmn issue view 42
```
