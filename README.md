# Dighub

> Advanced GitHub Dorking & Secret Hunting Tool

Dighub is a powerful CLI tool that performs advanced GitHub dorking to detect exposed secrets, credentials, webhooks and sensitive files inside public repositories.

---

## Installation

```bash
go install github.com/ahmetartuc/dighub@latest
sudo cp ~/go/bin/dighub /usr/bin
```
```bash
dighub -org <github-org> -token <your_github_pat>
```

---

## Disclaimer
This tool is intended for educational and authorized security research only.
Do not use it against systems or organizations without explicit permission.
