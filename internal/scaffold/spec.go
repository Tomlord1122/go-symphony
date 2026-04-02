package scaffold

import (
	"fmt"
	"strings"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
	"github.com/Tomlord1122/go-symphony/cmd/utils"
)

type OutputFormat string

const (
	OutputText OutputFormat = "text"
	OutputJSON OutputFormat = "json"
)

type ExecutionOptions struct {
	DryRun        bool
	NoInteractive bool
	SkipInstall   bool
	Output        OutputFormat
}

type FrontendSpec struct {
	Framework      flags.FrontendFramework       `json:"framework"`
	SvelteTemplate flags.SvelteKitTemplate       `json:"svelte_template,omitempty"`
	SvelteTypes    flags.SvelteKitTypes          `json:"svelte_types,omitempty"`
	PackageManager flags.SvelteKitPackageManager `json:"package_manager,omitempty"`
}

type CreateSpec struct {
	ProjectName      string             `json:"project_name"`
	ProjectType      flags.Framework    `json:"project_type"`
	DBDriver         flags.Database     `json:"db_driver"`
	AdvancedFeatures []string           `json:"advanced_features,omitempty"`
	GitMode          flags.Git          `json:"git_mode"`
	SupabaseMode     flags.SupabaseMode `json:"supabase_mode,omitempty"`
	Frontend         FrontendSpec       `json:"frontend"`
	Execution        ExecutionOptions   `json:"execution"`
}

func NormalizeAdvancedFeatures(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		result = append(result, normalized)
	}
	return result
}

func BuildSpec(
	name string,
	dbDriver flags.Database,
	advancedFeatures []string,
	gitMode flags.Git,
	supabaseMode flags.SupabaseMode,
	frontend flags.FrontendFramework,
	svelteTemplate flags.SvelteKitTemplate,
	svelteTypes flags.SvelteKitTypes,
	packageManager flags.SvelteKitPackageManager,
	exec ExecutionOptions,
) CreateSpec {
	return CreateSpec{
		ProjectName:      strings.TrimSpace(name),
		ProjectType:      flags.Gin,
		DBDriver:         dbDriver,
		AdvancedFeatures: NormalizeAdvancedFeatures(advancedFeatures),
		GitMode:          gitMode,
		SupabaseMode:     supabaseMode,
		Frontend: FrontendSpec{
			Framework:      frontend,
			SvelteTemplate: svelteTemplate,
			SvelteTypes:    svelteTypes,
			PackageManager: packageManager,
		},
		Execution: exec,
	}
}

func (s CreateSpec) RootDir() string {
	return utils.GetRootDir(s.ProjectName)
}

func (s CreateSpec) HasFeature(feature string) bool {
	needle := strings.ToLower(strings.TrimSpace(feature))
	for _, candidate := range s.AdvancedFeatures {
		if candidate == needle {
			return true
		}
	}
	return false
}

func ValidateSpec(spec CreateSpec) error {
	if spec.ProjectType == "" {
		spec.ProjectType = flags.Gin
	}

	if spec.ProjectName == "" {
		if spec.Execution.NoInteractive {
			return fmt.Errorf("project name is required when --no-interactive is set")
		}
		return nil
	}

	if !utils.ValidateModuleName(spec.ProjectName) {
		return fmt.Errorf("'%s' is not a valid module name. Please choose a different name", spec.ProjectName)
	}

	if spec.DBDriver == "" && spec.Execution.NoInteractive {
		return fmt.Errorf("database driver is required when --no-interactive is set")
	}

	if spec.GitMode == "" && spec.Execution.NoInteractive {
		return fmt.Errorf("git mode is required when --no-interactive is set")
	}

	if spec.Frontend.Framework == "" && spec.Execution.NoInteractive {
		return fmt.Errorf("frontend framework is required when --no-interactive is set")
	}

	if spec.SupabaseMode != "" && spec.DBDriver != flags.Supabase {
		return fmt.Errorf("supabase mode can only be used with the supabase driver")
	}

	if spec.DBDriver == flags.Supabase && spec.SupabaseMode == "" && spec.Execution.NoInteractive {
		return fmt.Errorf("supabase mode is required when driver is supabase and --no-interactive is set")
	}

	if spec.Frontend.Framework != flags.SvelteKitFrontend {
		if spec.Frontend.SvelteTemplate != "" || spec.Frontend.SvelteTypes != "" {
			return fmt.Errorf("sveltekit options can only be used when frontend is sveltekit")
		}
	}

	for _, feature := range spec.AdvancedFeatures {
		allowed := false
		for _, candidate := range flags.AllowedAdvancedFeatures {
			if feature == candidate {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("advanced feature to use. Allowed values: %s", strings.Join(flags.AllowedAdvancedFeatures, ", "))
		}
	}

	return nil
}
