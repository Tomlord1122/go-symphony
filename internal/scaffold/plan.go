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
	Kind          StepKind `json:"kind"`
	Name          string   `json:"name"`
	Path          string   `json:"path,omitempty"`
	Command       []string `json:"command,omitempty"`
	Optional      bool     `json:"optional,omitempty"`
	Description   string   `json:"description,omitempty"`
	RequiresTools []string `json:"requires_tools,omitempty"`
	SkippedReason string   `json:"skipped_reason,omitempty"`
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
			Step{Kind: StepRunCommand, Name: "Install database dependencies", Path: projectPath, Command: []string{"go", "get", "database dependencies"}, Optional: true, RequiresTools: []string{"go"}},
		)
	}

	if spec.HasFeature(flags.Sqlc) {
		steps = append(steps, Step{Kind: StepWriteFile, Name: "Write SQLC setup", Path: filepath.Join(projectPath, "sqlc.yaml")})
	}

	if spec.HasFeature(flags.Docker) {
		steps = append(steps,
			Step{Kind: StepWriteFile, Name: "Write Docker assets", Path: filepath.Join(projectPath, "Dockerfile")},
		)
		if spec.DBDriver != flags.None && spec.DBDriver != flags.Supabase {
			steps = append(steps, Step{Kind: StepWriteFile, Name: "Write Docker Compose assets", Path: filepath.Join(projectPath, "docker-compose.yml")})
		}
	}

	if spec.HasFeature(flags.GoProjectWorkflow) {
		steps = append(steps, Step{Kind: StepWriteFile, Name: "Write GitHub Actions workflows", Path: filepath.Join(projectPath, ".github", "workflows")})
	}

	if spec.Frontend.Framework == flags.SvelteKitFrontend {
		steps = append(steps,
			Step{Kind: StepInfo, Name: "SvelteKit bootstrap requirements", Description: "Frontend bootstrap uses external Node.js tooling", Optional: true, RequiresTools: []string{"node", "npx"}},
			Step{Kind: StepRunCommand, Name: "Bootstrap SvelteKit frontend", Path: baseDir, Command: []string{"npx", "sv", "create", spec.RootDir() + "-frontend"}, Optional: true, RequiresTools: []string{"node", "npx"}},
		)
	}

	if spec.Frontend.Framework == flags.NextJSFrontend {
		steps = append(steps,
			Step{Kind: StepInfo, Name: "Next.js bootstrap requirements", Description: "Frontend bootstrap uses external Node.js tooling", Optional: true, RequiresTools: []string{"node", "npx"}},
			Step{Kind: StepRunCommand, Name: "Bootstrap Next.js frontend", Path: baseDir, Command: []string{"npx", "create-next-app@latest", spec.RootDir() + "-frontend"}, Optional: true, RequiresTools: []string{"node", "npx"}},
		)
	}

	if spec.DBDriver == flags.Supabase {
		steps = append(steps,
			Step{Kind: StepInfo, Name: "Supabase bootstrap requirements", Description: "Supabase setup depends on the Supabase CLI", Optional: true, RequiresTools: []string{"supabase"}},
			Step{Kind: StepRunCommand, Name: "Initialize Supabase", Path: projectPath, Command: []string{"supabase", "init"}, Optional: true, RequiresTools: []string{"supabase"}},
		)
	}

	if !spec.Execution.SkipInstall {
		steps = append(steps,
			Step{Kind: StepRunCommand, Name: "Run go mod tidy", Path: projectPath, Command: []string{"go", "mod", "tidy"}, Optional: true, RequiresTools: []string{"go"}},
			Step{Kind: StepRunCommand, Name: "Run gofmt", Path: projectPath, Command: []string{"gofmt", "-s", "-w", "."}, Optional: true, RequiresTools: []string{"gofmt"}},
		)
	} else {
		steps = append(steps,
			Step{Kind: StepInfo, Name: "Skip install commands", Optional: true, SkippedReason: "skip-install enabled", Description: "Dependency installation and formatting commands will not run"},
		)
	}

	if spec.GitMode != "" && spec.GitMode != flags.Skip {
		steps = append(steps, Step{Kind: StepRunCommand, Name: "Initialize git repository", Path: projectPath, Command: []string{"git", "init"}, Optional: true, RequiresTools: []string{"git"}})
	}

	return Plan{Spec: spec, Steps: steps}
}
