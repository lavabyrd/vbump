package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const VERSION_FILE = "VERSION"

func main() {
	major := flag.Bool("major", false, "increment major version and reset minor and patch to 0")
	minor := flag.Bool("minor", false, "increment minor version and reset patch to 0")
	protocol := flag.String("protocol", "", "protocol name (e.g., solana)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                          # Increment patch version (e.g., 1.12.2 -> 1.12.3)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --minor                  # Increment minor version (e.g., 1.12.2 -> 1.13.0)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --major                  # Increment major version (e.g., 1.12.2 -> 2.0.0)\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --protocol=solana        # Use protocol-specific version bumping (e.g., solana-1.12.2 -> solana-1.12.3)\n", os.Args[0])
	}

	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	var versionPath string
	if *protocol != "" {
		versionPath = filepath.Join(dir, "plugins", *protocol, VERSION_FILE)
	} else {
		versionPath = filepath.Join(dir, VERSION_FILE)
	}

	if _, err := os.Stat(versionPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s file not found at %s\n", VERSION_FILE, versionPath)
		os.Exit(1)
	}

	currentVersion, err := readVersion(versionPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading version: %v\n", err)
		os.Exit(1)
	}

	originalVersion := currentVersion

	if *protocol != "" {
		prefix := *protocol + "-"
		if !strings.HasPrefix(currentVersion, prefix) {
			fmt.Fprintf(os.Stderr, "Error: version %s does not have expected prefix %s\n", currentVersion, prefix)
			os.Exit(1)
		}
		currentVersion = strings.TrimPrefix(currentVersion, prefix)
	}

	parts := strings.Split(currentVersion, ".")
	if len(parts) != 3 {
		fmt.Fprintf(os.Stderr, "Error: invalid version format. Expected X.Y.Z, got %s\n", currentVersion)
		os.Exit(1)
	}

	var newVersion string
	switch {
	case *major:
		newVersion = fmt.Sprintf("%d.0.0", atoi(parts[0])+1)
	case *minor:
		newVersion = fmt.Sprintf("%s.%d.0", parts[0], atoi(parts[1])+1)
	default:
		newVersion = fmt.Sprintf("%s.%s.%d", parts[0], parts[1], atoi(parts[2])+1)
	}

	makeVersion := newVersion
	if *protocol != "" {
		newVersion = *protocol + "-" + newVersion
	}

	if err := checkGitStatus(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	args := []string{"bump"}
	args = append(args, fmt.Sprintf("VERSION=%s", makeVersion))
	if *protocol != "" {
		args = append(args, fmt.Sprintf("PROTOCOL=%s", *protocol))
	}

	cmd := exec.Command("make", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running make bump: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully bumped version from %s to %s\n", originalVersion, newVersion)
}

func readVersion(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}
	return "", fmt.Errorf("empty version file")
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func checkGitStatus() error {
	mainBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	mainBranch, err := mainBranchCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	if strings.TrimSpace(string(mainBranch)) == "main" {
		return fmt.Errorf("cannot bump on main branch")
	}

	statusCmd := exec.Command("git", "diff-index", "--quiet", "HEAD")
	if err := statusCmd.Run(); err != nil {
		return fmt.Errorf("there are uncommitted changes")
	}

	return nil
}
