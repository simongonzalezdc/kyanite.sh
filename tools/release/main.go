package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	semver "github.com/Masterminds/semver/v3"
)

var (
	typeCategoryMap = map[string]string{
		"":           "Other Changes",
		"feat":       "Features",
		"feature":    "Features",
		"fix":        "Bug Fixes",
		"bugfix":     "Bug Fixes",
		"perf":       "Performance",
		"refactor":   "Refactors",
		"docs":       "Documentation",
		"doc":        "Documentation",
		"test":       "Tests",
		"tests":      "Tests",
		"build":      "Build System",
		"ci":         "Continuous Integration",
		"chore":      "Chores",
		"style":      "Style",
		"revert":     "Reverts",
		"deps":       "Dependencies",
		"dependency": "Dependencies",
		"security":   "Security",
		"release":    "Chores",
		"infra":      "Build System",
	}
	categoryOrder = []string{
		"Features",
		"Bug Fixes",
		"Security",
		"Performance",
		"Refactors",
		"Documentation",
		"Tests",
		"Build System",
		"Continuous Integration",
		"Dependencies",
		"Chores",
		"Style",
		"Reverts",
		"Other Changes",
	}
)

const changelogHeader = "# Changelog"

type commitInfo struct {
	Hash         string
	Author       string
	Email        string
	Subject      string
	Summary      string
	Category     string
	Breaking     bool
	BreakingNote string
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "verify-tag":
		err = runVerifyTag(os.Args[2:])
	case "notes", "generate-notes":
		err = runGenerateNotes(os.Args[2:])
	case "changelog":
		err = runUpdateChangelog(os.Args[2:])
	case "checksums":
		err = runGenerateChecksums(os.Args[2:])
	case "verify-checksums":
		err = runVerifyChecksums(os.Args[2:])
	default:
		usage()
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `release tool

Usage:
  release <command> [options]

Commands:
  verify-tag        Validate semantic version tags
  notes             Generate release notes from git commits
  changelog         Update or generate CHANGELOG.md content
  checksums         Generate SHA-256 checksums for release assets
  verify-checksums  Verify previously generated checksums
`)
}

func runVerifyTag(args []string) error {
	fs := flag.NewFlagSet("verify-tag", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var tag string
	var previous string
	fs.StringVar(&tag, "tag", "", "Release tag to validate (e.g. v1.2.3)")
	fs.StringVar(&previous, "previous", "", "Previous release tag (optional)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if tag == "" {
		return errors.New("--tag is required")
	}
	if !strings.HasPrefix(tag, "v") {
		return fmt.Errorf("release tag must be prefixed with 'v': %s", tag)
	}
	currentVersion, err := semver.StrictNewVersion(stripTagPrefix(tag))
	if err != nil {
		return fmt.Errorf("invalid tag %q: %w", tag, err)
	}

	if previous != "" {
		if !strings.HasPrefix(previous, "v") {
			return fmt.Errorf("previous tag must be prefixed with 'v': %s", previous)
		}
		prevVersion, err := semver.StrictNewVersion(stripTagPrefix(previous))
		if err != nil {
			return fmt.Errorf("invalid previous tag %q: %w", previous, err)
		}
		if !currentVersion.GreaterThan(prevVersion) {
			return fmt.Errorf("tag %s must be greater than previous release %s", tag, previous)
		}
	}

	fmt.Printf("Semantic version verified: %s\n", tag)
	return nil
}

func runGenerateNotes(args []string) error {
	fs := flag.NewFlagSet("notes", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var from string
	var to string
	var version string
	var output string
	fs.StringVar(&from, "from", "", "Previous tag or commit (exclusive)")
	fs.StringVar(&to, "to", "", "Target tag or commit (inclusive)")
	fs.StringVar(&version, "version", "", "Release version tag (e.g. v1.2.3)")
	fs.StringVar(&output, "output", "", "Path to write the release notes markdown")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if version == "" {
		version = to
	}
	if version == "" {
		return errors.New("either --version or --to must be provided")
	}
	if output == "" {
		return errors.New("--output is required")
	}
	if to == "" {
		to = version
	}

	commits, err := collectCommits(from, to)
	if err != nil {
		return err
	}

	repoURL, err := deriveRepositoryURL()
	if err != nil {
		return err
	}

	content := generateReleaseNotesContent(version, from, to, commits, repoURL)

	if err := ensureDir(filepath.Dir(output)); err != nil {
		return err
	}
	if err := os.WriteFile(output, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write release notes: %w", err)
	}

	fmt.Printf("Release notes written to %s (%d commits)\n", output, len(commits))
	return nil
}

func runUpdateChangelog(args []string) error {
	fs := flag.NewFlagSet("changelog", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var version string
	var notesPath string
	var changelogPath string
	var outputPath string
	fs.StringVar(&version, "version", "", "Release version (e.g. v1.2.3)")
	fs.StringVar(&notesPath, "notes", "", "Path to release notes markdown")
	fs.StringVar(&changelogPath, "changelog", "CHANGELOG.md", "Existing changelog file path")
	fs.StringVar(&outputPath, "output", "", "Output path for updated changelog (defaults to --changelog)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if version == "" {
		return errors.New("--version is required")
	}
	if notesPath == "" {
		return errors.New("--notes is required")
	}
	if outputPath == "" {
		outputPath = changelogPath
	}

	notesBytes, err := os.ReadFile(notesPath)
	if err != nil {
		return fmt.Errorf("read notes: %w", err)
	}
	entry := strings.TrimSpace(string(notesBytes))
	if entry == "" {
		return errors.New("release notes content is empty")
	}

	entryHeader := fmt.Sprintf("## %s", version)
	entryBracketHeader := fmt.Sprintf("## [%s]", version)
	if strings.Contains(entry, entryHeader) {
		entry = strings.Replace(entry, entryHeader, entryBracketHeader, 1)
	} else if !strings.Contains(entry, entryBracketHeader) {
		entry = entryBracketHeader + "\n\n" + entry
	}

	var existing string
	originalBytes, err := os.ReadFile(changelogPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("read changelog: %w", err)
		}
	} else {
		existing = string(originalBytes)
	}

	updated, err := composeChangelog(existing, entry, version)
	if err != nil {
		return err
	}

	if err := ensureDir(filepath.Dir(outputPath)); err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, []byte(updated), 0o600); err != nil {
		return fmt.Errorf("write changelog: %w", err)
	}

	fmt.Printf("Changelog updated at %s\n", outputPath)
	return nil
}

func runGenerateChecksums(args []string) error {
	fs := flag.NewFlagSet("checksums", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inputDir string
	var output string
	var algorithm string
	fs.StringVar(&inputDir, "input-dir", ".", "Directory containing release assets")
	fs.StringVar(&output, "output", "", "Path to write checksum file")
	fs.StringVar(&algorithm, "algorithm", "sha256", "Checksum algorithm (sha256 only)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if output == "" {
		return errors.New("--output is required")
	}
	if algorithm != "sha256" {
		return fmt.Errorf("unsupported algorithm %q (only sha256 is supported)", algorithm)
	}

	files := []string{}
	if err := filepath.WalkDir(inputDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Clean(path) == filepath.Clean(output) {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return fmt.Errorf("walk input dir: %w", err)
	}

	sort.Strings(files)

	var builder strings.Builder
	for _, file := range files {
		hash, err := hashFile(file)
		if err != nil {
			return fmt.Errorf("hash %s: %w", file, err)
		}
		rel, err := filepath.Rel(inputDir, file)
		if err != nil {
			return fmt.Errorf("relative path for %s: %w", file, err)
		}
		builder.WriteString(fmt.Sprintf("%s  %s\n", hash, filepath.ToSlash(rel)))
	}

	if err := ensureDir(filepath.Dir(output)); err != nil {
		return err
	}
	if err := os.WriteFile(output, []byte(builder.String()), 0o600); err != nil {
		return fmt.Errorf("write checksum file: %w", err)
	}

	fmt.Printf("Generated %d checksums at %s\n", len(files), output)
	return nil
}

func runVerifyChecksums(args []string) error {
	fs := flag.NewFlagSet("verify-checksums", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var inputDir string
	var checksumPath string
	fs.StringVar(&inputDir, "input-dir", ".", "Directory containing release assets")
	fs.StringVar(&checksumPath, "checksums", "", "Path to checksum file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if checksumPath == "" {
		return errors.New("--checksums is required")
	}

	file, err := os.Open(checksumPath)
	if err != nil {
		return fmt.Errorf("open checksums: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	verified := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return fmt.Errorf("invalid checksum format on line %d: %s", lineNo, line)
		}
		hash := fields[0]
		rel := fields[len(fields)-1]
		target := filepath.Join(inputDir, filepath.FromSlash(rel))

		expectedHash := strings.ToLower(hash)
		actualHash, err := hashFile(target)
		if err != nil {
			return fmt.Errorf("hash %s: %w", target, err)
		}
		if actualHash != expectedHash {
			return fmt.Errorf("checksum mismatch for %s: expected %s got %s", rel, expectedHash, actualHash)
		}
		verified++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read checksums: %w", err)
	}

	fmt.Printf("Verified %d checksum entries\n", verified)
	return nil
}

func collectCommits(from, to string) ([]commitInfo, error) {
	if to == "" {
		return nil, errors.New("--to is required")
	}

	const logFormat = "%H%x1f%an%x1f%ae%x1f%s%x1f%B%x1e"
	args := []string{"log", "--no-merges", "--reverse", fmt.Sprintf("--pretty=format:%s", logFormat)}
	if from == "" {
		args = append(args, to)
	} else {
		args = append(args, fmt.Sprintf("%s..%s", from, to))
	}

	cmd := exec.Command("git", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log range: %w", err)
	}

	records := bytes.Split(output, []byte{0x1e})
	commits := make([]commitInfo, 0, len(records))

	for _, record := range records {
		record = bytes.TrimSpace(record)
		if len(record) == 0 {
			continue
		}

		fields := bytes.SplitN(record, []byte{0x1f}, 5)
		if len(fields) != 5 {
			return nil, fmt.Errorf("unexpected git log record format: %q", string(record))
		}

		hash := string(fields[0])
		author := strings.TrimSpace(string(fields[1]))
		email := strings.TrimSpace(string(fields[2]))
		subject := strings.TrimSpace(string(fields[3]))
		fullBody := string(fields[4])

		body := extractBody(fullBody)

		category, summary, breaking, breakingNote := analyzeCommit(subject, body)

		commits = append(commits, commitInfo{
			Hash:         hash,
			Author:       author,
			Email:        email,
			Subject:      subject,
			Summary:      summary,
			Category:     category,
			Breaking:     breaking,
			BreakingNote: breakingNote,
		})
	}

	return commits, nil
}

func analyzeCommit(subject, body string) (string, string, bool, string) {
	trimmed := strings.TrimSpace(subject)
	if trimmed == "" {
		return "Other Changes", "Miscellaneous updates", false, ""
	}

	header := trimmed
	summary := trimmed
	scope := ""
	breaking := false
	breakingNote := ""

	if idx := strings.Index(header, ":"); idx != -1 {
		summary = strings.TrimSpace(header[idx+1:])
		header = header[:idx]
	}

	if idx := strings.Index(header, "("); idx != -1 {
		if end := strings.Index(header[idx:], ")"); end != -1 {
			scope = strings.TrimSpace(header[idx+1 : idx+end])
			header = header[:idx] + header[idx+end+1:]
		}
	}

	if strings.Contains(header, "!") {
		breaking = true
	}
	header = strings.TrimSuffix(header, "!")
	commitType := strings.ToLower(strings.TrimSpace(header))
	if commitType == "" {
		commitType = ""
	}

	category := mapCategory(commitType)

	if scope != "" {
		if summary == "" {
			summary = fmt.Sprintf("Updates for %s", scope)
		} else {
			summary = fmt.Sprintf("%s (%s)", summary, scope)
		}
	}

	if summary == "" {
		summary = trimmed
	}

	summary = ensureSentence(summary)

	if !breaking && body != "" {
		for _, line := range strings.Split(body, "\n") {
			upper := strings.ToUpper(strings.TrimSpace(line))
			if strings.HasPrefix(upper, "BREAKING CHANGE") || strings.HasPrefix(upper, "BREAKING-CHANGE") {
				breaking = true
				if idx := strings.Index(line, ":"); idx != -1 && len(line) > idx+1 {
					breakingNote = ensureSentence(strings.TrimSpace(line[idx+1:]))
				}
			}
		}
	}

	if breaking {
		if breakingNote == "" {
			breakingNote = summary
		}
	}

	return category, summary, breaking, breakingNote
}

func extractBody(fullBody string) string {
	if fullBody == "" {
		return ""
	}
	normalized := strings.ReplaceAll(fullBody, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	if len(lines) <= 1 {
		return ""
	}
	rest := strings.Join(lines[1:], "\n")
	return strings.TrimSpace(rest)
}

func generateReleaseNotesContent(version, from, to string, commits []commitInfo, repoURL string) string {
	var builder strings.Builder
	date := time.Now().UTC().Format("2006-01-02")

	builder.WriteString(fmt.Sprintf("## %s - %s\n\n", version, date))

	if from != "" {
		builder.WriteString(fmt.Sprintf("Compare: %s/compare/%s...%s\n\n", repoURL, from, to))
	}

	authors := make(map[string]struct{})
	commitsByCategory := make(map[string][]commitInfo)
	breakingChanges := make([]commitInfo, 0)

	for _, c := range commits {
		category := c.Category
		if category == "" {
			category = "Other Changes"
		}
		commitsByCategory[category] = append(commitsByCategory[category], c)
		if c.Breaking {
			breakingChanges = append(breakingChanges, c)
		}
		name := strings.TrimSpace(c.Author)
		if name == "" {
			name = c.Email
		}
		if name != "" {
			authors[name] = struct{}{}
		}
	}

	builder.WriteString("### Summary\n")
	builder.WriteString(fmt.Sprintf("- %d commits\n", len(commits)))
	builder.WriteString(fmt.Sprintf("- %d contributors\n\n", len(authors)))

	if len(commits) == 0 {
		builder.WriteString("_No commits were found for this release range._\n\n")
	}

	if len(breakingChanges) > 0 {
		builder.WriteString("### Breaking Changes\n")
		for _, c := range breakingChanges {
			note := c.BreakingNote
			if note == "" {
				note = c.Summary
			}
			builder.WriteString(fmt.Sprintf("- %s [`%s`](%s/commit/%s) — %s\n",
				note, shortHash(c.Hash), repoURL, c.Hash, c.Author))
		}
		builder.WriteString("\n")
	}

	for _, category := range categoryOrder {
		entries := commitsByCategory[category]
		if len(entries) == 0 {
			continue
		}
		builder.WriteString(fmt.Sprintf("### %s\n", category))
		for _, c := range entries {
			builder.WriteString(fmt.Sprintf("- %s [`%s`](%s/commit/%s) — %s\n",
				c.Summary, shortHash(c.Hash), repoURL, c.Hash, c.Author))
		}
		builder.WriteString("\n")
	}

	if len(authors) > 0 {
		builder.WriteString("### Contributors\n")
		names := make([]string, 0, len(authors))
		for name := range authors {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			builder.WriteString(fmt.Sprintf("- %s\n", name))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("### Verification\n")
	builder.WriteString("Download the release asset `checksums.txt` and run `sha256sum --check checksums.txt` (Linux) or `shasum -a 256 --check checksums.txt` (macOS) to verify the binary integrity.\n")

	content := strings.TrimRight(builder.String(), "\n") + "\n"
	return content
}

func composeChangelog(existing, entry, version string) (string, error) {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return "", errors.New("changelog entry is empty")
	}
	if !strings.HasPrefix(entry, "## [") {
		entry = fmt.Sprintf("## [%s]\n\n%s", version, entry)
	}

	entry = strings.TrimRight(entry, "\n") + "\n"

	if strings.TrimSpace(existing) == "" {
		return fmt.Sprintf("%s\n\n%s\n", changelogHeader, entry), nil
	}

	existing = strings.TrimLeft(existing, "\ufeff")
	if !strings.HasPrefix(strings.TrimSpace(existing), changelogHeader) {
		return "", fmt.Errorf("existing changelog must start with %q", changelogHeader)
	}

	body := strings.TrimSpace(existing[len(changelogHeader):])
	note, entries := splitChangelogBody(body)

	if strings.Contains(entries, fmt.Sprintf("## [%s]", version)) {
		return "", fmt.Errorf("changelog already contains entry for %s", version)
	}

	var builder strings.Builder
	builder.WriteString(changelogHeader)
	builder.WriteString("\n\n")

	if note != "" {
		builder.WriteString(note)
		builder.WriteString("\n\n")
	}

	builder.WriteString(entry)
	if entries != "" {
		builder.WriteString("\n")
		builder.WriteString(strings.TrimLeft(entries, "\n"))
		builder.WriteString("\n")
	} else {
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

func splitChangelogBody(body string) (string, string) {
	body = strings.TrimLeft(body, "\n")
	if body == "" {
		return "", ""
	}

	idx := strings.Index(body, "\n## ")
	if idx == -1 {
		trimmed := strings.TrimSpace(body)
		if strings.HasPrefix(trimmed, "## ") {
			return "", trimmed
		}
		return trimmed, ""
	}

	note := strings.TrimSpace(body[:idx])
	entries := strings.TrimSpace(body[idx+1:])
	return note, entries
}

func mapCategory(commitType string) string {
	if cat, ok := typeCategoryMap[commitType]; ok {
		return cat
	}
	return "Other Changes"
}

func ensureSentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size == 0 {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

func shortHash(hash string) string {
	if len(hash) <= 7 {
		return hash
	}
	return hash[:7]
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func ensureDir(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o750)
}

func deriveRepositoryURL() (string, error) {
	remoteCmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := remoteCmd.Output()
	if err != nil {
		if repo := os.Getenv("GITHUB_REPOSITORY"); repo != "" {
			return fmt.Sprintf("https://github.com/%s", strings.TrimSpace(repo)), nil
		}
		return "", fmt.Errorf("determine git remote: %w", err)
	}

	return normalizeRemoteURL(strings.TrimSpace(string(output)))
}

func normalizeRemoteURL(remote string) (string, error) {
	if remote == "" {
		return "", errors.New("empty remote URL")
	}

	if strings.HasPrefix(remote, "git@") {
		trimmed := strings.TrimSuffix(strings.TrimPrefix(remote, "git@"), ".git")
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("unable to parse remote %q", remote)
		}
		return fmt.Sprintf("https://%s/%s", parts[0], strings.TrimSuffix(parts[1], ".git")), nil
	}

	if strings.HasPrefix(remote, "https://") || strings.HasPrefix(remote, "http://") {
		trimmed := strings.TrimSuffix(remote, ".git")
		return trimmed, nil
	}

	return "", fmt.Errorf("unsupported remote format %q", remote)
}

func stripTagPrefix(tag string) string {
	return strings.TrimPrefix(tag, "v")
}
