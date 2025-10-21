package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// PlatformConfig represents a target platform configuration
type PlatformConfig struct {
	GOOS    string
	GOARCH  string
	Ext     string
	Archive string
}

// BuildResult represents the result of a build test
type BuildResult struct {
	Platform string
	Success  bool
	Duration time.Duration
	Error    string
	BinarySize int64
}

var platforms = []PlatformConfig{
	// Linux platforms
	{"linux", "amd64", "", "tar.gz"},
	{"linux", "arm64", "", "tar.gz"},
	{"linux", "arm", "", "tar.gz"},
	{"linux", "386", "", "tar.gz"},

	// macOS platforms (Darwin)
	{"darwin", "amd64", "", "tar.gz"},
	{"darwin", "arm64", "", "tar.gz"},

	// Windows platforms
	{"windows", "amd64", ".exe", "zip"},
	{"windows", "arm64", ".exe", "zip"},
	{"windows", "386", ".exe", "zip"},
}

func main() {
	fmt.Println("🔍 Cross-Platform Build Verification")
	fmt.Println("===================================")
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Platform: %s\n", runtime.GOOS)
	fmt.Println()

	results := make([]BuildResult, 0, len(platforms))

	for _, platform := range platforms {
		fmt.Printf("Building for %s/%s...\n", platform.GOOS, platform.GOARCH)

		result := testBuild(platform)
		results = append(results, result)

		if result.Success {
			fmt.Printf("  ✅ Success (%.2fs, %s)\n",
				result.Duration.Seconds(),
				formatSize(result.BinarySize))
		} else {
			fmt.Printf("  ❌ Failed: %s\n", result.Error)
		}
		fmt.Println()
	}

	printSummary(results)
}

func testBuild(platform PlatformConfig) BuildResult {
	start := time.Now()

	// Create output filename
	binaryName := fmt.Sprintf("noise-%s-%s%s", platform.GOOS, platform.GOARCH, platform.Ext)
	buildDir := "build/verify"
	binaryPath := filepath.Join(buildDir, binaryName)

	// Ensure build directory exists
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return BuildResult{
			Platform: fmt.Sprintf("%s/%s", platform.GOOS, platform.GOARCH),
			Success:  false,
			Duration: time.Since(start),
			Error:    fmt.Sprintf("failed to create build dir: %v", err),
		}
	}

	// Clean previous build
	os.Remove(binaryPath)

	// Build command
	args := []string{
		"build",
		"-trimpath",
		"-ldflags=-s -w",
		"-o", binaryPath,
	}

	// Add cross-compilation flags if different from current platform
	if platform.GOOS != runtime.GOOS || platform.GOARCH != runtime.GOARCH {
		args = append([]string{"env", fmt.Sprintf("GOOS=%s", platform.GOOS), fmt.Sprintf("GOARCH=%s", platform.GOARCH), "CGO_ENABLED=0"}, append([]string{"go"}, args...)...)
	} else {
		args = append([]string{"go"}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		return BuildResult{
			Platform: fmt.Sprintf("%s/%s", platform.GOOS, platform.GOARCH),
			Success:  false,
			Duration: duration,
			Error:    err.Error(),
		}
	}

	// Get binary size
	var binarySize int64
	if info, err := os.Stat(binaryPath); err == nil {
		binarySize = info.Size()
	}

	// Clean up
	os.Remove(binaryPath)
	os.Remove(buildDir)

	return BuildResult{
		Platform:   fmt.Sprintf("%s/%s", platform.GOOS, platform.GOARCH),
		Success:    true,
		Duration:   duration,
		BinarySize: binarySize,
	}
}

func printSummary(results []BuildResult) {
	fmt.Println("📊 Build Verification Summary")
	fmt.Println("============================")

	successful := 0
	failed := 0
	totalTime := time.Duration(0)
	var fastest, slowest BuildResult
	var largest, smallest BuildResult

	for _, result := range results {
		if result.Success {
			successful++
			totalTime += result.Duration

			if fastest.Duration == 0 || result.Duration < fastest.Duration {
				fastest = result
			}
			if result.Duration > slowest.Duration {
				slowest = result
			}

			if largest.BinarySize == 0 || result.BinarySize > largest.BinarySize {
				largest = result
			}
			if smallest.BinarySize == 0 || result.BinarySize < smallest.BinarySize {
				smallest = result
			}
		} else {
			failed++
		}
	}

	// Overall status
	fmt.Printf("✅ Successful: %d\n", successful)
	fmt.Printf("❌ Failed: %d\n", failed)
	fmt.Printf("📈 Success Rate: %.1f%%\n", float64(successful)/float64(len(results))*100)
	fmt.Println()

	if successful > 0 {
		fmt.Printf("⏱️  Average Build Time: %.2fs\n", totalTime.Seconds()/float64(successful))
		fmt.Printf("🚀 Fastest Build: %s (%.2fs)\n", fastest.Platform, fastest.Duration.Seconds())
		fmt.Printf("🐌 Slowest Build: %s (%.2fs)\n", slowest.Platform, slowest.Duration.Seconds())
		fmt.Println()

		fmt.Printf("📦 Largest Binary: %s (%s)\n", largest.Platform, formatSize(largest.BinarySize))
		fmt.Printf("📦 Smallest Binary: %s (%s)\n", smallest.Platform, formatSize(smallest.BinarySize))
		fmt.Println()
	}

	// Platform details
	fmt.Println("📋 Platform Details")
	fmt.Println("------------------")
	for _, result := range results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}

		fmt.Printf("%s %s: %.2fs", status, result.Platform, result.Duration.Seconds())
		if result.Success {
			fmt.Printf(" (%s)", formatSize(result.BinarySize))
		} else {
			fmt.Printf(" - %s", result.Error)
		}
		fmt.Println()
	}

	fmt.Println()

	// Recommendations
	if failed == 0 {
		fmt.Println("🎉 All platforms built successfully!")
		fmt.Println("✅ Week 8 cross-platform build automation is working correctly.")
	} else {
		fmt.Printf("⚠️  %d platform(s) failed to build. Check the errors above.\n", failed)
		fmt.Println("💡 Consider checking:")
		fmt.Println("   - Go version compatibility")
		fmt.Println("   - CGO dependencies for cross-compilation")
		fmt.Println("   - Platform-specific build requirements")
	}
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}