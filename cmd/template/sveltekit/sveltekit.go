package sveltekit

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Tomlord1122/go-symphony/v2/cmd/flags"
)

// SvelteKitManager handles SvelteKit project creation
type SvelteKitManager struct {
	ProjectPath    string
	FrontendName   string
	Template       flags.SvelteKitTemplate
	Types          flags.SvelteKitTypes
	PackageManager flags.SvelteKitPackageManager
	SkipAddOns     bool
	SkipInstall    bool
}

// NewSvelteKitManager creates a new SvelteKit manager
func NewSvelteKitManager(projectPath, frontendName string) *SvelteKitManager {
	return &SvelteKitManager{
		ProjectPath:    projectPath,
		FrontendName:   frontendName,
		Template:       flags.Demo,       // Default to demo template
		Types:          flags.TypeScript, // Default to TypeScript
		PackageManager: flags.PNPM,       // Default to pnpm (as requested by user)
		SkipAddOns:     false,
		SkipInstall:    false,
	}
}

// SetOptions configures the SvelteKit project options
func (sm *SvelteKitManager) SetOptions(template flags.SvelteKitTemplate, types flags.SvelteKitTypes, packageManager flags.SvelteKitPackageManager) {
	if template != "" {
		sm.Template = template
	}
	if types != "" {
		sm.Types = types
	}
	if packageManager != "" {
		sm.PackageManager = packageManager
	}
}

// Create creates a new SvelteKit project using npx sv create
func (sm *SvelteKitManager) Create() error {
	// Build the command arguments
	args := []string{"sv", "create"}

	// Add template option
	if sm.Template != "" {
		args = append(args, "--template", string(sm.Template))
	}

	// Add types option
	if sm.Types != "" && sm.Types != flags.NoTypes {
		args = append(args, "--types", string(sm.Types))
	}

	// Add package manager option
	if sm.PackageManager != "" {
		args = append(args, "--install", string(sm.PackageManager))
	}

	// Skip add-ons if requested
	if sm.SkipAddOns {
		args = append(args, "--no-add-ons")
	}

	// Skip install if requested
	if sm.SkipInstall {
		args = append(args, "--no-install")
	}

	// Add the project path
	frontendPath := filepath.Join(sm.ProjectPath, sm.FrontendName)
	args = append(args, frontendPath)

	fmt.Printf("Creating SvelteKit project: %s\n", sm.FrontendName)
	fmt.Printf("Command: npx %s\n", strings.Join(args, " "))

	// Execute the command
	cmd := exec.Command("npx", args...)
	cmd.Dir = sm.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // Allow interactive input if needed

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create SvelteKit project: %w", err)
	}

	return nil
}

// CreateInteractive creates a SvelteKit project with interactive prompts
func (sm *SvelteKitManager) CreateInteractive() error {
	frontendPath := filepath.Join(sm.ProjectPath, sm.FrontendName)

	fmt.Printf("Creating SvelteKit project: %s\n", sm.FrontendName)
	fmt.Println("This will run npx sv create with interactive prompts...")
	fmt.Println("You can press Ctrl+C to cancel at any time.")
	fmt.Println("--------------------------------")

	// Simple command for interactive mode
	cmd := exec.Command("npx", "sv", "create", frontendPath)
	cmd.Dir = sm.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// Check if it was cancelled by user
		if cmd.ProcessState != nil && !cmd.ProcessState.Success() {
			return fmt.Errorf("SvelteKit project creation was cancelled or failed")
		}
		return fmt.Errorf("failed to create SvelteKit project: %w", err)
	}

	return nil
}

// GetFrontendPath returns the full path to the frontend project
func (sm *SvelteKitManager) GetFrontendPath() string {
	return filepath.Join(sm.ProjectPath, sm.FrontendName)
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
