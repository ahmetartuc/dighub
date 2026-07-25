# DigHub

> **Advanced GitHub Dorking & Secret Hunting Tool**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

DigHub is a CLI tool that scans a GitHub organization or user's public repositories for exposed secrets, credentials, and sensitive files. It ships with 100+ dork patterns, concurrent scanning, and multiple output formats (terminal, JSON, CSV, HTML).

## Install

```bash
go install github.com/ahmetartuc/dighub@latest
```

Or build from source:

```bash
git clone https://github.com/ahmetartuc/dighub.git
cd dighub && go build -o dighub
```

## Usage

```bash
# Scan an organization
dighub -org <github-org> -token <your_github_pat>

# Scan a user
dighub -user <github-user> -token <your_github_pat>

# High-priority findings only, exported as HTML
dighub -org myorg -token ghp_xxx -priority high -output html -out-file report.html

# Faster scan, filter to specific categories
dighub -org myorg -token ghp_xxx -workers 10 -include "AWS,GitHub,SSH"
```

Both `-flag` and `--flag` syntax work.

## Options

| Flag | Description |
|------|-------------|
| `-org, -o` | GitHub organization to scan |
| `-user, -u` | GitHub user to scan (alternative to `-org`) |
| `-token, -t` | GitHub Personal Access Token **(required)** |
| `-output, -f` | Output format: `terminal`, `json`, `csv`, `html` (default: `terminal`) |
| `-out-file, -w` | Output file path (auto-generated if omitted) |
| `-priority, -p` | Filter by priority: `all`, `high`, `medium`, `low` (default: `all`) |
| `-include, -i` | Include only these dork categories (comma-separated) |
| `-exclude, -e` | Exclude these dork categories (comma-separated) |
| `-workers, -W` | Concurrent workers, 1–20 (default: `5`) |
| `-delay, -d` | Delay between requests in seconds (default: `2`) |
| `-verbose, -v` | Verbose output |
| `-quiet, -q` | Quiet mode — only show matches |

## GitHub Token

1. Go to [Settings → Developer Settings → Personal Access Tokens](https://github.com/settings/tokens).
2. Generate a classic token with the `public_repo` scope.
3. Pass it with the `-token` flag.

## Disclaimer

This tool is intended for **authorized security research and educational purposes only**. Only scan repositories and organizations you have permission to test, and always comply with applicable laws and GitHub's Terms of Service. The authors are not responsible for misuse.

## License

MIT — see [LICENSE](LICENSE). Contributions welcome; see [CONTRIBUTING.md](CONTRIBUTING.md).
