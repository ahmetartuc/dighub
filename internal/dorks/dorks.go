package dorks

import "strings"

// Priority levels for dorks
const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// Dork represents a search pattern with metadata
type Dork struct {
	Pattern     string
	Priority    string
	Description string
	Category    string
}

// GetDorks returns all dorks with metadata
func GetDorks() []Dork {
	return []Dork{
		// HIGH PRIORITY - Critical secrets and credentials
		{Pattern: "AWS_ACCESS_KEY_ID", Priority: PriorityHigh, Category: "AWS", Description: "AWS access key"},
		{Pattern: "AWS_SECRET_ACCESS_KEY", Priority: PriorityHigh, Category: "AWS", Description: "AWS secret key"},
		{Pattern: "filename:.env AWS_SECRET_ACCESS_KEY", Priority: PriorityHigh, Category: "AWS", Description: "AWS secret in .env"},
		{Pattern: "filename:.env AWS_ACCESS_KEY_ID", Priority: PriorityHigh, Category: "AWS", Description: "AWS access key in .env"},
		{Pattern: "ghp_", Priority: PriorityHigh, Category: "GitHub", Description: "GitHub personal access token"},
		{Pattern: "gho_", Priority: PriorityHigh, Category: "GitHub", Description: "GitHub OAuth token"},
		{Pattern: "ghu_", Priority: PriorityHigh, Category: "GitHub", Description: "GitHub user token"},
		{Pattern: "ghr_", Priority: PriorityHigh, Category: "GitHub", Description: "GitHub refresh token"},
		{Pattern: "filename:id_rsa", Priority: PriorityHigh, Category: "SSH", Description: "SSH private key"},
		{Pattern: "filename:id_dsa", Priority: PriorityHigh, Category: "SSH", Description: "DSA private key"},
		{Pattern: "filename:id_ed25519", Priority: PriorityHigh, Category: "SSH", Description: "Ed25519 private key"},
		{Pattern: "extension:pem PRIVATE KEY", Priority: PriorityHigh, Category: "Certificates", Description: "PEM private key"},
		{Pattern: "extension:key PRIVATE KEY", Priority: PriorityHigh, Category: "Certificates", Description: "Private key file"},
		{Pattern: "filename:.env PRIVATE_KEY", Priority: PriorityHigh, Category: "Encryption", Description: "Private key in .env"},
		{Pattern: "filename:.env SECRET_KEY", Priority: PriorityHigh, Category: "Encryption", Description: "Secret key in .env"},
		{Pattern: "filename:.env JWT_SECRET", Priority: PriorityHigh, Category: "Authentication", Description: "JWT secret"},
		{Pattern: "filename:.env DB_PASSWORD", Priority: PriorityHigh, Category: "Database", Description: "Database password"},
		{Pattern: "DATABASE_URL", Priority: PriorityHigh, Category: "Database", Description: "Database URL"},
		{Pattern: "filename:.env STRIPE_SECRET_KEY", Priority: PriorityHigh, Category: "Payment", Description: "Stripe secret key"},
		{Pattern: "filename:.env PAYPAL_CLIENT_SECRET", Priority: PriorityHigh, Category: "Payment", Description: "PayPal secret"},
		{Pattern: "filename:.env OPENAI_API_KEY", Priority: PriorityHigh, Category: "API Keys", Description: "OpenAI API key"},
		{Pattern: "filename:.git-credentials", Priority: PriorityHigh, Category: "Git", Description: "Git credentials"},
		{Pattern: "filename:.npmrc _auth", Priority: PriorityHigh, Category: "NPM", Description: "NPM authentication"},
		{Pattern: "filename:.dockercfg", Priority: PriorityHigh, Category: "Docker", Description: "Docker config"},
		{Pattern: "filename:.docker/config.json", Priority: PriorityHigh, Category: "Docker", Description: "Docker auth config"},
		{Pattern: "filename:.env ADMIN_PASSWORD", Priority: PriorityHigh, Category: "Admin", Description: "Admin password"},
		{Pattern: "filename:.env ROOT_PASSWORD", Priority: PriorityHigh, Category: "Admin", Description: "Root password"},
		{Pattern: "filename:.env MASTER_KEY", Priority: PriorityHigh, Category: "Encryption", Description: "Master key"},
		
		// MEDIUM PRIORITY - Sensitive configurations and webhooks
		{Pattern: "discord.com/api/webhooks", Priority: PriorityMedium, Category: "Webhooks", Description: "Discord webhook"},
		{Pattern: "discordapp.com/api/webhooks", Priority: PriorityMedium, Category: "Webhooks", Description: "Discord app webhook"},
		{Pattern: "hooks.slack.com/services", Priority: PriorityMedium, Category: "Webhooks", Description: "Slack webhook"},
		{Pattern: "outlook.office.com/webhook", Priority: PriorityMedium, Category: "Webhooks", Description: "Office 365 webhook"},
		{Pattern: "teams.microsoft.com/webhook", Priority: PriorityMedium, Category: "Webhooks", Description: "Teams webhook"},
		{Pattern: "filename:.env WEBHOOK_URL", Priority: PriorityMedium, Category: "Webhooks", Description: "Generic webhook URL"},
		{Pattern: "filename:.env DISCORD_WEBHOOK", Priority: PriorityMedium, Category: "Webhooks", Description: "Discord webhook in .env"},
		{Pattern: "filename:.env SLACK_WEBHOOK_URL", Priority: PriorityMedium, Category: "Webhooks", Description: "Slack webhook in .env"},
		{Pattern: "filename:.env API_KEY", Priority: PriorityMedium, Category: "API Keys", Description: "Generic API key"},
		{Pattern: "filename:.env TOKEN", Priority: PriorityMedium, Category: "Tokens", Description: "Generic token"},
		{Pattern: "filename:.env MAIL_PASSWORD", Priority: PriorityMedium, Category: "Email", Description: "Email password"},
		{Pattern: "filename:.env SMTP_PASSWORD", Priority: PriorityMedium, Category: "Email", Description: "SMTP password"},
		{Pattern: "filename:.env SENDGRID_API_KEY", Priority: PriorityMedium, Category: "Email", Description: "SendGrid API key"},
		{Pattern: "filename:.env MAILGUN_API_KEY", Priority: PriorityMedium, Category: "Email", Description: "Mailgun API key"},
		{Pattern: "filename:.env TWILIO_AUTH_TOKEN", Priority: PriorityMedium, Category: "SMS", Description: "Twilio auth token"},
		{Pattern: "filename:.env FIREBASE", Priority: PriorityMedium, Category: "Firebase", Description: "Firebase config"},
		{Pattern: "filename:firebase.json", Priority: PriorityMedium, Category: "Firebase", Description: "Firebase JSON"},
		{Pattern: "filename:firebase-adminsdk.json", Priority: PriorityMedium, Category: "Firebase", Description: "Firebase Admin SDK"},
		{Pattern: "filename:google-services.json", Priority: PriorityMedium, Category: "Google", Description: "Google services config"},
		{Pattern: "filename:client_secret.json", Priority: PriorityMedium, Category: "OAuth", Description: "OAuth client secret"},
		{Pattern: "filename:.aws/credentials", Priority: PriorityMedium, Category: "AWS", Description: "AWS credentials file"},
		{Pattern: "filename:credentials.json", Priority: PriorityMedium, Category: "Credentials", Description: "Generic credentials"},
		{Pattern: "filename:secrets.yml", Priority: PriorityMedium, Category: "Secrets", Description: "Secrets YAML"},
		{Pattern: "filename:secrets.yaml", Priority: PriorityMedium, Category: "Secrets", Description: "Secrets YAML"},
		{Pattern: "filename:wp-config.php", Priority: PriorityMedium, Category: "WordPress", Description: "WordPress config"},
		{Pattern: "filename:.env MONGO_URI", Priority: PriorityMedium, Category: "Database", Description: "MongoDB URI"},
		{Pattern: "filename:.env MONGODB_URI", Priority: PriorityMedium, Category: "Database", Description: "MongoDB URI"},
		{Pattern: "filename:.env POSTGRES_PASSWORD", Priority: PriorityMedium, Category: "Database", Description: "PostgreSQL password"},
		{Pattern: "filename:.env MYSQL_ROOT_PASSWORD", Priority: PriorityMedium, Category: "Database", Description: "MySQL root password"},
		{Pattern: "filename:.env REDIS_PASSWORD", Priority: PriorityMedium, Category: "Database", Description: "Redis password"},
		{Pattern: "filename:database.yml", Priority: PriorityMedium, Category: "Database", Description: "Database YAML config"},
		{Pattern: "filename:.github/workflows token", Priority: PriorityMedium, Category: "CI/CD", Description: "GitHub Actions token"},
		{Pattern: "filename:.github/workflows GITHUB_TOKEN", Priority: PriorityMedium, Category: "CI/CD", Description: "GitHub token in workflow"},
		{Pattern: "filename:.travis.yml", Priority: PriorityMedium, Category: "CI/CD", Description: "Travis CI config"},
		{Pattern: "filename:.gitlab-ci.yml", Priority: PriorityMedium, Category: "CI/CD", Description: "GitLab CI config"},
		{Pattern: "filename:terraform.tfvars", Priority: PriorityMedium, Category: "Infrastructure", Description: "Terraform variables"},
		{Pattern: "filename:kubeconfig", Priority: PriorityMedium, Category: "Kubernetes", Description: "Kubernetes config"},
		{Pattern: "filename:.kube/config", Priority: PriorityMedium, Category: "Kubernetes", Description: "Kubectl config"},
		{Pattern: "filename:.env CLOUDFLARE_API_KEY", Priority: PriorityMedium, Category: "CDN", Description: "Cloudflare API key"},
		{Pattern: "filename:.env VERCEL_API_KEY", Priority: PriorityMedium, Category: "Hosting", Description: "Vercel API key"},
		{Pattern: "filename:.env NETLIFY_AUTH_TOKEN", Priority: PriorityMedium, Category: "Hosting", Description: "Netlify auth token"},
		{Pattern: "filename:.env HEROKU_API_KEY", Priority: PriorityMedium, Category: "Hosting", Description: "Heroku API key"},
		
		// LOW PRIORITY - General configuration files and patterns
		{Pattern: "filename:config.json api_key", Priority: PriorityLow, Category: "Config", Description: "API key in config"},
		{Pattern: "filename:settings.py SECRET_KEY", Priority: PriorityLow, Category: "Django", Description: "Django secret key"},
		{Pattern: "filename:application.yml token", Priority: PriorityLow, Category: "Config", Description: "Token in application config"},
		{Pattern: "filename:settings.yaml token", Priority: PriorityLow, Category: "Config", Description: "Token in settings"},
		{Pattern: "filename:.netrc password", Priority: PriorityLow, Category: "Network", Description: "Netrc password"},
		{Pattern: "filename:.ssh/config", Priority: PriorityLow, Category: "SSH", Description: "SSH config"},
		{Pattern: "filename:authorized_keys", Priority: PriorityLow, Category: "SSH", Description: "SSH authorized keys"},
		{Pattern: "filename:package.json token", Priority: PriorityLow, Category: "NPM", Description: "Token in package.json"},
		{Pattern: "filename:.npmrc authToken", Priority: PriorityLow, Category: "NPM", Description: "NPM auth token"},
		{Pattern: "filename:.pypirc password", Priority: PriorityLow, Category: "Python", Description: "PyPI password"},
		{Pattern: "filename:config.js password", Priority: PriorityLow, Category: "Config", Description: "Password in JS config"},
		{Pattern: "filename:config.json password", Priority: PriorityLow, Category: "Config", Description: "Password in JSON config"},
		{Pattern: "filename:settings.ini", Priority: PriorityLow, Category: "Config", Description: "INI settings file"},
		{Pattern: "filename:application.properties", Priority: PriorityLow, Category: "Config", Description: "Properties file"},
		{Pattern: "filename:.bash_history password", Priority: PriorityLow, Category: "History", Description: "Password in bash history"},
		{Pattern: "filename:.zsh_history token", Priority: PriorityLow, Category: "History", Description: "Token in zsh history"},
		{Pattern: "extension:log password", Priority: PriorityLow, Category: "Logs", Description: "Password in logs"},
		{Pattern: "extension:log token", Priority: PriorityLow, Category: "Logs", Description: "Token in logs"},
		{Pattern: "extension:log secret", Priority: PriorityLow, Category: "Logs", Description: "Secret in logs"},
		{Pattern: "extension:json password", Priority: PriorityLow, Category: "JSON", Description: "Password in JSON"},
		{Pattern: "extension:yaml password", Priority: PriorityLow, Category: "YAML", Description: "Password in YAML"},
		{Pattern: "extension:txt secret", Priority: PriorityLow, Category: "Text", Description: "Secret in text file"},
		{Pattern: "extension:ini password", Priority: PriorityLow, Category: "INI", Description: "Password in INI"},
		{Pattern: "filename:debug.log", Priority: PriorityLow, Category: "Logs", Description: "Debug log file"},
		{Pattern: "filename:error.log", Priority: PriorityLow, Category: "Logs", Description: "Error log file"},
		{Pattern: "filename:local.env", Priority: PriorityLow, Category: "Environment", Description: "Local environment"},
		{Pattern: "filename:prod.env", Priority: PriorityLow, Category: "Environment", Description: "Production environment"},
		{Pattern: "filename:staging.env", Priority: PriorityLow, Category: "Environment", Description: "Staging environment"},
		{Pattern: "filename:backup.sql", Priority: PriorityLow, Category: "Backup", Description: "SQL backup"},
		{Pattern: "filename:dump.sql", Priority: PriorityLow, Category: "Backup", Description: "SQL dump"},
		{Pattern: "filename:db.dump", Priority: PriorityLow, Category: "Backup", Description: "Database dump"},
		{Pattern: "filename:backup.tar", Priority: PriorityLow, Category: "Backup", Description: "Tar backup"},
		{Pattern: "extension:bak SECRET_KEY", Priority: PriorityLow, Category: "Backup", Description: "Secret in backup"},
		{Pattern: "filename:.env~", Priority: PriorityLow, Category: "Backup", Description: "Backup .env file"},
		{Pattern: "extension:swp SECRET", Priority: PriorityLow, Category: "Temp", Description: "Secret in swap file"},
		{Pattern: "extension:orig SECRET", Priority: PriorityLow, Category: "Temp", Description: "Secret in orig file"},
	}
}

// FilterByPriority returns dorks filtered by priority level
func FilterByPriority(dorks []Dork, priority string) []Dork {
	if priority == "all" {
		return dorks
	}
	
	var filtered []Dork
	for _, dork := range dorks {
		if dork.Priority == priority {
			filtered = append(filtered, dork)
		}
	}
	return filtered
}

// FilterByInclude returns dorks that match include patterns
func FilterByInclude(dorks []Dork, patterns []string) []Dork {
	if len(patterns) == 0 {
		return dorks
	}
	
	var filtered []Dork
	for _, dork := range dorks {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(dork.Pattern), strings.ToLower(pattern)) || 
			   strings.Contains(strings.ToLower(dork.Category), strings.ToLower(pattern)) || 
			   strings.Contains(strings.ToLower(dork.Description), strings.ToLower(pattern)) {
				filtered = append(filtered, dork)
				break
			}
		}
	}
	return filtered
}

// FilterByExclude returns dorks that don't match exclude patterns
func FilterByExclude(dorks []Dork, patterns []string) []Dork {
	if len(patterns) == 0 {
		return dorks
	}
	
	var filtered []Dork
	for _, dork := range dorks {
		exclude := false
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(dork.Pattern), strings.ToLower(pattern)) || 
			   strings.Contains(strings.ToLower(dork.Category), strings.ToLower(pattern)) || 
			   strings.Contains(strings.ToLower(dork.Description), strings.ToLower(pattern)) {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, dork)
		}
	}
	return filtered
}
