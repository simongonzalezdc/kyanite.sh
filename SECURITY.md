# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

Only the current major version receives security updates and bug fixes.

## Reporting a Vulnerability

### How to Report

If you discover a security vulnerability in Focus.sh, please report it responsibly:

1. **Email**: Send a detailed report to `security@focus.sh`
2. **GitHub**: Use the [Security Advisory](https://github.com/kyanite/focus/security/advisories) feature
3. **Private Issue**: Create a private issue in the repository

### What to Include

Please provide as much information as possible:

- **Vulnerability Type**: (e.g., SQL injection, XSS, authentication bypass)
- **Impact**: What could happen if exploited?
- **Steps to Reproduce**: Detailed steps to trigger the vulnerability
- **Proof of Concept**: Code or screenshots demonstrating the issue
- **Affected Versions**: Which versions are vulnerable?
- **Environment**: OS, Go version, AI provider used

### Response Timeline

- **Initial Response**: Within 48 hours
- **Detailed Assessment**: Within 7 days
- **Fix Released**: Within 30 days (depending on severity)
- **Public Disclosure**: After fix is released, unless coordinated disclosure is needed

## Security Model

### Data Protection

Focus.sh is designed with security in mind:

#### Local Data Storage
- **Tasks**: Stored locally in JSON format
- **Configuration**: Local configuration files
- **AI Prompts**: Processed locally (when using Ollama)
- **Logs**: Local log files only

#### AI Integration Security

**Ollama (Local LLM)**
- ✅ All processing happens locally
- ✅ No data sent to external servers
- ✅ Full control over model and data

**OpenRouter (Fallback)**
- ⚠️ API requests sent to OpenRouter servers
- ⚠️ Data processed by third-party models
- ✅ No data stored by OpenRouter beyond session
- ✅ API key protection with environment variables

#### Network Communications
- **Ollama**: Local network only (localhost)
- **OpenRouter**: HTTPS API calls only
- **Google Calendar**: OAuth 2.0 authentication
- **MCP Server**: Local TCP connections

### Authentication & Authorization

#### API Keys
- Stored in environment variables (`.env` file)
- Not logged or tracked
- Local file system access only

#### Google Calendar Integration
- OAuth 2.0 flow with user consent
- Limited calendar permissions
- Token storage in local configuration
- Token refresh handled automatically

### File System Access

Focus.sh accesses these locations:
- **Configuration**: `~/.config/focus/` or user-specified
- **Data**: `~/.local/share/focus/` or user-specified  
- **Logs**: `~/.local/log/focus/` or user-specified
- **Temp**: System temp directory for operations

### Input Validation

- **User Input**: Sanitized for command-line injection
- **File Paths**: Validated and canonicalized
- **AI Prompts**: Limited length and content filtering
- **Configuration**: Validated on startup

## Security Best Practices for Users

### Installation Security

1. **Download from official sources only**:
   - GitHub releases: https://github.com/kyanite/focus/releases
   - Package managers (when available)

2. **Verify releases**:
   ```bash
   # Verify checksum (when provided)
   sha256sum focus-v1.0.0-linux-amd64.tar.gz
   
   # Verify GPG signature (when available)
   gpg --verify focus-v1.0.0-linux-amd64.tar.gz.sig
   ```

3. **Avoid running as root** unless absolutely necessary

### Configuration Security

1. **Protect API keys**:
   ```bash
   # Set proper file permissions
   chmod 600 ~/.config/focus/.env
   
   # Never commit .env files to version control
   echo ".env" >> .gitignore
   ```

2. **Use local AI when possible**:
   - Prefer Ollama over OpenRouter for sensitive data
   - Review OpenRouter's privacy policy

3. **Secure Google Calendar access**:
   - Review requested permissions
   - Revoke access when no longer needed

### Operational Security

1. **Keep software updated**:
   ```bash
   # Check for updates
   focus --version
   
   # Upgrade to latest version
   # Follow installation instructions for your platform
   ```

2. **Regular backups**:
   ```bash
   # Backup data and configuration
   cp -r ~/.config/focus ~/.config/focus.backup
   cp -r ~/.local/share/focus ~/.local/share/focus.backup
   ```

3. **Monitor logs for suspicious activity**:
   ```bash
   # Check recent logs
   tail -f ~/.local/log/focus/focus.log
   ```

## Known Security Considerations

### AI Model Security
- **Local Models**: Data stays local, but model integrity depends on source
- **Remote Models**: Third-party processing, potential data exposure
- **Model Updates**: Verify model sources before updating

### Network Dependencies
- **Ollama**: Requires local installation and network access
- **OpenRouter**: Requires internet connectivity
- **Google Calendar**: Requires OAuth authentication and internet

### File System
- **Configuration Files**: Plain text storage (consider encryption for sensitive data)
- **Task Data**: JSON format (readable by anyone with file access)
- **Logs**: May contain sensitive information (AI prompts, responses)

## Threat Model

### Primary Threats Addressed
1. **Data Exposure**: Minimized through local processing and secure defaults
2. **Unauthorized Access**: Limited through file permissions and authentication
3. **Code Injection**: Mitigated through input validation and sanitization
4. **Man-in-the-Middle**: Prevented through HTTPS for external communications

### Limitations
1. **Physical Access**: No protection against physical access to files
2. **System Compromise**: Relies on underlying OS security
3. **AI Model Attacks**: Dependent on model provider security
4. **Social Engineering**: User education required for phishing prevention

## Security Updates

### How Updates Are Handled
- **Critical Security Issues**: Immediate patches and releases
- **Non-Critical Issues**: Included in next scheduled release
- **Dependencies**: Regular updates for Go modules and dependencies

### Notification Methods
- **GitHub Releases**: Security updates announced in release notes
- **Security Advisories**: Published for critical vulnerabilities
- **Documentation**: Updated security best practices

## Responsible Disclosure

This security policy is designed to encourage responsible disclosure and ensure that security vulnerabilities are addressed promptly and appropriately.

### Coordinated Disclosure
- We coordinate with security researchers to ensure responsible disclosure
- Public disclosure typically occurs after a fix is available
- Exceptions made for active exploits or user safety concerns

### Legal Protection
- We commit to not take legal action against researchers who follow this policy
- Research conducted in good faith is protected and appreciated

## Questions

For security-related questions:
- **General Questions**: Create a GitHub Discussion
- **Vulnerability Reports**: Follow the reporting process above
- **Security Team**: `security@focus.sh`

---

Thank you for helping keep Focus.sh secure! 🛡️