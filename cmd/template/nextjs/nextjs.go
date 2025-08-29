package nextjs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
)

// NextJSManager handles Next.js project creation
type NextJSManager struct {
	ProjectPath    string
	FrontendName   string
	TypeScript     bool
	ESLint         bool
	TailwindCSS    bool
	SrcDirectory   bool
	AppRouter      bool
	ImportAlias    string
	PackageManager flags.SvelteKitPackageManager // Reuse the same package manager types
}

// NewNextJSManager creates a new Next.js manager
func NewNextJSManager(projectPath, frontendName string) *NextJSManager {
	return &NextJSManager{
		ProjectPath:    projectPath,
		FrontendName:   frontendName,
		TypeScript:     true,       // Default to TypeScript
		ESLint:         true,       // Default to ESLint
		TailwindCSS:    false,      // Default to no Tailwind
		SrcDirectory:   false,      // Default to no src directory
		AppRouter:      true,       // Default to App Router
		ImportAlias:    "@/*",      // Default import alias
		PackageManager: flags.PNPM, // Default to pnpm (as requested by user)
	}
}

// SetOptions configures the Next.js project options
func (nm *NextJSManager) SetOptions(packageManager flags.SvelteKitPackageManager, typescript, eslint, tailwind, srcDir, appRouter bool, importAlias string) {
	if packageManager != "" {
		nm.PackageManager = packageManager
	}
	nm.TypeScript = typescript
	nm.ESLint = eslint
	nm.TailwindCSS = tailwind
	nm.SrcDirectory = srcDir
	nm.AppRouter = appRouter
	if importAlias != "" {
		nm.ImportAlias = importAlias
	}
}

// Create creates a new Next.js project using npx create-next-app@latest
func (nm *NextJSManager) Create() error {
	// Build the command arguments
	args := []string{"create-next-app@latest"}

	// Add the project path
	frontendPath := filepath.Join(nm.ProjectPath, nm.FrontendName)
	args = append(args, frontendPath)

	// Add TypeScript option
	if nm.TypeScript {
		args = append(args, "--typescript")
	} else {
		args = append(args, "--javascript")
	}

	// Add ESLint option
	if nm.ESLint {
		args = append(args, "--eslint")
	} else {
		args = append(args, "--no-eslint")
	}

	// Add Tailwind CSS option
	if nm.TailwindCSS {
		args = append(args, "--tailwind")
	} else {
		args = append(args, "--no-tailwind")
	}

	// Add src directory option
	if nm.SrcDirectory {
		args = append(args, "--src-dir")
	} else {
		args = append(args, "--no-src-dir")
	}

	// Add App Router option
	if nm.AppRouter {
		args = append(args, "--app")
	} else {
		args = append(args, "--no-app")
	}

	// Add import alias
	if nm.ImportAlias != "" {
		args = append(args, "--import-alias", nm.ImportAlias)
	}

	// Add package manager option
	if nm.PackageManager != "" {
		args = append(args, "--use-"+string(nm.PackageManager))
	}

	fmt.Printf("Creating Next.js project: %s\n", nm.FrontendName)
	fmt.Printf("Command: npx %s\n", strings.Join(args, " "))

	// Execute the command
	cmd := exec.Command("npx", args...)
	cmd.Dir = nm.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Allow interactive input if needed

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create Next.js project: %w", err)
	}

	return nil
}

// CreateInteractive creates a Next.js project with interactive prompts
func (nm *NextJSManager) CreateInteractive() error {
	frontendPath := filepath.Join(nm.ProjectPath, nm.FrontendName)

	fmt.Printf("Creating Next.js project: %s\n", nm.FrontendName)
	fmt.Println("This will run npx create-next-app@latest with interactive prompts...")
	fmt.Println("You can press Ctrl+C to cancel at any time.")
	fmt.Println("--------------------------------")

	// Simple command for interactive mode
	cmd := exec.Command("npx", "create-next-app@latest", frontendPath)
	cmd.Dir = nm.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Check if it was cancelled by user
		if cmd.ProcessState != nil && !cmd.ProcessState.Success() {
			return fmt.Errorf("next.js project creation was cancelled or failed")
		}
		return fmt.Errorf("failed to create Next.js project: %w", err)
	}

	return nil
}

// GetFrontendPath returns the full path to the frontend project
func (nm *NextJSManager) GetFrontendPath() string {
	return filepath.Join(nm.ProjectPath, nm.FrontendName)
}

// CheckNodeJS checks if Node.js is installed
func CheckNodeJS() error {
	_, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("node.js is not installed. please install it from https://nodejs.org/")
	}
	return nil
}

// CheckNPX checks if npx is available
func CheckNPX() error {
	_, err := exec.LookPath("npx")
	if err != nil {
		return fmt.Errorf("npx is not available. Please ensure Node.js is properly installed")
	}
	return nil
}

// CheckPNPM checks if pnpm is installed (user preference)
func CheckPNPM() error {
	_, err := exec.LookPath("pnpm")
	if err != nil {
		return fmt.Errorf("pnpm is not installed. Please install it with: npm install -g pnpm")
	}
	return nil
}
