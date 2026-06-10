# Security Policy

## Supported Versions

The `noise.sh` project follows semantic versioning. Security fixes are provided for the latest minor release and any release explicitly marked as **LTS** in the changelog or release notes.

| Version | Supported |
| ------- | --------- |
| Latest main branch | ✅ |
| Latest tagged release | ✅ |
| Older releases | ❌ |

If you are running an unsupported version, please upgrade before requesting security assistance.

## Reporting a Vulnerability

We take all security reports seriously. Please follow the process below to disclose vulnerabilities responsibly:

1. **Do not create a public issue.**
2. Email the maintainers at `security@noise.sh` (or the address configured in repository secrets).
3. Include the following details:
   - Project version and commit hash
   - Environment information (OS, architecture, build flags)
   - Step-by-step reproduction instructions
   - Impact assessment and any known mitigations

You should receive an acknowledgement within **72 hours**. We aim to provide an initial assessment within **7 calendar days** and regular updates until resolution.

If encrypted communication is required, request the GPG public key in your initial email.

## Security Response Process

1. Triage and reproduce the report.
2. Assign severity using the latest CVSS specification.
3. Develop and test a fix with regression, unit, and integration coverage.
4. Issue a patched release and update documentation, release notes, and advisories.
5. Credit the reporter (if desired) after coordinated disclosure.

## Security Testing & Tooling

Automated security scans run via GitHub Actions:

- [`security.yml`](.github/workflows/security.yml) executes:
  - Dependency vulnerability scanning with [`govulncheck`](https://go.dev/blog/vuln) (JSON + text reports)
  - Static Application Security Testing (SAST) using [`gosec`](https://github.com/securego/gosec) with SARIF upload to GitHub Code Scanning
  - License compliance verification via [`go-licenses`](https://github.com/google/go-licenses)
  - Module integrity verification using `go mod verify`

- [`ci.yml`](.github/workflows/ci.yml) runs linting, tests, cross-platform builds, and release packaging.

- [`dependabot.yml`](.github/workflows/dependabot.yml) evaluates dependency freshness and produces actionable update reports.

Security workflows run on every pull request and on a scheduled cadence to detect regressions early.

## Handling Secrets

- Never commit secrets or credentials to the repository.
- Prefer GitHub Actions secrets or external secret management for CI workflows.
- Rotate credentials immediately if exposure is suspected.

## Hardening Guidelines

- Build binaries using the provided CI workflows to benefit from code signing verification and reproducible builds.
- Use `go mod verify` after cloning to confirm module integrity.
- Run `make lint` and `go test ./...` locally before submitting contributions.
- Review third-party dependencies for license compatibility and security posture.

## Contact

- Security team email: `security@noise.sh`
- General project contact: `maintainers@noise.sh`
- Documentation: [docs/](docs/)

Thank you for helping to keep `noise.sh` secure.