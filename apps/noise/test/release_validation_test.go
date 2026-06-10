package noise

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	semver "github.com/Masterminds/semver/v3"
)

func TestReleaseVersionInjection(t *testing.T) {
	releaseVersion := os.Getenv("RELEASE_VERSION")
	if releaseVersion == "" {
		t.Skip("RELEASE_VERSION not set; skipping release validation test")
	}

	stripped := strings.TrimPrefix(strings.TrimSpace(releaseVersion), "v")
	if stripped == "" {
		t.Fatalf("empty release version after trimming prefix: %q", releaseVersion)
	}

	if _, err := semver.StrictNewVersion(stripped); err != nil {
		t.Fatalf("release version %q is not a valid semantic version: %v", releaseVersion, err)
	}

	buildDate := "1970-01-01"
	buildCommit := "validation"

	cmd := exec.Command("go", "run",
		"-ldflags", fmt.Sprintf("-X main.version=%s -X main.commit=%s -X main.date=%s", stripped, buildCommit, buildDate),
		"./cmd/noise",
		"--version",
	)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to run noise --version: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, stripped) {
		t.Fatalf("version output %q does not contain version %q", output, stripped)
	}
	if !strings.Contains(output, buildCommit) {
		t.Fatalf("version output %q does not contain commit %q", output, buildCommit)
	}
	if !strings.Contains(output, buildDate) {
		t.Fatalf("version output %q does not contain build date %q", output, buildDate)
	}
}

func TestReleaseTagVerification(t *testing.T) {
	releaseVersion := os.Getenv("RELEASE_VERSION")
	if releaseVersion == "" {
		t.Skip("RELEASE_VERSION not set; skipping release validation test")
	}

	cmd := exec.Command("go", "run", "./tools/release/main.go", "verify-tag", "--tag", releaseVersion)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("release tag verification failed: %v", err)
	}
}
