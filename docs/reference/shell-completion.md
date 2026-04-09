---
title: Shell Completion
description: Set up shell auto-completion for rmn in Bash, Zsh, Fish, and PowerShell.
---

# Shell Completion

```bash
# Bash
source <(rmn completion bash)

# Zsh
rmn completion zsh > "${fpath[1]}/_rmn"

# Fish
rmn completion fish | source

# PowerShell
rmn completion powershell | Out-String | Invoke-Expression
```
