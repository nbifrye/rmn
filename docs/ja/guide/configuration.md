---
title: 設定
description: rmnの認証設定、設定ファイルの場所、コマンド単位のオーバーライド。
---

# 設定

## 認証

```bash
# フラグを使って設定
rmn auth login --url https://your-redmine.example.com --api-key YOUR_API_KEY

# または対話形式で
rmn auth login

# 設定を確認
rmn auth status
```

## 設定ファイル

設定は `~/.config/rmn/config.json`（`$XDG_CONFIG_HOME`が設定されている場合は `$XDG_CONFIG_HOME/rmn/config.json`）に保存されます：

```json
{
  "redmine_url": "https://your-redmine.example.com",
  "api_key": "your-api-key-here"
}
```

設定ファイルは `0600` パーミッションで作成されます。rmnは安全でないパーミッションの設定ファイルの読み取りを拒否します。

## コマンド単位のオーバーライド

グローバルフラグ `--redmine-url` と `--api-key` で、単一コマンドの保存済み設定をオーバーライドできます：

```bash
rmn issue list --redmine-url https://other-redmine.example.com --api-key OTHER_KEY
```

::: warning
平文HTTP接続時、rmnはAPIキーが平文で送信されることを警告します。常にHTTPSを使用してください。
:::
