// Package utils provides extra utility
// for the program
package utils

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const ProgramName = "go-symphony"

// NonInteractiveCommand creates the command string from a flagSet
// to be used for getting the equivalent non-interactive shell command
func NonInteractiveCommand(use string, flagSet *pflag.FlagSet) string {
	nonInteractiveCommand := fmt.Sprintf("%s %s", ProgramName, use)

	visitFn := func(flag *pflag.Flag) {
		if flag.Name != "help" {
			if flag.Name == "feature" {
				featureFlagsString := ""
				// Creates string representation for the feature flags to be
				// concatenated with the nonInteractiveCommand
				for _, k := range strings.Split(flag.Value.String(), ",") {
					if k != "" {
						featureFlagsString += fmt.Sprintf(" --feature %s", k)
					}
				}
				nonInteractiveCommand += featureFlagsString
			} else if flag.Value.Type() == "bool" {
				if flag.Value.String() == "true" {
					nonInteractiveCommand = fmt.Sprintf("%s --%s", nonInteractiveCommand, flag.Name)
				}
			} else {
				flagValue := flag.Value.String()
				// Only include flags with non-empty values or flags that were explicitly changed from default
				if flagValue != "" && flag.Changed {
					// Skip SvelteKit flags if SvelteKit frontend is not enabled
					if isSvelteKitFlag(flag.Name) && !isSvelteKitEnabled(flagSet) {
						return
					}
					// Skip Supabase mode flag if Supabase driver is not selected
					if flag.Name == "supabase-mode" && !isSupabaseEnabled(flagSet) {
						return
					}
					nonInteractiveCommand = fmt.Sprintf("%s --%s %s", nonInteractiveCommand, flag.Name, flagValue)
				}
			}
		}
	}

	flagSet.SortFlags = false
	flagSet.VisitAll(visitFn)

	return nonInteractiveCommand
}

// isSvelteKitFlag checks if the flag is related to SvelteKit
func isSvelteKitFlag(flagName string) bool {
	svelteKitFlags := []string{"sveltekit-template", "sveltekit-types", "sveltekit-package-manager"}
	return slices.Contains(svelteKitFlags, flagName)
}

// isSvelteKitEnabled checks if SvelteKit frontend framework is selected
func isSvelteKitEnabled(flagSet *pflag.FlagSet) bool {
	frontendFlag := flagSet.Lookup("frontend")
	if frontendFlag == nil {
		return false
	}
	return frontendFlag.Value.String() == "sveltekit"
}

// isSupabaseEnabled checks if Supabase driver is selected
func isSupabaseEnabled(flagSet *pflag.FlagSet) bool {
	driverFlag := flagSet.Lookup("driver")
	if driverFlag == nil {
		return false
	}
	return driverFlag.Value.String() == "supabase"
}

func RegisterStaticCompletions(cmd *cobra.Command, flag string, options []string) {
	err := cmd.RegisterFlagCompletionFunc(flag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return options, cobra.ShellCompDirectiveNoFileComp
	})

	if err != nil {
		log.Printf("warning: could not register completion for --%s: %v", flag, err)
	}
}

// ExecuteCmd provides a shorthand way to run a shell command
func ExecuteCmd(name string, args []string, dir string) error {
	command := exec.Command(name, args...)
	command.Dir = dir
	var out bytes.Buffer
	var stdErr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stdErr
	if err := command.Run(); err != nil {
		return fmt.Errorf("%v\n%v", err, stdErr.String())
	}
	return nil
}

// InitGoMod initializes go.mod with the given project name
// in the selected directory
func InitGoMod(projectName string, appDir string) error {
	if err := ExecuteCmd("go",
		[]string{"mod", "init", projectName},
		appDir); err != nil {
		return err
	}

	return nil
}

// GoGetPackage runs "go get" for a given package in the
// selected directory
func GoGetPackage(appDir string, packages []string) error {
	for _, packageName := range packages {
		if err := ExecuteCmd("go",
			[]string{"get", "-u", packageName},
			appDir); err != nil {
			return err
		}
	}

	return nil
}

// GoFmt runs "gofmt" in a selected directory using the
// simplify and overwrite flags
func GoFmt(appDir string) error {
	if err := ExecuteCmd("gofmt",
		[]string{"-s", "-w", "."},
		appDir); err != nil {
		return err
	}

	return nil
}

// GoModReplace runs "go mod edit -replace" in the selected
// replace_payload e.g: github.com/gocql/gocql=github.com/scylladb/gocql@v1.14.4
func GoModReplace(appDir string, replace string) error {
	if err := ExecuteCmd("go",
		[]string{"mod", "edit", "-replace", replace},
		appDir,
	); err != nil {
		return err
	}

	return nil
}

func GoTidy(appDir string) error {
	err := ExecuteCmd("go", []string{"mod", "tidy"}, appDir)
	if err != nil {
		return err
	}
	return nil
}

func CheckGitConfig(key string) (bool, error) {
	cmd := exec.Command("git", "config", "--get", key)
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// The command failed to run.
			if exitError.ExitCode() == 1 {
				// The 'git config --get' command returns 1 if the key was not found.
				return false, nil
			}
		}
		// Some other error occurred.
		return false, err
	}
	// The command ran successfully, so the key is set.
	return true, nil
}

// ValidateModuleName returns true if it's a valid module name.
// It allows any number of / and . characters in between.
func ValidateModuleName(moduleName string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+(?:[\\/.][a-zA-Z0-9_-]+)*$", moduleName)
	return matched
}

// GetRootDir returns the project directory name from the module path.
// Returns the last token by splitting the moduleName with /
func GetRootDir(moduleName string) string {
	tokens := strings.Split(moduleName, "/")
	return tokens[len(tokens)-1]
}
