# Security Policy

**syntax.sh Security Policy**
**Version:** 1.0
**Last Updated:** November 2025

---

## Reporting Security Vulnerabilities

### How to Report

If you discover a security vulnerability in syntax.sh, please report it responsibly:

**DO NOT** open a public GitHub issue for security vulnerabilities.

**Instead, please email:**
- **Email:** security@kyanite.sh
- **Subject:** `[SECURITY] syntax.sh - Brief Description`

**Include in your report:**
1. Description of the vulnerability
2. Steps to reproduce
3. Potential impact
4. Suggested fix (if any)
5. Your name/handle (for acknowledgment)

### Response Timeline

- **Acknowledgment:** Within 48 hours
- **Initial Assessment:** Within 5 business days
- **Fix Timeline:** Depends on severity (see below)

### Severity Levels

| Severity | Description | Fix Timeline |
|----------|-------------|--------------|
| **Critical** | Remote code execution, data loss | 24-48 hours |
| **High** | Data exposure, auth bypass | 1 week |
| **Medium** | DoS, information disclosure | 2 weeks |
| **Low** | Minor issues, best practices | Next release |

---

## Supported Versions

| Version | Supported | Security Updates |
|---------|-----------|------------------|
| 1.0.x   | ✅ Yes    | Active |
| < 1.0   | ❌ No     | Upgrade required |

**Policy:**
- Current major version (v1.x) receives all security updates
- Previous major versions receive critical security updates for 6 months after new major release
- Older versions are not supported

---

## Security Features

### 1. Local-First Architecture

**Design Philosophy:**
- All data stored locally on user's machine
- No cloud sync (eliminates server-side attacks)
- No telemetry or data collection
- Full user control over data

**Benefits:**
- No remote attack surface
- No data breaches from centralized servers
- Privacy by design

### 2. File Permissions

**Configuration Files:**
```
~/.config/syntax/config.toml    0600 (owner read/write only)
~/.config/syntax/api_keys.enc   0600 (if API mode enabled)
```

**Data Directories:**
```
~/.local/share/syntax/          0700 (owner only)
~/.local/share/syntax/projects/ 0700 (owner only)
```

**Rationale:**
- Prevents other users on shared systems from accessing data
- Protects API keys from unauthorized access

**Implementation:**
```go
// Always set restrictive permissions
os.MkdirAll(configDir, 0700)
os.WriteFile(configFile, data, 0600)
```

### 3. Input Sanitization

**File Paths:**

All user-provided file paths are sanitized to prevent:
- Path traversal attacks (`../../../etc/passwd`)
- Directory escaping
- Symbolic link attacks

```go
func SanitizeFilename(input string) string {
    // Remove path separators
    safe := strings.ReplaceAll(input, "/", "-")
    safe = strings.ReplaceAll(safe, "\\", "-")
    safe = strings.ReplaceAll(safe, "..", "")

    // Remove special characters
    safe = strings.Map(func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsNumber(r) ||
           r == '-' || r == '_' || r == ' ' {
            return r
        }
        return -1
    }, safe)

    // Limit length (filesystem constraint)
    if len(safe) > 255 {
        safe = safe[:255]
    }

    return filepath.Clean(safe)
}

// Validation example
func ValidateProjectPath(path string) error {
    cleaned := filepath.Clean(path)

    // Reject path traversal
    if strings.Contains(cleaned, "..") {
        return errors.New("path traversal detected")
    }

    // Ensure within allowed directory
    if !strings.HasPrefix(cleaned, getAllowedDir()) {
        return errors.New("path outside allowed directory")
    }

    return nil
}
```

**User Input:**

Character names, scene titles, etc. are sanitized:
```go
func SanitizeUserInput(input string, maxLen int) string {
    // Trim whitespace
    input = strings.TrimSpace(input)

    // Limit length (prevent DoS)
    if len(input) > maxLen {
        input = input[:maxLen]
    }

    // Remove null bytes (can cause issues in C libraries)
    input = strings.ReplaceAll(input, "\x00", "")

    // Remove control characters
    input = strings.Map(func(r rune) rune {
        if unicode.IsControl(r) && r != '\n' && r != '\t' {
            return -1
        }
        return r
    }, input)

    return input
}
```

### 4. API Key Storage

**Problem:** AI assistant requires API keys (if using external API)

**Solution:** OS-native secure storage

**Primary Method:** System Keyring
```go
import "github.com/zalando/go-keyring"

func StoreAPIKey(provider, key string) error {
    service := "syntax.sh"
    username := provider  // "openrouter", "openai", etc.

    return keyring.Set(service, username, key)
}

func GetAPIKey(provider string) (string, error) {
    service := "syntax.sh"
    username := provider

    return keyring.Get(service, username)
}
```

**Platform Implementation:**
- **macOS:** Keychain
- **Windows:** Credential Manager
- **Linux:** Secret Service (GNOME Keyring, KWallet)

**Fallback Method:** Encrypted file (if keyring unavailable)
```go
func encryptAPIKey(key string) ([]byte, error) {
    // Use AES-256-GCM for encryption
    // Key derived from machine-specific data + user password
    // Store in ~/.config/syntax/api_keys.enc with 0600 permissions

    // IMPORTANT: Show warning to user
    fmt.Println("WARNING: API key stored in encrypted file.")
    fmt.Println("For better security, use OS keyring.")

    return encryptedData, nil
}
```

**Never:** Plain text storage

### 5. Denial of Service Prevention

**Resource Limits:**

```go
const (
    MaxFileSize      = 50 * 1024 * 1024  // 50MB per file
    MaxCharacters    = 1000              // Per project
    MaxScenes        = 10000             // Per project
    MaxLocations     = 1000              // Per project
    MaxSearchResults = 100               // Limit search output
    MaxUndo          = 100               // Undo history limit
)

func LoadScene(path string) (*Scene, error) {
    // Check file size before loading
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if info.Size() > MaxFileSize {
        return nil, errors.New("file too large")
    }

    // Load with size limit
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    return ParseScene(data)
}
```

**Rationale:**
- Prevent accidental memory exhaustion
- Protect against malicious crafted files
- Ensure responsive UI

### 6. Dependency Security

**Vulnerability Scanning:**

Required before each release:
```bash
# Official Go vulnerability scanner
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Check dependency tree
go list -json -m all | nancy sleuth
```

**Dependency Policy:**
- Minimal dependencies (reduce attack surface)
- Only well-maintained packages
- Regular security audits
- Immediate updates for CVEs

**Current Dependencies (v1.0):**
```
github.com/charmbracelet/bubbletea  # TUI framework
github.com/charmbracelet/lipgloss   # Styling
github.com/adrg/frontmatter         # YAML parsing
github.com/adrg/xdg                 # Cross-platform paths
github.com/zalando/go-keyring       # Secure key storage
github.com/jung-kurt/gofpdf         # PDF export
github.com/nguyenthenguyen/docx     # DOCX export
```

All dependencies reviewed for:
- ✅ Active maintenance
- ✅ No known CVEs
- ✅ Minimal sub-dependencies
- ✅ Clear license (MIT/Apache 2.0)

### 7. Secure File Operations

**Atomic Writes:**

Prevent data corruption/loss during saves:
```go
func AtomicWriteFile(path string, data []byte, perm os.FileMode) error {
    // Write to temporary file
    tmpPath := path + ".tmp"
    if err := os.WriteFile(tmpPath, data, perm); err != nil {
        return err
    }

    // Sync to disk (ensure durability)
    if err := syncFile(tmpPath); err != nil {
        os.Remove(tmpPath)
        return err
    }

    // Atomic rename (POSIX guarantees atomicity)
    if err := os.Rename(tmpPath, path); err != nil {
        os.Remove(tmpPath)
        return err
    }

    return nil
}

func syncFile(path string) error {
    f, err := os.OpenFile(path, os.O_RDWR, 0)
    if err != nil {
        return err
    }
    defer f.Close()

    return f.Sync()
}
```

**Benefits:**
- No partial writes visible to readers
- Power loss during write doesn't corrupt data
- Original file unchanged if write fails

---

## Threat Model

### In-Scope Threats

1. **Local File Manipulation**
   - Mitigation: File permissions (0600/0700)
   - Mitigation: Input validation

2. **Path Traversal Attacks**
   - Mitigation: Path sanitization
   - Mitigation: Whitelist allowed directories

3. **Malicious Project Files**
   - Mitigation: Schema validation
   - Mitigation: Size limits
   - Mitigation: Parse errors handled gracefully

4. **API Key Exposure**
   - Mitigation: OS keyring storage
   - Mitigation: Never log keys
   - Mitigation: Memory cleared after use

5. **Dependency Vulnerabilities**
   - Mitigation: Regular scans
   - Mitigation: Minimal dependencies
   - Mitigation: Rapid update process

### Out-of-Scope Threats

1. **Physical Access to Machine**
   - If attacker has physical access, no software can protect data
   - Recommendation: Use full-disk encryption (OS-level)

2. **Compromised OS**
   - If OS is compromised, all applications are at risk
   - Recommendation: Keep OS updated, use antivirus

3. **Supply Chain Attacks**
   - We use Go's module verification
   - All dependencies pinned in go.sum
   - Regular review of dependency updates

4. **Network Attacks**
   - No network functionality except optional AI API
   - AI API uses HTTPS (TLS) for all connections
   - User warned when API mode enabled

---

## Security Best Practices for Users

### 1. Keep Software Updated

```bash
# Check for updates
syntax --version

# Update to latest from the public repository
git clone https://github.com/simongonzalezdc/Syntax.sh.git
cd Syntax.sh
go build -o bin/syntax ./cmd/syntax
```

### 2. Use Local Mode (Default)

- AI assistant in local mode (Ollama) by default
- No data leaves your machine
- No API keys required

### 3. Protect API Keys

If using external AI API:
- Never share API keys
- Revoke keys if compromised
- Use read-only or limited-scope keys when possible

### 4. Regular Backups

```bash
# Manual backup
syntax backup

# Backups stored in: ~/.local/share/syntax/projects/{project}/.backups/
```

### 5. Secure Your System

- Use full-disk encryption (FileVault, BitLocker, LUKS)
- Keep OS and software updated
- Use strong user account password
- Enable firewall

---

## Security Checklist for Developers

**Before Release:**

- [ ] Run `govulncheck ./...` (no vulnerabilities)
- [ ] Review all dependencies for updates
- [ ] Test file permission enforcement
- [ ] Test input sanitization (fuzz testing)
- [ ] Verify atomic file writes work correctly
- [ ] Check no sensitive data in logs
- [ ] Verify API keys not in memory dumps
- [ ] Test on all platforms (Linux, macOS, Windows)
- [ ] Security audit of new code
- [ ] Update this document if threat model changes

**Code Review Checklist:**

- [ ] No hardcoded credentials
- [ ] All user input sanitized
- [ ] File operations use atomic writes
- [ ] Error messages don't leak sensitive info
- [ ] No arbitrary code execution
- [ ] Resource limits enforced
- [ ] Dependencies up to date

---

## Incident Response Plan

### If Vulnerability Discovered

1. **Assess Severity**
   - Critical/High: Immediate action
   - Medium/Low: Schedule for next release

2. **Develop Fix**
   - Create private branch
   - Develop and test fix
   - Peer review

3. **Coordinate Disclosure**
   - Notify reporter
   - Prepare security advisory
   - Plan coordinated release

4. **Release Patch**
   - Tag new version (e.g., v1.0.1)
   - Publish GitHub release
   - Update security advisory
   - Notify users (GitHub, docs)

5. **Post-Mortem**
   - Document incident
   - Update security practices
   - Improve testing/review process

### Example Timeline (Critical Vulnerability)

- **Hour 0:** Vulnerability reported
- **Hour 2:** Acknowledged, assessment begins
- **Hour 6:** Severity confirmed (critical)
- **Hour 12:** Fix developed
- **Hour 18:** Fix tested
- **Hour 24:** Patch released
- **Hour 48:** Security advisory published

---

## Compliance

### Data Privacy

**GDPR Compliance:**
- ✅ No user data collected
- ✅ No analytics or telemetry
- ✅ All data stored locally
- ✅ User has full control

**CCPA Compliance:**
- ✅ No personal data processing
- ✅ No data sales
- ✅ Local-only storage

### License Compliance

**Dependencies:**
- All dependencies use permissive licenses (MIT, Apache 2.0)
- License compatibility verified
- Attribution provided in LICENSES.md

---

## Security Contacts

**Security Issues:**
- Email: security@kyanite.sh
- Response Time: 48 hours

**General Questions:**
- GitHub Discussions: github.com/simongonzalezdc/Syntax.sh/discussions
- Documentation: docs.kyanite.sh

---

## Acknowledgments

We thank the security research community for responsible disclosure practices. Security researchers will be acknowledged (with permission) in:
- Security advisories
- Release notes
- SECURITY.md (this file)

### Hall of Fame

*(None yet - be the first!)*

---

## Changelog

### v1.0 (November 2025)
- Initial security policy
- Defined threat model
- Established disclosure process
- Documented security features

---

**Last Reviewed:** November 2025
**Next Review:** March 2026 (or after any security incident)
