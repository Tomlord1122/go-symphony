package supabase

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Tomlord1122/go-symphony/v2/cmd/flags"
	tpl "github.com/Tomlord1122/go-symphony/v2/cmd/template"
)

//go:embed files/supabase_env.tmpl
var SupabaseEnvTemplate []byte

//go:embed files/supabase_sqlc.yaml.tmpl
var SupabaseSqlcConfigTemplate []byte

// SupabaseManager handles Supabase CLI operations
type SupabaseManager struct {
	ProjectPath string
	Mode        flags.SupabaseMode
}

// NewSupabaseManager creates a new Supabase manager
func NewSupabaseManager(projectPath string, mode flags.SupabaseMode) *SupabaseManager {
	return &SupabaseManager{
		ProjectPath: projectPath,
		Mode:        mode,
	}
}

// Init initializes Supabase project using supabase init command
// Handles interactive prompts for VS Code and IntelliJ Deno settings
func (sm *SupabaseManager) Init() error {
	cmd := exec.Command("supabase", "init")
	cmd.Dir = sm.ProjectPath

	// Create a pipe to handle stdin for interactive prompts
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	defer stdin.Close()

	// Capture stdout to monitor for prompts
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	cmd.Stderr = os.Stderr

	fmt.Println("Initializing Supabase project...")

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start Supabase init: %w", err)
	}

	// Handle interactive prompts by monitoring output and responding appropriately
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if err != nil {
				break
			}

			output := string(buffer[:n])
			fmt.Print(output) // Still show output to user

			// Check for VS Code Deno prompt
			if strings.Contains(output, "Generate VS Code settings for Deno") {
				io.WriteString(stdin, "N\n")
			}
			// Check for IntelliJ Deno prompt
			if strings.Contains(output, "Generate IntelliJ Settings for Deno") {
				io.WriteString(stdin, "N\n")
			}
		}
	}()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("failed to initialize Supabase: %w", err)
	}

	return nil
}

// Start starts local Supabase instance (only for local-db mode)
func (sm *SupabaseManager) Start() error {
	if sm.Mode != flags.LocalDB {
		return nil // Skip for init-only mode
	}

	cmd := exec.Command("supabase", "start")
	cmd.Dir = sm.ProjectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting local Supabase instance...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Supabase: %w", err)
	}

	return nil
}

// CreateMigration creates a new migration using supabase migration new
func (sm *SupabaseManager) CreateMigration(name string) (string, error) {
	cmd := exec.Command("supabase", "migration", "new", name)
	cmd.Dir = sm.ProjectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create migration: %w", err)
	}

	// Parse output to get migration path
	outputStr := strings.TrimSpace(string(output))
	fmt.Printf("Created migration: %s\n", outputStr)

	return outputStr, nil
}

// GetMigrationPath returns the migration directory path
func (sm *SupabaseManager) GetMigrationPath() string {
	return filepath.Join(sm.ProjectPath, "supabase", "migrations")
}

// GenerateSupabaseEnv creates .env file with Supabase configuration including global environment variables
func (sm *SupabaseManager) GenerateSupabaseEnv() error {
	envPath := filepath.Join(sm.ProjectPath, ".env")

	file, err := os.Create(envPath)
	if err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}
	defer file.Close()

	// Combine global environment template with Supabase-specific template
	envBytes := [][]byte{
		tpl.GlobalEnvTemplate(),
		SupabaseEnvTemplate,
	}

	combinedTemplate := bytes.Join(envBytes, []byte("\n"))
	_, err = file.Write(combinedTemplate)
	if err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	fmt.Println("Generated .env file with Supabase configuration")
	return nil
}

// GenerateSupabaseSqlcConfig creates sqlc.yaml with Supabase migration path
func (sm *SupabaseManager) GenerateSupabaseSqlcConfig() error {
	migrationPath := sm.GetMigrationPath()

	tmpl, err := template.New("sqlc").Parse(string(SupabaseSqlcConfigTemplate))
	if err != nil {
		return fmt.Errorf("failed to parse sqlc template: %w", err)
	}

	data := struct {
		MigrationPath string
	}{
		MigrationPath: migrationPath,
	}

	sqlcPath := filepath.Join(sm.ProjectPath, "sqlc.yaml")
	file, err := os.Create(sqlcPath)
	if err != nil {
		return fmt.Errorf("failed to create sqlc.yaml: %w", err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("failed to write sqlc.yaml: %w", err)
	}

	fmt.Println("Generated sqlc.yaml with Supabase migration path")
	return nil
}

// CheckSupabaseCLI checks if Supabase CLI is installed
func CheckSupabaseCLI() error {
	_, err := exec.LookPath("supabase")
	if err != nil {
		return fmt.Errorf("supabase CLI is not installed. Please install it from https://supabase.com/docs/guides/cli")
	}
	return nil
}
