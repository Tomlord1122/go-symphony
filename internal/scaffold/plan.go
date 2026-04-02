package scaffold

import (
	"path/filepath"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
)

type StepKind string

const (
	StepCreateDir  StepKind = "create_dir"
	StepWriteFile  StepKind = "write_file"
	StepRunCommand StepKind = "run_command"
	StepInfo       StepKind = "info"
)

type Step struct {
	Kind        StepKind `json:"kind"`
	Name        string   `json:"name"`
	Path        string   `json:"path,omitempty"`
	Command     []string `json:"command,omitempty"`
	Optional    bool     `json:"optional,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Plan struct {
	Spec  CreateSpec `json:"spec"`
	Steps []Step     `json:"steps"`
}

func BuildPlan(spec CreateSpec, baseDir string) Plan {
	projectPath := filepath.Join(baseDir, spec.RootDir())
	steps := []Step{
		{Kind: StepCreateDir, Name: "Create project root", Path: projectPath},
		{Kind: StepRunCommand, Name: "Initialize Go module", Path: projectPath, Command: []string{"go", "mod", "init", spec.ProjectName}},
		{Kind: StepWriteFile, Name: "Write backend scaffold files", Path: projectPath, Description: "Generate Go backend templates, README, Makefile, and config files"},
	}

	if spec.DBDriver != "" && spec.DBDriver != flags.None {
		steps = append(steps,
			Step{Kind: StepRunCommand, Name: "Install database dependencies", Path: projectPath, Command: []string{"go", "get", "database dependencies"}, Optional: true},
		)
	}

	if spec.HasFeature(flags.Sqlc) {
		steps = append(steps, Step{Kind: StepWriteFile, Name: "Write SQLC setup", Path: filepath.Join(projectPath, "sqlc.yaml")})
	}

	if spec.HasFeature(flags.Docker) {
		steps = append(steps, Step{Kind: StepWriteFile, Name: "Write Docker assets", Path: filepath.Join(projectPath, "Dockerfile")})
	}

	if spec.HasFeature(flags.GoProjectWorkflow) {
		steps = append(steps, Step{Kind: StepWriteFile, Name: "Write GitHub Actions workflows", Path: filepath.Join(projectPath, ".github", "workflows")})
	}

	if spec.Frontend.Framework == flags.SvelteKitFrontend {
		steps = append(steps, Step{Kind: StepRunCommand, Name: "Bootstrap SvelteKit frontend", Path: baseDir, Command: []string{"npx", "sv", "create", spec.RootDir() + "-frontend"}, Optional: true})
	}

	if spec.Frontend.Framework == flags.NextJSFrontend {
		steps = append(steps, Step{Kind: StepRunCommand, Name: "Bootstrap Next.js frontend", Path: baseDir, Command: []string{"npx", "create-next-app@latest", spec.RootDir() + "-frontend"}, Optional: true})
	}

	if spec.DBDriver == flags.Supabase {
		steps = append(steps, Step{Kind: StepRunCommand, Name: "Initialize Supabase", Path: projectPath, Command: []string{"supabase", "init"}, Optional: true})
	}

	if !spec.Execution.SkipInstall {
		steps = append(steps,
			Step{Kind: StepRunCommand, Name: "Run go mod tidy", Path: projectPath, Command: []string{"go", "mod", "tidy"}, Optional: true},
			Step{Kind: StepRunCommand, Name: "Run gofmt", Path: projectPath, Command: []string{"gofmt", "-s", "-w", "."}, Optional: true},
		)
	}

	if spec.GitMode != "" && spec.GitMode != flags.Skip {
		steps = append(steps, Step{Kind: StepRunCommand, Name: "Initialize git repository", Path: projectPath, Command: []string{"git", "init"}, Optional: true})
	}

	return Plan{Spec: spec, Steps: steps}
}
