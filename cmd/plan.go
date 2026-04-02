package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
	"github.com/Tomlord1122/go-symphony/cmd/program"
	"github.com/Tomlord1122/go-symphony/cmd/ui/textinput"
	"github.com/Tomlord1122/go-symphony/cmd/utils"
	"github.com/Tomlord1122/go-symphony/internal/scaffold"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Preview scaffold steps for a project configuration",
	Run: func(cmd *cobra.Command, args []string) {
		spec, err := buildPlanSpecFromCommand(cmd)
		if err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			cobra.CheckErr(err)
		}

		plan := scaffold.BuildPlan(spec, currentWorkingDir)
		if err := scaffold.WritePlan(os.Stdout, plan, spec.Execution.Output); err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	var flagDBDriver flags.Database
	var advancedFeatures flags.AdvancedFeatures
	var flagGit flags.Git
	var flagSupabaseMode flags.SupabaseMode
	var flagFrontendFramework flags.FrontendFramework
	var flagSvelteKitTemplate flags.SvelteKitTemplate
	var flagSvelteKitTypes flags.SvelteKitTypes
	var flagSvelteKitPackageManager flags.SvelteKitPackageManager

	planCmd.Flags().StringP("name", "n", "", "Name of project to create")
	planCmd.Flags().VarP(&flagDBDriver, "driver", "d", "Database driver to plan for")
	planCmd.Flags().Var(&advancedFeatures, "feature", "Advanced feature to include")
	planCmd.Flags().VarP(&flagGit, "git", "g", "Git mode to plan for")
	planCmd.Flags().Var(&flagSupabaseMode, "supabase-mode", "Supabase mode to plan for")
	planCmd.Flags().Var(&flagFrontendFramework, "frontend", "Frontend framework to plan for")
	planCmd.Flags().Var(&flagSvelteKitTemplate, "sveltekit-template", "SvelteKit template to plan for")
	planCmd.Flags().Var(&flagSvelteKitTypes, "sveltekit-types", "SvelteKit typing mode to plan for")
	planCmd.Flags().Var(&flagSvelteKitPackageManager, "sveltekit-package-manager", "SvelteKit package manager to plan for")
	planCmd.Flags().Bool("skip-install", false, "Skip dependency installation and formatting commands in the plan")
	planCmd.Flags().String("output", string(scaffold.OutputText), "Output format: text or json")

	utils.RegisterStaticCompletions(planCmd, "driver", flags.AllowedDBDrivers)
	utils.RegisterStaticCompletions(planCmd, "feature", flags.AllowedAdvancedFeatures)
	utils.RegisterStaticCompletions(planCmd, "git", flags.AllowedGitsOptions)
	utils.RegisterStaticCompletions(planCmd, "supabase-mode", flags.AllowedSupabaseModes)
	utils.RegisterStaticCompletions(planCmd, "frontend", flags.AllowedFrontendFrameworks)
	utils.RegisterStaticCompletions(planCmd, "output", []string{string(scaffold.OutputText), string(scaffold.OutputJSON)})
	utils.RegisterStaticCompletions(planCmd, "sveltekit-template", flags.AllowedSvelteKitTemplates)
	utils.RegisterStaticCompletions(planCmd, "sveltekit-types", flags.AllowedSvelteKitTypes)
	utils.RegisterStaticCompletions(planCmd, "sveltekit-package-manager", flags.AllowedSvelteKitPackageManagers)
}

func buildPlanSpecFromCommand(cmd *cobra.Command) (scaffold.CreateSpec, error) {
	flagSkipInstall, err := cmd.Flags().GetBool("skip-install")
	if err != nil {
		log.Fatal("failed to retrieve skip-install flag")
	}

	featureFlags := cmd.Flag("feature").Value.String()
	featureValues := []string{}
	if featureFlags != "" {
		featureValues = strings.Split(featureFlags, ",")
	}

	project := &program.Project{
		ProjectName:     cmd.Flag("name").Value.String(),
		DBDriver:        flags.Database(cmd.Flag("driver").Value.String()),
		GitOptions:      flags.Git(cmd.Flag("git").Value.String()),
		AdvancedOptions: make(map[string]bool),
	}
	for _, feature := range featureValues {
		project.AdvancedOptions[strings.ToLower(strings.TrimSpace(feature))] = true
	}

	spec := buildScaffoldSpec(cmd, project, nil, scaffold.ExecutionOptions{
		DryRun:        true,
		NoInteractive: true,
		SkipInstall:   flagSkipInstall,
		Output:        scaffold.OutputFormat(cmd.Flag("output").Value.String()),
	})

	if err := scaffold.ValidateSpec(spec); err != nil {
		return scaffold.CreateSpec{}, err
	}

	return spec, nil
}
