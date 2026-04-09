---
title: Usage
description: Manage Redmine issues with rmn — list, view, create, update, close, and delete.
---

# Usage

## Listing issues

```bash
rmn issue list                                    # List open issues (default limit: 25)
rmn issue list -p my-project                      # Filter by project
rmn issue list -s closed                          # Filter by status (open, closed, *, or status ID)
rmn issue list -a me                              # Issues assigned to you
rmn issue list -t 2                               # Filter by tracker ID
rmn issue list --sort updated_on:desc             # Sort by column
rmn issue list -l 50 --offset 100                 # Pagination
rmn issue list -p my-project -s closed -a me      # Combine filters
```

## Viewing issues

```bash
rmn issue view 42                                 # View issue details
```

## Creating issues

```bash
rmn issue create -p my-project -s "Bug report"
rmn issue create -p my-project -s "Feature request" -d "Detailed description" \
  -t 2 --priority 3 -a 5 --start-date 2025-01-01 --due-date 2025-03-31
```

All create flags: `--project/-p` (required), `--subject/-s` (required), `--description/-d`, `--tracker/-t`, `--priority`, `--assignee/-a`, `--category`, `--version`, `--parent`, `--start-date`, `--due-date`, `--estimated-hours`, `--done-ratio`.

## Updating issues

```bash
rmn issue update 42 --status 3                    # Change status
rmn issue update 42 -n "Work in progress"         # Add a note
rmn issue update 42 --done-ratio 50 --priority 2  # Update multiple fields
```

Only specified fields are changed; omitted fields remain unchanged. All create flags are available, plus `--status` and `--notes/-n`.

## Closing issues

```bash
rmn issue close 42                                # Close (status ID 5 by default)
rmn issue close 42 --status 6                     # Close with custom status ID
rmn issue close 42 -n "Fixed in v1.2"             # Close with a note
```

## Deleting issues

```bash
rmn issue delete 42                               # Delete with confirmation prompt
rmn issue delete 42 -y                            # Skip confirmation
```

## Command Aliases

| Command              | Aliases        |
|----------------------|----------------|
| `rmn issue list`     | `ls`           |
| `rmn issue view`     | `show`, `get`  |
| `rmn issue create`   | `new`          |
| `rmn issue delete`   | `rm`           |

```bash
rmn issue ls                    # Same as: rmn issue list
rmn issue show 42               # Same as: rmn issue view 42
rmn issue new -p proj -s "Bug"  # Same as: rmn issue create ...
rmn issue rm 42                 # Same as: rmn issue delete 42
```

## Global Flags

| Flag             | Description                              |
|------------------|------------------------------------------|
| `--output`       | Output format: `table` (default) or `json` |
| `--redmine-url`  | Override Redmine instance URL            |
| `--api-key`      | Override Redmine API key                 |

## JSON Output

Use `--output json` on any command for machine-readable output, useful for scripting and piping:

```bash
rmn issue list --output json                      # JSON array of issues
rmn issue view 42 --output json                   # Full issue as JSON
rmn issue list -p my-project --output json | jq '.issues[].subject'
```
