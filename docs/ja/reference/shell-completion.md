---
title: シェル補完
description: Bash、Zsh、Fish、PowerShellでrmnのシェル自動補完を設定。
---

# シェル補完

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
