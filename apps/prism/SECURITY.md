# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in Prism.sh, please report it privately:

### How to Report

1. **Email:** Send details to security@kyanite.sh (if available)
2. **GitHub Security Advisory:** Use GitHub's [private vulnerability reporting](https://github.com/simongonzalezdc/Prism.sh/security/advisories/new)

### What to Include

Please include the following information:

- **Type of issue:** (e.g., buffer overflow, SQL injection, cross-site scripting)
- **Full paths of source file(s)** related to the manifestation of the issue
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the issue:** including how an attacker might exploit it

### What to Expect

- **Response Time:** We aim to respond within 48 hours
- **Status Updates:** We'll keep you informed of progress
- **Credit:** We'll credit you in the security advisory (unless you prefer otherwise)
- **Fix Timeline:** We'll work to release a fix as quickly as possible

## Security Best Practices for Users

### Safe Installation

```bash
# Always verify the source
git clone https://github.com/simongonzalezdc/Prism.sh.git
cd Prism.sh
make build

# Or build from source
git clone https://github.com/simongonzalezdc/Prism.sh.git
cd prism
go build -o bin/prism ./cmd/prism
```

### File Permissions

Prism.sh creates files with secure permissions:

- Directories: `0755` (rwxr-xr-x)
- Configuration files: `0644` (rw-r--r--)
- Data files: `0644` (rw-r--r--)

### Configuration Files

Configuration files are stored in:

- **Linux:** `~/.config/prism/`
- **macOS:** `~/Library/Application Support/prism/`
- **Windows:** `%APPDATA%/prism/`

**Do not:**
- Share configuration files publicly
- Store sensitive data in palette names or descriptions
- Run prism with elevated privileges unless necessary

## Known Security Considerations

### File Operations

- **Atomic Writes:** All file writes use atomic operations (write to temp, then rename)
- **Path Traversal:** All file paths are sanitized using `filepath.Clean()`
- **File Locking:** Advisory locks prevent concurrent access issues

### Input Validation

- **Hex Colors:** Validated with regex before processing
- **File Paths:** Sanitized to prevent path traversal
- **JSON Imports:** Validated against schema before processing

### No Network Access

Prism.sh does not make network connections. All features work offline.

### Dependencies

We regularly update dependencies to patch known vulnerabilities:

```bash
# Check for vulnerabilities
go list -m all | nancy sleuth

# Update dependencies
go get -u ./...
go mod tidy
```

## Security Features

### Built-in Protections

1. **No Panics:** All errors are handled gracefully
2. **Memory Safety:** Go's memory safety prevents buffer overflows
3. **No CGO:** Pure Go implementation prevents C-related vulnerabilities
4. **No Eval:** No dynamic code execution
5. **Input Sanitization:** All user inputs are validated

### Secure Defaults

- Files created with restrictive permissions
- No auto-execution of imported content
- No automatic network requests
- No shell command execution

## Vulnerability Disclosure Policy

### Our Commitment

- We will respond to your report within 48 hours
- We will provide a fix within 30 days for critical issues
- We will credit reporters in security advisories
- We will not take legal action against researchers who:
  - Act in good faith
  - Report vulnerabilities privately
  - Don't access others' data
  - Don't degrade service

### Severity Levels

- **Critical:** Remote code execution, data loss
- **High:** Privilege escalation, information disclosure
- **Medium:** Denial of service, input validation issues
- **Low:** Information leakage, minor issues

### CVE Assignment

For critical and high-severity issues, we will:
- Request a CVE ID
- Publish a security advisory
- Coordinate disclosure timeline with reporter

## Third-Party Dependencies

We monitor our dependencies for vulnerabilities:

### Current Dependencies

```
github.com/charmbracelet/bubbletea v1.3.10
github.com/charmbracelet/lipgloss v1.1.0
```

### Dependency Policy

- Update dependencies monthly
- Review security advisories weekly
- Pin major versions to prevent breaking changes
- Test thoroughly before updating

## Incident Response

In case of a security incident:

1. **Assess:** Determine severity and impact
2. **Contain:** Prevent further exploitation
3. **Fix:** Develop and test patch
4. **Notify:** Inform affected users
5. **Release:** Deploy fix in new version
6. **Review:** Conduct post-mortem

## Security Audit History

- **v1.0.0 (November 2025):** Initial release - no known vulnerabilities

## Security Tools

We use these tools for security:

- **gosec:** Go security checker
- **nancy:** Dependency vulnerability scanner
- **CodeQL:** Static analysis (GitHub)

```bash
# Run security checks
gosec ./...
go list -m all | nancy sleuth
```

## Contact

- **Security Issues:** security@kyanite.sh (or GitHub Security Advisory)
- **General Questions:** Open a GitHub Discussion
- **Bug Reports:** Open a GitHub Issue

## Acknowledgments

We thank the security researchers who have helped make Prism.sh more secure:

- (None yet - be the first!)

---

**Last Updated:** November 2025

Thank you for helping keep Prism.sh and its users safe!
