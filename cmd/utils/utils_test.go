package utils

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestValidateModuleName(t *testing.T) {
	passTestCases := []string{
		"github.com/user/project",
		"github.com/user/projec1-hyphen",
		"github.com/user/projecT_under_Score",
		"github.com/user/project.hyphen3",
		"project",
		"ProJEct",
		"PRo_45-.4Jc",
		"PRo_/4J/c",
	}
	for _, testCase := range passTestCases {
		ok := ValidateModuleName(testCase)
		if !ok {
			t.Errorf("testing:%s expected:true got:%v", testCase, ok)
		}
	}

	failTestCases := []string{
		"",
		"/",
		".",
		"//",
		"/project",
		"ProJEct/",
		"PRo_$4Jc",
		"PRo_@J/c",
	}
	for _, testCase := range failTestCases {
		ok := ValidateModuleName(testCase)
		if ok {
			t.Errorf("testing:%s expected:false got:%v", testCase, ok)
		}
	}
}

func TestGeRootDir(t *testing.T) {
	testCases := map[string]string{
		"github.com/user/pro-ject": "pro-ject",
		"pro-ject":                 "pro-ject",
		"/":                        "",
		"":                         "",
		"//":                       "",
		"@":                        "@",
	}

	for intput, output := range testCases {
		rootDir := GetRootDir(intput)
		if rootDir != output {
			t.Errorf("testing:%s expected:%s got:%s", intput, output, rootDir)
		}
	}
}

func TestNonInteractiveCommand(t *testing.T) {
	tests := []struct {
		name        string
		setupFlags  func(*pflag.FlagSet)
		expected    string
		description string
	}{
		{
			name: "basic command with required flags",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.String("name", "", "project name")
				flags.String("driver", "", "database driver")
				flags.String("git", "", "git option")
				flags.Set("name", "test-project")
				flags.Set("driver", "postgres")
				flags.Set("git", "commit")
			},
			expected:    "go-symphony create --name test-project --driver postgres --git commit",
			description: "Should include basic flags with values",
		},
		{
			name: "exclude empty sveltekit flags when sveltekit not enabled",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.String("name", "", "project name")
				flags.String("driver", "", "database driver")
				flags.String("frontend", "", "frontend framework")
				flags.String("sveltekit-template", "", "sveltekit template")
				flags.String("sveltekit-types", "", "sveltekit types")
				flags.Set("name", "test-project")
				flags.Set("driver", "postgres")
				flags.Set("frontend", "none")
				// Don't set sveltekit flags values - they should be excluded
			},
			expected:    "go-symphony create --name test-project --driver postgres --frontend none",
			description: "Should exclude SvelteKit flags when SvelteKit feature is not enabled",
		},
		{
			name: "include sveltekit flags when sveltekit enabled",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.String("name", "", "project name")
				flags.String("driver", "", "database driver")
				flags.String("frontend", "", "frontend framework")
				flags.String("sveltekit-template", "", "sveltekit template")
				flags.String("sveltekit-types", "", "sveltekit types")
				flags.Set("name", "test-project")
				flags.Set("driver", "postgres")
				flags.Set("frontend", "sveltekit")
				flags.Set("sveltekit-template", "skeleton")
				flags.Set("sveltekit-types", "typescript")
			},
			expected:    "go-symphony create --name test-project --driver postgres --frontend sveltekit --sveltekit-template skeleton --sveltekit-types typescript",
			description: "Should include SvelteKit flags when SvelteKit feature is enabled",
		},
		{
			name: "exclude supabase-mode when driver is not supabase",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.String("name", "", "project name")
				flags.String("driver", "", "database driver")
				flags.String("supabase-mode", "", "supabase mode")
				flags.Set("name", "test-project")
				flags.Set("driver", "postgres")
				flags.Set("supabase-mode", "local-db")
			},
			expected:    "go-symphony create --name test-project --driver postgres",
			description: "Should exclude supabase-mode flag when driver is not supabase",
		},
		{
			name: "include supabase-mode when driver is supabase",
			setupFlags: func(flags *pflag.FlagSet) {
				flags.String("name", "", "project name")
				flags.String("driver", "", "database driver")
				flags.String("supabase-mode", "", "supabase mode")
				flags.Set("name", "test-project")
				flags.Set("driver", "supabase")
				flags.Set("supabase-mode", "local-db")
			},
			expected:    "go-symphony create --name test-project --driver supabase --supabase-mode local-db",
			description: "Should include supabase-mode flag when driver is supabase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.setupFlags(flags)

			result := NonInteractiveCommand("create", flags)
			if result != tt.expected {
				t.Errorf("Test %s failed:\nExpected: %s\nGot:      %s\nDescription: %s",
					tt.name, tt.expected, result, tt.description)
			}
		})
	}
}
