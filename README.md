#  DigHub

> **Advanced GitHub Dorking & Secret Hunting Tool**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/ahmetartuc/dighub?style=for-the-badge)](https://github.com/ahmetartuc/dighub/stargazers)

Dighub is a powerful CLI tool that performs advanced GitHub dorking to detect exposed secrets, credentials, webhooks and sensitive files inside public repositories. With concurrent scanning, multiple output formats, and intelligent filtering, Dighub helps security researchers and DevOps teams identify security vulnerabilities efficiently.

## Features

- 🚀 **Concurrent Scanning** - Up to 20 parallel workers for 10x faster scanning
- 🎯 **Smart Filtering** - Filter by priority (high/medium/low), include/exclude patterns
- 📊 **Multiple Output Formats** - Terminal, JSON, CSV, HTML reports
- 🎨 **Beautiful Terminal Output** - Colored output with progress bars
- ⚡ **Rate Limit Handling** - Automatic retry with intelligent wait times
- 🔍 **100+ Dork Patterns** - Comprehensive detection for AWS, GitHub, SSH keys, databases, webhooks, and more
- 📈 **Detailed Statistics** - Track findings by priority and category
- 🎛️ **Flexible Configuration** - Extensive CLI flags for customization

## Installation

### From Source

```bash
git clone https://github.com/ahmetartuc/dighub.git
cd dighub
go mod download
go build -o dighub
sudo mv dighub /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/ahmetartuc/dighub@latest
```

## Quick Start

### Basic Usage

```bash
# Scan an organization
dighub -org <github-org> -token <your_github_pat>

# Scan a user
dighub -user <github-user> -token <your_github_pat>
```

### Advanced Usage

```bash
# High priority findings only, save as JSON
dighub -org myorg -token ghp_xxx -priority high -output json

# Concurrent scanning with 10 workers
dighub -org myorg -token ghp_xxx -workers 10

# Filter specific patterns
dighub -org myorg -token ghp_xxx -include "AWS,GitHub" -exclude "backup,log"

# Export to HTML report
dighub -org myorg -token ghp_xxx -output html -out-file report.html

# Quiet mode (URLs only)
dighub -org myorg -token ghp_xxx -quiet

# Verbose mode (detailed output)
dighub -org myorg -token ghp_xxx -verbose
```

## Command Line Options

### Required Flags
```
-org, -o          GitHub organization to scan
-user, -u         GitHub user to scan (alternative to org)
-token, -t        GitHub Personal Access Token (required)
```

### Output Options
```
-output, -f       Output format: terminal, json, csv, html (default: terminal)
-out-file, -w     Output file path (auto-generated if not specified)
-quiet, -q        Quiet mode - only show matches
-verbose, -v      Verbose output with detailed information
-no-color, -n     Disable colored output
```

### Filtering Options
```
-priority, -p     Priority level: all, high, medium, low (default: all)
-include, -i      Include specific dorks (comma-separated patterns)
-exclude, -e      Exclude specific dorks (comma-separated patterns)
```

### Performance Options
```
-workers, -W      Number of concurrent workers (1-20, default: 5)
-rate-limit, -r   Requests per minute (default: 30)
-delay, -d        Delay between requests in seconds (default: 2)
```

## Dork Categories

Dighub includes dorks for detecting:

### High Priority 🔴
- **AWS Credentials** - Access keys, secret keys, session tokens
- **GitHub Tokens** - Personal access tokens, OAuth tokens
- **SSH Keys** - Private keys (RSA, DSA, Ed25519)
- **Private Keys & Certificates** - PEM files, key files
- **Database Credentials** - Connection strings, passwords
- **Payment Gateway Secrets** - Stripe, PayPal, Braintree
- **API Keys** - OpenAI, Cloudflare, Vercel, and more

### Medium Priority 🟡
- **Webhooks** - Discord, Slack, Teams, Office 365
- **Email Services** - SendGrid, Mailgun, SMTP credentials
- **Cloud Services** - Firebase, Google Services, Azure
- **CI/CD Configs** - GitHub Actions, Travis, GitLab CI
- **Infrastructure** - Terraform, Kubernetes configs

### Low Priority 🔵
- **Configuration Files** - Settings, properties, ini files
- **Log Files** - Debug logs, error logs
- **Backup Files** - SQL dumps, database backups
- **History Files** - Bash history, zsh history

## Output Formats

### Terminal Output (Default)
Colored, organized output grouped by priority with progress bar.

### JSON Output
```json
{
  "scan_info": {
    "target": "myorg",
    "scan_date": "2024-01-14T...",
    "duration": "5m23s",
    "total_dorks": 100
  },
  "summary": {
    "total_matches": 45,
    "unique_files": 23,
    "high_priority": 12,
    "medium_priority": 20,
    "low_priority": 13
  },
  "findings": [...]
}
```

### CSV Output
Structured CSV with all finding details for easy analysis in spreadsheets.

### HTML Output
Beautiful, interactive HTML report with:
- Executive summary with statistics
- Color-coded priority badges
- Sortable findings
- Direct links to GitHub files
- Responsive design

## GitHub Token Setup

1. Go to [GitHub Settings → Developer Settings → Personal Access Tokens](https://github.com/settings/tokens)
2. Click "Generate new token (classic)"
3. Select scopes: `public_repo` (for public repos only)
4. Copy the token (starts with `ghp_`)
5. Use it with the `-token` flag

## Performance Tips

1. **Use concurrent workers** for faster scanning:
   ```bash
   dighub -org myorg -token xxx -workers 10
   ```

2. **Filter by priority** to focus on critical findings:
   ```bash
   dighub -org myorg -token xxx -priority high
   ```

3. **Use specific includes** to target what matters:
   ```bash
   dighub -org myorg -token xxx -include "AWS,GitHub,SSH"
   ```

4. **Adjust rate limits** based on your token limits:
   ```bash
   dighub -org myorg -token xxx -rate-limit 50 -delay 1
   ```

## Examples

### Example 1: Quick Security Audit
```bash
dighub -org mycompany -token ghp_xxx -priority high -output html
```

### Example 2: Comprehensive Scan with JSON Export
```bash
dighub -org mycompany -token ghp_xxx -workers 10 -output json -out-file security-audit.json
```

### Example 3: Focus on AWS Credentials
```bash
dighub -org mycompany -token ghp_xxx -include "AWS" -verbose
```

### Example 4: Exclude False Positives
```bash
dighub -org mycompany -token ghp_xxx -exclude "test,example,demo"
```

## Security Best Practices

1. **Never commit secrets** - Use environment variables or secret managers
2. **Rotate exposed credentials immediately** - If Dighub finds secrets, rotate them
3. **Use .gitignore** - Prevent sensitive files from being committed
4. **Enable GitHub secret scanning** - GitHub's built-in protection
5. **Regular audits** - Run Dighub regularly on your repos

## 📋 Scan Results Interpretation

### High Priority Findings 🔴
**Action Required**: Immediately rotate/revoke these credentials
- Direct access to critical systems
- Can cause data breaches or service disruptions

### Medium Priority Findings 🟡
**Review Required**: Assess risk and take appropriate action
- May provide indirect access or information disclosure
- Should be removed from public repositories

### Low Priority Findings 🔵
**Best Practice**: Clean up for security hygiene
- Generally configuration files or less sensitive data
- Should still be reviewed and removed if not needed

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Adding New Dorks

To add new dork patterns, edit `internal/dorks/dorks.go`:

```go
{Pattern: "filename:.env NEW_SECRET", Priority: PriorityHigh, Category: "Category", Description: "Description"},
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

**This tool is intended for authorized security research and educational purposes only.**

- Only use Dighub on repositories and organizations you have permission to test
- Do not use it for malicious purposes or unauthorized access
- The authors are not responsible for misuse of this tool
- Always comply with applicable laws and GitHub's Terms of Service
- Respect rate limits and API usage policies


**Made with ❤️ for the security community**

If you find this tool useful, please consider giving it a ⭐ on GitHub!
